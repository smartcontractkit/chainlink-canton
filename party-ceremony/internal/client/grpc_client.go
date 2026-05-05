package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-canton/party-ceremony/internal/helpers"

	participantv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/admin/participant/v30"
	cryptoadminv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/admin/v30"
	cryptov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/crypto/v30"
	protov30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/protocol/v30"
	topoadminv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/topology/admin/v30"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

const defaultHost = "localhost"
const defaultPort = 5001

// Dial opens a gRPC connection to a Canton admin API endpoint using the
// host, port, and optional JWT bearer token from cfg.
func Dial(cfg ClientConfig) (*grpc.ClientConn, error) {
	host := cfg.AdminHost
	if host == "" {
		host = defaultHost
	}
	port := cfg.AdminPort
	if port == 0 {
		port = defaultPort
	}
	target := fmt.Sprintf("%s:%d", host, port)

	creds := transportCredentials(host, port)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	if cfg.AdminJWT != "" {
		opts = append(opts, grpc.WithUnaryInterceptor(jwtUnaryInterceptor(cfg.AdminJWT)))
		opts = append(opts, grpc.WithStreamInterceptor(jwtStreamInterceptor(cfg.AdminJWT)))
	}

	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		return nil, fmt.Errorf("grpc dial %s: %w", target, err)
	}

	return conn, nil
}

// transportCredentials returns TLS credentials when dialing port 443
// (DevNet exposes public Admin/Ledger gRPC through TLS on 443),
// and plaintext/insecure credentials otherwise.
func transportCredentials(host string, port int) credentials.TransportCredentials {
	if port == 443 {
		return credentials.NewClientTLSFromCert(nil, host)
	}

	return insecure.NewCredentials()
}

func jwtUnaryInterceptor(token string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token), method, req, reply, cc, opts...)
	}
}

func jwtStreamInterceptor(token string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return streamer(metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token), desc, cc, method, opts...)
	}
}

// GRPCCantonClient implements [CantonClient] using real Canton Admin gRPC APIs.
type GRPCCantonClient struct {
	conn   *grpc.ClientConn
	topo   topoadminv30.TopologyManagerWriteServiceClient
	reader topoadminv30.TopologyManagerReadServiceClient
	vault  cryptoadminv30.VaultServiceClient
	id     topoadminv30.IdentityInitializationServiceClient
	pkg    participantv30.PackageServiceClient
}

// NewGRPCClient creates a new [GRPCCantonClient] from an established gRPC connection.
func NewGRPCClient(conn *grpc.ClientConn) *GRPCCantonClient {
	return &GRPCCantonClient{
		conn:   conn,
		topo:   topoadminv30.NewTopologyManagerWriteServiceClient(conn),
		reader: topoadminv30.NewTopologyManagerReadServiceClient(conn),
		vault:  cryptoadminv30.NewVaultServiceClient(conn),
		id:     topoadminv30.NewIdentityInitializationServiceClient(conn),
		pkg:    participantv30.NewPackageServiceClient(conn),
	}
}

// ── Store helpers ────────────────────────────────────────────────────────────

var authorizedStore = &topoadminv30.StoreId{
	Store: &topoadminv30.StoreId_Authorized_{
		Authorized: &topoadminv30.StoreId_Authorized{},
	},
}

func synchronizerStore(syncID string) *topoadminv30.StoreId {
	return &topoadminv30.StoreId{
		Store: &topoadminv30.StoreId_Synchronizer{
			Synchronizer: &topoadminv30.Synchronizer{
				Kind: &topoadminv30.Synchronizer_Id{Id: syncID},
			},
		},
	}
}

func storeForSync(syncID string) *topoadminv30.StoreId {
	if syncID == "" {
		return authorizedStore
	}

	return synchronizerStore(syncID)
}

func headStateBaseQuery(store *topoadminv30.StoreId) *topoadminv30.BaseQuery {
	return &topoadminv30.BaseQuery{
		Store: store,
		TimeQuery: &topoadminv30.BaseQuery_HeadState{
			HeadState: &emptypb.Empty{},
		},
	}
}

// ── Identity ─────────────────────────────────────────────────────────────────

func (c *GRPCCantonClient) GetParticipantUID(ctx context.Context) (string, error) {
	resp, err := c.id.GetId(ctx, &topoadminv30.GetIdRequest{})
	if err != nil {
		return "", fmt.Errorf("GetId: %w", err)
	}

	return resp.GetUniqueIdentifier(), nil
}

func (c *GRPCCantonClient) GetParticipantID(ctx context.Context) (string, error) {
	uid, err := c.GetParticipantUID(ctx)
	if err != nil {
		return "", err
	}
	// UID format: "<id>::<fingerprint>"; return the portion before "::"
	if before, _, found := strings.Cut(uid, "::"); found {
		return before, nil
	}

	return uid, nil
}

// ── Key Management ───────────────────────────────────────────────────────────

func (c *GRPCCantonClient) GenerateSigningKey(
	ctx context.Context,
	name string,
	usage []cryptov30.SigningKeyUsage,
) (*cryptov30.SigningPublicKey, error) {
	resp, err := c.vault.GenerateSigningKey(ctx, &cryptoadminv30.GenerateSigningKeyRequest{
		Name:  name,
		Usage: usage,
	})
	if err != nil {
		return nil, fmt.Errorf("GenerateSigningKey: %w", err)
	}

	return resp.GetPublicKey(), nil
}

func (c *GRPCCantonClient) RegisterKmsSigningKey(
	ctx context.Context,
	kmsKeyID string,
	name string,
	usage []cryptov30.SigningKeyUsage,
) (*cryptov30.SigningPublicKey, error) {
	resp, err := c.vault.RegisterKmsSigningKey(ctx, &cryptoadminv30.RegisterKmsSigningKeyRequest{
		KmsKeyId: kmsKeyID,
		Name:     name,
		Usage:    usage,
	})
	if err != nil {
		return nil, fmt.Errorf("RegisterKmsSigningKey: %w", err)
	}

	return resp.GetPublicKey(), nil
}

func (c *GRPCCantonClient) GetNamespaceFingerprint(ctx context.Context, keyName string, synchronizerID string, knownOwners []string) (string, error) {
	// Step 1: find ALL keys named keyName in the vault.
	vaultResp, err := c.vault.ListMyKeys(ctx, &cryptoadminv30.ListMyKeysRequest{
		Filters: &cryptoadminv30.ListKeysFilters{
			Name:  keyName,
			Usage: []cryptov30.SigningKeyUsage{cryptov30.SigningKeyUsage_SIGNING_KEY_USAGE_NAMESPACE},
		},
	})
	if err != nil {
		return "", fmt.Errorf("ListMyKeys: %w", err)
	}

	var allKeyBytes [][]byte
	for _, km := range vaultResp.GetPrivateKeysMetadata() {
		pkn := km.GetPublicKeyWithName()
		if pkn.GetName() == keyName {
			if spk := pkn.GetPublicKey().GetSigningPublicKey(); spk != nil {
				allKeyBytes = append(allKeyBytes, spk.GetPublicKey())
			}
		}
	}
	if len(allKeyBytes) == 0 {
		return "", fmt.Errorf("no NAMESPACE signing key named %q found in vault", keyName)
	}

	// Build a fast lookup set for knownOwners (when provided).
	ownerSet := make(map[string]struct{}, len(knownOwners))
	for _, o := range knownOwners {
		ownerSet[o] = struct{}{}
	}

	// Step 2: list NSDs in the synchronizer store.
	nsdResp, err := c.reader.ListNamespaceDelegation(ctx, &topoadminv30.ListNamespaceDelegationRequest{
		BaseQuery: headStateBaseQuery(storeForSync(synchronizerID)),
	})
	if err != nil {
		return "", fmt.Errorf("ListNamespaceDelegation: %w", err)
	}

	// Step 3: find the NSD whose namespace is in knownOwners (if provided)
	// AND whose target key matches one of the participant's vault keys.
	for _, r := range nsdResp.GetResults() {
		item := r.GetItem()
		ns := item.GetNamespace()
		if len(ownerSet) > 0 {
			if _, ok := ownerSet[ns]; !ok {
				continue
			}
		}
		for _, kb := range allKeyBytes {
			if bytes.Equal(item.GetTargetKey().GetPublicKey(), kb) {
				return ns, nil
			}
		}
	}

	return "", fmt.Errorf(
		"no namespace delegation for key %q found; ensure propose_delegation has been submitted", keyName)
}

// GetNamespaceKeyName discovers the human-readable name of this participant's
// NAMESPACE signing key that belongs to the decentralized namespace identified
// by knownOwners. It lists all NAMESPACE keys in the vault (without filtering
// by name), cross-references them with namespace delegations, and returns the
// name of the matching key.
func (c *GRPCCantonClient) GetNamespaceKeyName(ctx context.Context, synchronizerID string, knownOwners []string) (string, error) {
	// Step 1: list ALL NAMESPACE signing keys in the vault (no name filter).
	vaultResp, err := c.vault.ListMyKeys(ctx, &cryptoadminv30.ListMyKeysRequest{
		Filters: &cryptoadminv30.ListKeysFilters{
			Usage: []cryptov30.SigningKeyUsage{cryptov30.SigningKeyUsage_SIGNING_KEY_USAGE_NAMESPACE},
		},
	})
	if err != nil {
		return "", fmt.Errorf("ListMyKeys (NAMESPACE): %w", err)
	}

	type vaultEntry struct {
		name     string
		keyBytes []byte
	}

	var entries []vaultEntry
	for _, km := range vaultResp.GetPrivateKeysMetadata() {
		pkn := km.GetPublicKeyWithName()
		if spk := pkn.GetPublicKey().GetSigningPublicKey(); spk != nil {
			entries = append(entries, vaultEntry{
				name:     pkn.GetName(),
				keyBytes: spk.GetPublicKey(),
			})
		}
	}
	if len(entries) == 0 {
		return "", fmt.Errorf("no NAMESPACE signing keys found in vault")
	}

	// Build a fast lookup set for knownOwners.
	ownerSet := make(map[string]struct{}, len(knownOwners))
	for _, o := range knownOwners {
		ownerSet[o] = struct{}{}
	}

	// Step 2: list NSDs in the synchronizer store.
	nsdResp, err := c.reader.ListNamespaceDelegation(ctx, &topoadminv30.ListNamespaceDelegationRequest{
		BaseQuery: headStateBaseQuery(storeForSync(synchronizerID)),
	})
	if err != nil {
		return "", fmt.Errorf("ListNamespaceDelegation: %w", err)
	}

	// Step 3: find the NSD whose namespace is in knownOwners AND whose
	// target key matches one of the vault entries. Return the key's name.
	for _, r := range nsdResp.GetResults() {
		item := r.GetItem()
		ns := item.GetNamespace()
		if len(ownerSet) > 0 {
			if _, ok := ownerSet[ns]; !ok {
				continue
			}
		}
		for _, entry := range entries {
			if bytes.Equal(item.GetTargetKey().GetPublicKey(), entry.keyBytes) {
				return entry.name, nil
			}
		}
	}

	return "", fmt.Errorf("no NAMESPACE key in vault matches any DNS owner in %v", knownOwners)
}

func (c *GRPCCantonClient) GetProtocolKeyFingerprint(ctx context.Context, knownSigningKeys []string) (string, string, error) {
	// Step 1: list all PROTOCOL signing keys in the vault.
	vaultResp, err := c.vault.ListMyKeys(ctx, &cryptoadminv30.ListMyKeysRequest{
		Filters: &cryptoadminv30.ListKeysFilters{
			Usage: []cryptov30.SigningKeyUsage{cryptov30.SigningKeyUsage_SIGNING_KEY_USAGE_PROTOCOL},
		},
	})
	if err != nil {
		return "", "", fmt.Errorf("ListMyKeys (PROTOCOL): %w", err)
	}

	var vaultKeys [][]byte
	for _, km := range vaultResp.GetPrivateKeysMetadata() {
		if spk := km.GetPublicKeyWithName().GetPublicKey().GetSigningPublicKey(); spk != nil {
			vaultKeys = append(vaultKeys, spk.GetPublicKey())
		}
	}
	if len(vaultKeys) == 0 {
		return "", "", fmt.Errorf("no PROTOCOL signing keys found in vault")
	}

	// Step 2: decode each known signing key and cross-reference with vault keys.
	for _, skB64 := range knownSigningKeys {
		skBytes, decErr := base64.StdEncoding.DecodeString(skB64)
		if decErr != nil {
			continue
		}
		var pk cryptov30.SigningPublicKey
		if unmErr := proto.Unmarshal(skBytes, &pk); unmErr != nil {
			continue
		}
		for _, vaultKeyBytes := range vaultKeys {
			if bytes.Equal(pk.GetPublicKey(), vaultKeyBytes) {
				fp, fpErr := helpers.GetPublicKeyFingerprint(vaultKeyBytes)
				if fpErr != nil {
					return "", "", fmt.Errorf("computing fingerprint: %w", fpErr)
				}

				return fp, skB64, nil
			}
		}
	}

	return "", "", fmt.Errorf("no PROTOCOL signing key in vault matches the party's signing keys")
}

// ── Topology Write ───────────────────────────────────────────────────────────

func (c *GRPCCantonClient) Authorize(
	ctx context.Context,
	serial uint32,
	mapping *protov30.TopologyMapping,
	synchronizerID string,
	mustFullyAuthorize bool,
	signedBy ...string,
) (*protov30.SignedTopologyTransaction, error) {
	input := &topoadminv30.AuthorizeRequest{
		Type: &topoadminv30.AuthorizeRequest_Proposal_{
			Proposal: &topoadminv30.AuthorizeRequest_Proposal{
				Change:  protov30.Enums_TOPOLOGY_CHANGE_OP_ADD_REPLACE,
				Serial:  serial,
				Mapping: mapping,
			},
		},
		MustFullyAuthorize: mustFullyAuthorize,
		Store:              storeForSync(synchronizerID),
		SignedBy:           signedBy,
	}
	resp, err := c.topo.Authorize(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("Authorize: %w", err)
	}

	return resp.GetTransaction(), nil
}

func (c *GRPCCantonClient) SignTransactions(
	ctx context.Context,
	txs []*protov30.SignedTopologyTransaction,
	synchronizerID string,
) ([]*protov30.SignedTopologyTransaction, error) {
	resp, err := c.topo.SignTransactions(ctx, &topoadminv30.SignTransactionsRequest{
		Transactions: txs,
		// signed_by left empty: Canton auto-selects appropriate namespace keys.
		Store: storeForSync(synchronizerID),
	})
	if err != nil {
		return nil, fmt.Errorf("SignTransactions: %w", err)
	}

	return resp.GetTransactions(), nil
}

func (c *GRPCCantonClient) AddTransactions(
	ctx context.Context,
	txs []*protov30.SignedTopologyTransaction,
	synchronizerID string,
) error {
	_, err := c.topo.AddTransactions(ctx, &topoadminv30.AddTransactionsRequest{
		Transactions: txs,
		Store:        storeForSync(synchronizerID),
	})
	if err != nil {
		return fmt.Errorf("AddTransactions: %w", err)
	}

	return nil
}

// ── Topology Read ────────────────────────────────────────────────────────────

func (c *GRPCCantonClient) DNSExists(ctx context.Context, namespace string, synchronizerID string) (bool, error) {
	resp, err := c.reader.ListDecentralizedNamespaceDefinition(ctx,
		&topoadminv30.ListDecentralizedNamespaceDefinitionRequest{
			BaseQuery:       headStateBaseQuery(storeForSync(synchronizerID)),
			FilterNamespace: namespace,
		})
	if err != nil {
		return false, fmt.Errorf("ListDecentralizedNamespaceDefinition: %w", err)
	}
	for _, r := range resp.GetResults() {
		if r.GetItem().GetDecentralizedNamespace() == namespace {
			if r.GetContext().GetOperation() == protov30.Enums_TOPOLOGY_CHANGE_OP_ADD_REPLACE {
				return true, nil
			}
		}
	}

	return false, nil
}

func (c *GRPCCantonClient) NSDExists(ctx context.Context, namespace string, synchronizerID string) (bool, error) {
	resp, err := c.reader.ListNamespaceDelegation(ctx, &topoadminv30.ListNamespaceDelegationRequest{
		BaseQuery:       headStateBaseQuery(storeForSync(synchronizerID)),
		FilterNamespace: namespace,
	})
	if err != nil {
		return false, fmt.Errorf("ListNamespaceDelegation: %w", err)
	}
	for _, r := range resp.GetResults() {
		if r.GetItem().GetNamespace() == namespace {
			if r.GetContext().GetOperation() == protov30.Enums_TOPOLOGY_CHANGE_OP_ADD_REPLACE {
				return true, nil
			}
		}
	}

	return false, nil
}

func (c *GRPCCantonClient) P2PExists(ctx context.Context, partyUID string, synchronizerID string) (bool, error) {
	resp, err := c.reader.ListPartyToParticipant(ctx,
		&topoadminv30.ListPartyToParticipantRequest{
			BaseQuery:   headStateBaseQuery(storeForSync(synchronizerID)),
			FilterParty: partyUID,
		})
	if err != nil {
		return false, fmt.Errorf("ListPartyToParticipant: %w", err)
	}
	for _, r := range resp.GetResults() {
		if strings.Contains(r.GetItem().GetParty(), partyUID) {
			if r.GetContext().GetOperation() == protov30.Enums_TOPOLOGY_CHANGE_OP_ADD_REPLACE {
				return true, nil
			}
		}
	}

	return false, nil
}

func (c *GRPCCantonClient) GetDNS(ctx context.Context, namespace string, synchronizerID string) (*DNSState, error) {
	resp, err := c.reader.ListDecentralizedNamespaceDefinition(ctx,
		&topoadminv30.ListDecentralizedNamespaceDefinitionRequest{
			BaseQuery:       headStateBaseQuery(storeForSync(synchronizerID)),
			FilterNamespace: namespace,
		})
	if err != nil {
		return nil, fmt.Errorf("ListDecentralizedNamespaceDefinition: %w", err)
	}

	for _, r := range resp.GetResults() {
		if r.GetItem().GetDecentralizedNamespace() == namespace &&
			r.GetContext().GetOperation() == protov30.Enums_TOPOLOGY_CHANGE_OP_ADD_REPLACE {
			return &DNSState{
				DecentralizedNamespace: r.GetItem().GetDecentralizedNamespace(),
				Owners:                 r.GetItem().GetOwners(),
				Threshold:              r.GetItem().GetThreshold(),
				Serial:                 r.GetContext().GetSerial(),
			}, nil
		}
	}

	return nil, fmt.Errorf("DecentralizedNamespaceDefinition not found for namespace %q", namespace)
}

func (c *GRPCCantonClient) GetP2P(ctx context.Context, partyUID string, synchronizerID string) (*P2PState, error) {
	resp, err := c.reader.ListPartyToParticipant(ctx,
		&topoadminv30.ListPartyToParticipantRequest{
			BaseQuery:   headStateBaseQuery(storeForSync(synchronizerID)),
			FilterParty: partyUID,
		})
	if err != nil {
		return nil, fmt.Errorf("ListPartyToParticipant: %w", err)
	}

	for _, r := range resp.GetResults() {
		if strings.Contains(r.GetItem().GetParty(), partyUID) &&
			r.GetContext().GetOperation() == protov30.Enums_TOPOLOGY_CHANGE_OP_ADD_REPLACE {
			hps := r.GetItem().GetParticipants()
			participants := make([]P2PParticipantInfo, len(hps))
			for i, hp := range hps {
				participants[i] = P2PParticipantInfo{
					ParticipantUID: hp.GetParticipantUid(),
					Permission:     hp.GetPermission().String(),
				}
			}

			var signingKeys *P2PSigningKeysInfo
			if skwt := r.GetItem().GetPartySigningKeys(); skwt != nil && len(skwt.GetKeys()) > 0 {
				keys := make([]string, len(skwt.GetKeys()))
				for j, k := range skwt.GetKeys() {
					kb, mErr := proto.Marshal(k)
					if mErr != nil {
						return nil, fmt.Errorf("marshalling party signing key %d: %w", j, mErr)
					}
					keys[j] = base64.StdEncoding.EncodeToString(kb)
				}
				signingKeys = &P2PSigningKeysInfo{
					Keys:      keys,
					Threshold: skwt.GetThreshold(),
				}
			}

			return &P2PState{
				Party:            r.GetItem().GetParty(),
				Participants:     participants,
				Threshold:        r.GetItem().GetThreshold(),
				Serial:           r.GetContext().GetSerial(),
				PartySigningKeys: signingKeys,
			}, nil
		}
	}

	return nil, fmt.Errorf("PartyToParticipant mapping not found for party %q", partyUID)
}

func (c *GRPCCantonClient) ListDecentralizedNamespaces(ctx context.Context, synchronizerID string) ([]*DNSState, error) {
	resp, err := c.reader.ListDecentralizedNamespaceDefinition(ctx,
		&topoadminv30.ListDecentralizedNamespaceDefinitionRequest{
			BaseQuery: headStateBaseQuery(storeForSync(synchronizerID)),
		})
	if err != nil {
		return nil, fmt.Errorf("ListDecentralizedNamespaceDefinition: %w", err)
	}

	var results []*DNSState
	for _, r := range resp.GetResults() {
		if r.GetContext().GetOperation() == protov30.Enums_TOPOLOGY_CHANGE_OP_ADD_REPLACE {
			results = append(results, &DNSState{
				DecentralizedNamespace: r.GetItem().GetDecentralizedNamespace(),
				Owners:                 r.GetItem().GetOwners(),
				Threshold:              r.GetItem().GetThreshold(),
				Serial:                 r.GetContext().GetSerial(),
			})
		}
	}

	return results, nil
}

// ── Package Management ───────────────────────────────────────────────────────

func (c *GRPCCantonClient) UploadDar(ctx context.Context, darBytes []byte) (string, error) {
	resp, err := c.pkg.UploadDar(ctx, &participantv30.UploadDarRequest{
		Dars: []*participantv30.UploadDarRequest_UploadDarData{
			{Bytes: darBytes},
		},
		VetAllPackages:     true,
		SynchronizeVetting: true,
	})
	if err != nil {
		return "", fmt.Errorf("UploadDar: %w", err)
	}

	darIDs := resp.GetDarIds()
	if len(darIDs) == 0 {
		return "", fmt.Errorf("UploadDar: no DAR IDs returned")
	}

	return darIDs[0], nil
}
