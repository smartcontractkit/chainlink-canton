package cmd

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/admin"
	v30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/protocol/v30"
	versionv1 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/version/v1"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-canton/examples/canton-ccip-cli/internal/cantonops"
	"github.com/smartcontractkit/chainlink-canton/examples/canton-ccip-cli/ledger/usbwallet"
)

func newCantonExternalPartyCommand(g *Globals) *cobra.Command {
	c := &cobra.Command{
		Use:   "external-party",
		Short: "External Party management",
	}
	c.AddCommand(newGetFingerprintCmd(g))
	c.AddCommand(newPrepareTopologyTransactionCmd(g))
	c.AddCommand(newSignTopologyTransactionCmd(g))
	c.AddCommand(newSubmitTopologyTransactionCmd(g))

	return c
}

func newGetFingerprintCmd(g *Globals) *cobra.Command {
	var (
		derivationPathFlag string
	)
	c := &cobra.Command{
		Use:   "get-fingerprint",
		Short: "Get the fingerprint of a Ledger account",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			_, err := g.Resolve(ctx, true)
			if err != nil {
				return err
			}

			derivationPath, err := cantonops.ParseDerivationPath(derivationPathFlag)
			if err != nil {
				return fmt.Errorf("failed to parse derivation path: %w", err)
			}
			fmt.Printf("Getting fingerprint for derivation path: %v\n", derivationPath.String())

			fmt.Println("Looking for connected Ledger devices...")
			ledgerHub, err := usbwallet.NewCantonLedgerHub()
			if err != nil {
				return fmt.Errorf("failed to create Ledger hub: %w", err)
			}
			wallets := ledgerHub.Wallets()
			if len(wallets) == 0 {
				return fmt.Errorf("no Ledger wallets found. Please connect a Ledger device and try again")
			}
			fmt.Printf("Found %d connected Ledger wallets\n", len(wallets))
			wallet := wallets[0]

			err = wallet.Open("")
			if err != nil {
				return fmt.Errorf("failed to open Ledger wallet: %w", err)
			}
			defer wallet.Close()

			// Get the public key twice. First without confirmation to be able to display it and then prompt the user to confirm it matches what's shown on the Ledger
			// Due to how the Canton Ledger app is implemented, it'll show the party as "ldg::{fingerprint}"

			pubkey, err := wallet.GetPublicKey(derivationPath, false)
			if err != nil {
				return fmt.Errorf("failed to get public key from Ledger: %w", err)
			}
			fingerprint, err := usbwallet.Fingerprint(pubkey)
			if err != nil {
				return fmt.Errorf("failed to compute fingerprint: %w", err)
			}

			fmt.Printf("Public Key: %v\n", hex.EncodeToString(pubkey))
			fmt.Printf("Confirm this matches the party shown on your Ledger: ldg::%v\n", fingerprint)

			_, err = wallet.GetPublicKey(derivationPath, true)
			if err != nil {
				return fmt.Errorf("failed to get public key from Ledger: %w", err)
			}
			publicKeyDER, err := x509.MarshalPKIXPublicKey(ed25519.PublicKey(pubkey))
			if err != nil {
				return fmt.Errorf("failed to marshal public key to DER: %w", err)
			}

			fmt.Println("Public key confirmed.")
			fmt.Printf("Share the public key with your operator to continue: %v\n", hex.EncodeToString(publicKeyDER))

			return nil
		},
	}
	c.Flags().StringVar(&derivationPathFlag, "derivation-path", "m/44'/6767'/0'/0'/0'", "Derivation path for the Ledger account or derivation index. Defaults to m/44'/6767'/0'/0'/0'. If only a number is specified, will increment last component, e.g. 42=m/44'/6767'/0'/0'/42'")

	return c
}

type topologyTransaction struct {
	PartyId              string   `json:"partyId"`
	PublicKeyFingerprint string   `json:"publicKeyFingerprint"`
	MultiHash            string   `json:"multiHash"`
	TopologyTransactions []string `json:"topologyTransactions"`
	Signature            string   `json:"signature,omitempty"`
}

func newPrepareTopologyTransactionCmd(g *Globals) *cobra.Command {
	var (
		synchronizerId string
	)
	c := &cobra.Command{
		Use:   "prepare-topology",
		Short: "Prepare the topology transaction for a given party ID and print to output",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 || args[0] == "" || args[1] == "" {
				return fmt.Errorf("requires two arguments: party hint and DER public key")
			}
			partyHint := args[0]
			pubKeyHex := args[1]
			pubKeyDER, err := hex.DecodeString(strings.TrimPrefix(pubKeyHex, "0x"))
			if err != nil {
				return fmt.Errorf("failed to decode public key hex: %w", err)
			}
			_, err = x509.ParsePKIXPublicKey(pubKeyDER)
			if err != nil {
				return fmt.Errorf("failed to parse public key from DER: %w", err)
			}

			ctx := cmd.Context()
			b, err := g.Resolve(ctx, false)
			if err != nil {
				return err
			}

			connectedSynchronizersResp, err := b.Participant.LedgerServices.State.GetConnectedSynchronizers(ctx, &apiv2.GetConnectedSynchronizersRequest{})
			if err != nil {
				return fmt.Errorf("failed to get connected synchronizers: %w", err)
			}
			synchronizerId := ""
			for _, synchronizer := range connectedSynchronizersResp.GetConnectedSynchronizers() {
				if synchronizer.GetSynchronizerAlias() == "global" {
					synchronizerId = synchronizer.GetSynchronizerId()
					break
				}
			}
			if synchronizerId == "" {
				return fmt.Errorf("no synchronizer with alias 'global' found among connected synchronizers, use --synchronizer-id to manually specify a synchronizer ID")
			}

			generateResponse, err := b.Participant.LedgerServices.Admin.PartyManagement.GenerateExternalPartyTopology(ctx, &admin.GenerateExternalPartyTopologyRequest{
				Synchronizer: synchronizerId,
				PartyHint:    partyHint,
				PublicKey: &apiv2.SigningPublicKey{
					Format:  apiv2.CryptoKeyFormat_CRYPTO_KEY_FORMAT_DER_X509_SUBJECT_PUBLIC_KEY_INFO,
					KeyData: pubKeyDER,
					KeySpec: apiv2.SigningKeySpec_SIGNING_KEY_SPEC_EC_CURVE25519,
				},
			})
			if err != nil {
				return fmt.Errorf("failed to generate external party topology: %w", err)
			}

			tx := topologyTransaction{
				PartyId:              generateResponse.GetPartyId(),
				PublicKeyFingerprint: generateResponse.GetPublicKeyFingerprint(),
				MultiHash:            hex.EncodeToString(generateResponse.GetMultiHash()),
				TopologyTransactions: make([]string, len(generateResponse.GetTopologyTransactions())),
				Signature:            "",
			}
			for i, b := range generateResponse.GetTopologyTransactions() {
				tx.TopologyTransactions[i] = hex.EncodeToString(b)
			}

			data, err := json.Marshal(tx)
			if err != nil {
				return fmt.Errorf("failed to marshal topology transaction to JSON: %w", err)
			}
			fmt.Println("Topology transaction prepared:")
			fmt.Println(string(data))

			return nil
		},
	}
	c.Flags().StringVar(&synchronizerId, "synchronizer-id", "", "Synchronizer ID to use for the topology transaction. Defaults to the global synchronizer.")

	return c
}

func newSignTopologyTransactionCmd(g *Globals) *cobra.Command {
	var (
		derivationPathFlag string
	)
	c := &cobra.Command{
		Use:   "sign-topology",
		Short: "Sign the topology transaction for a given party",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 || args[0] == "" {
				return fmt.Errorf("requires one argument: the topology transaction JSON")
			}
			topologyJSON := args[0]
			var tx topologyTransaction
			if err := json.Unmarshal([]byte(topologyJSON), &tx); err != nil {
				return fmt.Errorf("failed to unmarshal topology transaction JSON: %w", err)
			}

			ctx := cmd.Context()
			_, err := g.Resolve(ctx, true)
			if err != nil {
				return err
			}

			// Parse prepared transaction
			topologyTransactions := make([][]byte, len(tx.TopologyTransactions))
			var publicKeys [][]byte
			for i, transaction := range tx.TopologyTransactions {
				fmt.Printf("Parsing topology transaction %d...\n", i)
				b, err := hex.DecodeString(strings.TrimPrefix(transaction, "0x"))
				if err != nil {
					return fmt.Errorf("failed to decode topology transaction %d: %w", i, err)
				}
				topologyTransactions[i] = b
				versionedMessage := versionv1.UntypedVersionedMessage{}
				err = proto.Unmarshal(b, &versionedMessage)
				if err != nil {
					return fmt.Errorf("failed to unmarshal versioned message: %w", err)
				}
				topologyTransaction := v30.TopologyTransaction{}
				err = proto.Unmarshal(versionedMessage.GetData(), &topologyTransaction)
				if err != nil {
					return fmt.Errorf("failed to unmarshal topology transaction: %w", err)
				}

				tw := table.NewWriter()
				tw.SetStyle(table.StyleLight)
				tw.Style().Options.SeparateRows = true

				tw.AppendRow(table.Row{"Operation", topologyTransaction.GetOperation().String(), topologyTransaction.GetOperation().String()}, table.RowConfig{AutoMerge: true})
				switch m := topologyTransaction.GetMapping().GetMapping().(type) {
				case *v30.TopologyMapping_PartyToParticipant:
					tw.AppendRow(table.Row{"Mapping Type", "PartyToParticipant", "PartyToParticipant"}, table.RowConfig{AutoMerge: true})
					tw.AppendRow(table.Row{"Party", m.PartyToParticipant.GetParty(), m.PartyToParticipant.GetParty()}, table.RowConfig{AutoMerge: true})
					tw.AppendRow(table.Row{"Confirmation Threshold", m.PartyToParticipant.GetThreshold(), m.PartyToParticipant.GetThreshold()}, table.RowConfig{AutoMerge: true})
					tw.AppendRow(table.Row{"Participants", "Participants", "Participants"}, table.RowConfig{AutoMerge: true})
					for i2, participant := range m.PartyToParticipant.GetParticipants() {
						tw.AppendRow(table.Row{strconv.Itoa(i2), "ID", participant.GetParticipantUid()})
						tw.AppendRow(table.Row{strconv.Itoa(i2), "Permission", participant.GetPermission().String()})
					}
					tw.AppendRow(table.Row{"Signing Keys", "Signing Keys", "Signing Keys"}, table.RowConfig{AutoMerge: true})
					tw.AppendRow(table.Row{"Signing Threshold", m.PartyToParticipant.GetPartySigningKeys().GetThreshold(), m.PartyToParticipant.GetPartySigningKeys().GetThreshold()}, table.RowConfig{AutoMerge: true})
					for i2, key := range m.PartyToParticipant.GetPartySigningKeys().GetKeys() {
						tw.AppendRow(table.Row{strconv.Itoa(i2), "Public Key", hex.EncodeToString(key.GetPublicKey())})
						publicKeys = append(publicKeys, key.GetPublicKey())
						tw.AppendRow(table.Row{strconv.Itoa(i2), "Key Spec", key.GetKeySpec().String()})
						tw.AppendRow(table.Row{strconv.Itoa(i2), "Key Spec", key.GetKeySpec().String()})
						usages := make([]string, len(key.GetUsage()))
						for i3, usage := range key.GetUsage() {
							usages[i3] = usage.String()
						}
						tw.AppendRow(table.Row{strconv.Itoa(i2), "Usage", strings.Join(usages, ", ")})
					}
				default:
					return fmt.Errorf("unsupported mapping type: %T", m)
				}
				tw.SetColumnConfigs([]table.ColumnConfig{
					{Number: 1, AutoMerge: true, Align: text.AlignLeft},
					{Number: 2, AutoMerge: true, Align: text.AlignLeft},
					{Number: 3, AutoMerge: true, Align: text.AlignLeft},
				})
				fmt.Println(tw.Render())
			}

			// Calculate Multihash
			multihash := topologyTransactionMultihash(topologyTransactions)

			// Connect Ledger
			fmt.Println("Looking for connected Ledger devices...")
			ledgerHub, err := usbwallet.NewCantonLedgerHub()
			if err != nil {
				return fmt.Errorf("failed to create Ledger hub: %w", err)
			}
			wallets := ledgerHub.Wallets()
			if len(wallets) == 0 {
				return fmt.Errorf("no Ledger wallets found. Please connect a Ledger device and try again")
			}
			fmt.Printf("Found %d connected Ledger wallets\n", len(wallets))
			wallet := wallets[0]

			err = wallet.Open("")
			if err != nil {
				return fmt.Errorf("failed to open Ledger wallet: %w", err)
			}
			defer wallet.Close()

			// Parse derivation path
			derivationPath, err := cantonops.ParseDerivationPath(derivationPathFlag)
			if err != nil {
				return fmt.Errorf("failed to parse derivation path: %w", err)
			}
			fmt.Printf("Getting public key for derivation path: %v\n", derivationPath.String())

			// Get the public key at the derivation path and compare with expected key
			pubkey, err := wallet.GetPublicKey(derivationPath, false)
			if err != nil {
				return fmt.Errorf("failed to get public key from Ledger: %w", err)
			}
			publicKeyDER, err := x509.MarshalPKIXPublicKey(ed25519.PublicKey(pubkey))
			if err != nil {
				return fmt.Errorf("failed to marshal public key to DER: %w", err)
			}
			fmt.Printf("Using public key %s to sign transaction\n", hex.EncodeToString(publicKeyDER))

			if !slices.ContainsFunc(publicKeys, func(i []byte) bool {
				return bytes.Equal(i, publicKeyDER)
			}) {
				return fmt.Errorf("public key from Ledger does not match any of the expected signing keys in the topology transaction")
			}

			// Sign hash
			fmt.Printf("📜 Confirm the signature on your Ledger device for the multihash: %s\n", strings.ToUpper(hex.EncodeToString(multihash)))
			signature, err := wallet.SignHash(derivationPath, multihash)
			if err != nil {
				return fmt.Errorf("failed to sign multihash with Ledger: %w", err)
			}

			tx.Signature = hex.EncodeToString(signature)

			data, err := json.Marshal(tx)
			if err != nil {
				return fmt.Errorf("failed to marshal signed topology transaction to JSON: %w", err)
			}
			fmt.Println("Topology transaction signed:")
			fmt.Println(string(data))

			return nil
		},
	}
	c.Flags().StringVar(&derivationPathFlag, "derivation-path", "m/44'/6767'/0'/0'/0'", "Derivation path for the Ledger account or derivation index. Defaults to m/44'/6767'/0'/0'/0'. If only a number is specified, will increment last component, e.g. 42=m/44'/6767'/0'/0'/42'")

	return c
}

// hashPurposeTopologyTransactionSignature is Canton's hash purpose for topology transaction signatures.
const hashPurposeTopologyTransactionSignature uint32 = 11

// hashPurposeMultiTopologyTransaction is Canton's hash purpose for combined topology transaction signatures.
const hashPurposeMultiTopologyTransaction uint32 = 55

func cantonSHA256Multihash(purpose uint32, content []byte) []byte {
	var hashPurpose [4]byte
	binary.BigEndian.PutUint32(hashPurpose[:], purpose)

	digest := sha256.New()
	digest.Write(hashPurpose[:])
	digest.Write(content)

	// 0x12 0x20 is the multihash prefix for a 32 byte SHA-256 digest.
	return append([]byte{0x12, 0x20}, digest.Sum(nil)...)
}

func topologyTransactionMultihash(topologyTransactions [][]byte) []byte {
	transactionHashes := make([][]byte, 0, len(topologyTransactions))
	for _, transaction := range topologyTransactions {
		transactionHashes = append(transactionHashes, cantonSHA256Multihash(hashPurposeTopologyTransactionSignature, transaction))
	}

	sort.Slice(transactionHashes, func(i, j int) bool {
		return bytes.Compare(transactionHashes[i], transactionHashes[j]) < 0
	})

	encoded := make([]byte, 4)
	binary.BigEndian.PutUint32(encoded, uint32(len(transactionHashes))) //nolint:gosec
	for _, transactionHash := range transactionHashes {
		var hashLen [4]byte
		binary.BigEndian.PutUint32(hashLen[:], uint32(len(transactionHash))) //nolint:gosec
		encoded = append(encoded, hashLen[:]...)                             //nolint:makezero // first four bytes are preallocated
		encoded = append(encoded, transactionHash...)                        //nolint:makezero // first four bytes are preallocated
	}

	return cantonSHA256Multihash(hashPurposeMultiTopologyTransaction, encoded)
}

func newSubmitTopologyTransactionCmd(g *Globals) *cobra.Command {
	var (
		synchronizerId string
	)
	c := &cobra.Command{
		Use:   "submit-topology",
		Short: "Submit the signed topology transaction for a given party ID and execute it",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 || args[0] == "" {
				return fmt.Errorf("requires one argument: the topology transaction JSON")
			}
			topologyJSON := args[0]
			var tx topologyTransaction
			if err := json.Unmarshal([]byte(topologyJSON), &tx); err != nil {
				return fmt.Errorf("failed to unmarshal topology transaction JSON: %w", err)
			}

			ctx := cmd.Context()
			b, err := g.Resolve(ctx, false)
			if err != nil {
				return err
			}

			connectedSynchronizersResp, err := b.Participant.LedgerServices.State.GetConnectedSynchronizers(ctx, &apiv2.GetConnectedSynchronizersRequest{})
			if err != nil {
				return fmt.Errorf("failed to get connected synchronizers: %w", err)
			}
			synchronizerId := ""
			for _, synchronizer := range connectedSynchronizersResp.GetConnectedSynchronizers() {
				if synchronizer.GetSynchronizerAlias() == "global" {
					synchronizerId = synchronizer.GetSynchronizerId()
					break
				}
			}
			if synchronizerId == "" {
				return fmt.Errorf("no synchronizer with alias 'global' found among connected synchronizers, use --synchronizer-id to manually specify a synchronizer ID")
			}

			// Parse input
			onboardingTransactions := make([]*admin.AllocateExternalPartyRequest_SignedTransaction, len(tx.TopologyTransactions))
			for i, transaction := range tx.TopologyTransactions {
				txBytes, err := hex.DecodeString(strings.TrimPrefix(transaction, "0x"))
				if err != nil {
					return fmt.Errorf("failed to decode topology transaction %d: %w", i, err)
				}
				onboardingTransactions[i] = &admin.AllocateExternalPartyRequest_SignedTransaction{
					Transaction: txBytes,
				}
			}
			signature, err := hex.DecodeString(strings.TrimPrefix(tx.Signature, "0x"))
			if err != nil {
				return fmt.Errorf("failed to decode signature: %w", err)
			}

			// Submit transaction
			generateResponse, err := b.Participant.LedgerServices.Admin.PartyManagement.AllocateExternalParty(ctx, &admin.AllocateExternalPartyRequest{
				Synchronizer:           synchronizerId,
				OnboardingTransactions: onboardingTransactions,
				MultiHashSignatures: []*apiv2.Signature{{
					Format:               apiv2.SignatureFormat_SIGNATURE_FORMAT_CONCAT,
					Signature:            signature,
					SignedBy:             tx.PublicKeyFingerprint,
					SigningAlgorithmSpec: apiv2.SigningAlgorithmSpec_SIGNING_ALGORITHM_SPEC_ED25519,
				}},
			})
			if err != nil {
				return fmt.Errorf("failed to generate external party topology: %w", err)
			}

			// Verify the party has been created
			resp, err := b.Participant.LedgerServices.Admin.PartyManagement.ListKnownParties(ctx, &admin.ListKnownPartiesRequest{
				FilterParty: generateResponse.GetPartyId(),
			})
			if err != nil {
				return fmt.Errorf("failed to list known parties: %w", err)
			}
			if len(resp.GetPartyDetails()) != 1 {
				return fmt.Errorf("expected 1 known party for %s, got %d", generateResponse.GetPartyId(), len(resp.GetPartyDetails()))
			}
			fmt.Printf("Party %s created successfully.\n", generateResponse.GetPartyId())

			return nil
		},
	}
	c.Flags().StringVar(&synchronizerId, "synchronizer-id", "", "Synchronizer ID to use for the topology transaction. Defaults to the global synchronizer.")

	return c
}
