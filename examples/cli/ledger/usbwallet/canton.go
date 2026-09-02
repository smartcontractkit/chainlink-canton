package usbwallet

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"sync"

	interactivepb "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2/interactive"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
	"github.com/karalabe/hid"
)

const CLA byte = 0xe0

// maxAPDUDataLength is the maximum number of payload bytes that fit into a single APDU,
// since Lc is a single byte.
const maxAPDUDataLength = 255

type cantonInstruction byte

const (
	INS_GET_VERSION    cantonInstruction = 0x03
	INS_GET_APP_NAME   cantonInstruction = 0x04
	INS_GET_PUBLIC_KEY cantonInstruction = 0x05
	INS_SIGN_TX        cantonInstruction = 0x06
)

type cantonParam1 byte

const (
	P1_NONE                           cantonParam1 = 0x00
	P1_CONFIRM                        cantonParam1 = 0x01
	P1_SIGN_HASH                      cantonParam1 = 0x00
	P1_SIGN_UNTYPED_VERSIONED_MESSAGE cantonParam1 = 0x01
	P1_SIGN_PREPARED_TRANSACTION      cantonParam1 = 0x02
)

type cantonParam2 byte

const (
	P2_NONE        cantonParam2 = 0x00
	P2_FIRST       cantonParam2 = 0x01
	P2_MORE        cantonParam2 = 0x02
	P2_MESSAGE_END cantonParam2 = 0x04
)

// challengeLength is the size of the optional attestation challenge
// (16 byte nonce + 8 byte timestamp).
const challengeLength = 24

const swOK uint16 = 0x9000

// LedgerStatusError is returned when the device replies with a status word other than 0x9000.
type LedgerStatusError uint16

func (e LedgerStatusError) Error() string {
	if desc, ok := cantonStatusWords[uint16(e)]; ok {
		return fmt.Sprintf("ledger: device returned status word 0x%04X (%s)", uint16(e), desc)
	}

	return fmt.Sprintf("ledger: device returned status word 0x%04X", uint16(e))
}

// cantonStatusWords maps status words of the Canton Ledger app to human-readable descriptions.
// See https://github.com/LedgerHQ/app-canton/blob/develop/doc/APDU.md#status-words
var cantonStatusWords = map[uint16]string{
	0x6985: "denied by user",
	0x6A80: "incorrect data; when signing a raw hash this means blind signing is disabled in the Canton app settings",
	0x6A86: "wrong P1/P2",
	0x6A87: "wrong data length",
	0x6D00: "invalid instruction",
	0x6E00: "invalid CLA",
	0xB000: "wrong response length",
	0xB001: "display BIP32 path failed",
	0xB002: "display address failed",
	0xB003: "display amount failed",
	0xB004: "wrong transaction length",
	0xB005: "transaction parsing failed",
	0xB006: "transaction hashing failed",
	0xB007: "bad state",
	0xB008: "signature failed",
	0xB009: "challenge signature failed",
}

// errLedgerReplyInvalidHeader is the error message returned by a Ledger data exchange
// if the device replies with a mismatching header.
var errLedgerReplyInvalidHeader = errors.New("ledger: invalid reply header")

// errLedgerInvalidVersionReply is the error message returned by a Ledger version retrieval
// when a response does arrive, but it does not contain the expected data.
var errLedgerInvalidVersionReply = errors.New("ledger: invalid version reply")

var errLedgerInvalidAppNameReply = errors.New("ledger: invalid app name reply")

// hashPurposePublicKeyFingerprint is Canton's hash purpose for key fingerprints.
const hashPurposePublicKeyFingerprint uint32 = 12

// Fingerprint returns the Canton fingerprint of a raw Ed25519 public key, i.e. the hex encoded
// SHA-256 multihash of the hash purpose concatenated with the raw key.
func Fingerprint(publicKey []byte) (string, error) {
	if len(publicKey) != ed25519.PublicKeySize {
		return "", fmt.Errorf("ledger: invalid public key length %d, want %d", len(publicKey), ed25519.PublicKeySize)
	}

	var purpose [4]byte
	binary.BigEndian.PutUint32(purpose[:], hashPurposePublicKeyFingerprint)

	digest := sha256.New()
	digest.Write(purpose[:])
	digest.Write(publicKey)

	// 0x12 0x20 is the multihash prefix for a 32 byte SHA-256 digest.
	return hex.EncodeToString(append([]byte{0x12, 0x20}, digest.Sum(nil)...)), nil
}

type CantonWallet interface {
	URL() accounts.URL
	Status() (string, error)
	Open(passphrase string) error
	Close() error

	GetVersion() ([3]byte, error)
	GetAppName() (string, error)
	GetPublicKey(path accounts.DerivationPath, display bool) ([]byte, error)
	SignHash(path accounts.DerivationPath, hash []byte) ([]byte, error)
	SignPreparedTransaction(path accounts.DerivationPath, transaction *interactivepb.PreparedTransaction) ([]byte, error)
	SignTopologyTransactions(path accounts.DerivationPath, transactions [][]byte, challenge []byte) (signature []byte, challengeSignature []byte, err error)
}

var _ CantonWallet = &cantonWallet{}

type cantonWallet struct {
	log log.Logger

	hub    *Hub
	url    *accounts.URL  // Textual URL uniquely identifying this wallet
	info   hid.DeviceInfo // Known USB device infos about the wallet
	device hid.Device     // USB device advertising itself as a hardware wallet
	mux    sync.Mutex
}

func (w *cantonWallet) URL() accounts.URL {
	return *w.url
}

func (w *cantonWallet) Status() (string, error) {
	w.mux.Lock()
	defer w.mux.Unlock()

	if w.device == nil {
		return "Closed", accounts.ErrWalletClosed
	}

	// Get App version
	version, err := w.ledgerVersion()
	if err != nil {
		return "Unable to detect Canton app version", err
	}

	return fmt.Sprintf("Canton app v%v.%v.%v online", version[0], version[1], version[2]), nil
}

func (w *cantonWallet) Open(passphrase string) error {
	w.mux.Lock()
	defer w.mux.Unlock()

	// Make sure the actual device connection is done only once
	if w.device != nil {
		return accounts.ErrWalletAlreadyOpen
	}
	device, err := w.info.Open()
	if err != nil {
		return err
	}
	w.device = device

	// Get App version
	version, err := w.ledgerVersion()
	if err != nil {
		return err
	}

	// Fail if no version is returned
	if version == [3]byte{0, 0, 0} {
		_ = device.Close()
		w.device = nil
		return fmt.Errorf("could not detect Canton App version, make sure the App is open")
	}

	return nil
}

func (w *cantonWallet) Close() error {
	w.mux.Lock()
	defer w.mux.Unlock()

	// Don't fail if device isn't opened
	if w.device != nil {
		if err := w.device.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (w *cantonWallet) ledgerVersion() ([3]byte, error) {
	reply, err := w.ledgerExchange(INS_GET_VERSION, P1_NONE, P2_NONE, nil)
	if err != nil {
		return [3]byte{}, err
	}
	if len(reply) != 3 {
		return [3]byte{}, errLedgerInvalidVersionReply
	}
	var version [3]byte
	copy(version[:], reply[:])

	return version, nil
}

func (w *cantonWallet) GetVersion() ([3]byte, error) {
	w.mux.Lock()
	defer w.mux.Unlock()

	if w.device == nil {
		return [3]byte{}, accounts.ErrWalletClosed
	}
	version, err := w.ledgerVersion()

	return version, err
}

func (w *cantonWallet) GetAppName() (string, error) {
	w.mux.Lock()
	defer w.mux.Unlock()

	if w.device == nil {
		return "", accounts.ErrWalletClosed
	}

	reply, err := w.ledgerExchange(INS_GET_APP_NAME, P1_NONE, P2_NONE, nil)
	if err != nil {
		return "", err
	}
	if len(reply) == 0 {
		return "", errLedgerInvalidAppNameReply
	}

	return string(reply), nil
}

func (w *cantonWallet) GetPublicKey(path accounts.DerivationPath, display bool) ([]byte, error) {
	w.mux.Lock()
	defer w.mux.Unlock()

	if w.device == nil {
		return nil, accounts.ErrWalletClosed
	}

	// If display==true, Ledger will prompt the user to confirm
	// Note that this doesn't really work on Canton as the ledger will show the PartyId that the
	// key would resolve to with Ledger Live's Canton offering, which is using their validators.
	p1 := P1_NONE
	if display {
		p1 = P1_CONFIRM
	}
	reply, err := w.ledgerExchange(INS_GET_PUBLIC_KEY, p1, P2_NONE, encodeDerivationPath(path))
	if err != nil {
		return nil, err
	}
	// Reply layout: pub_key_len (1) | pub_key (32) | chain_code_len (1) | chain_code (32)
	if len(reply) < 1 || len(reply) < 1+int(reply[0]) {
		return nil, errors.New("ledger: reply lacks public key entry")
	}
	pubKey := reply[1 : 1+int(reply[0])]
	if len(pubKey) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("ledger: unexpected public key length %d, want %d", len(pubKey), ed25519.PublicKeySize)
	}

	return pubKey, nil
}

// SignHash signs a raw 32 byte hash or a 34 byte multi-hash (2 prefix bytes + SHA-256).
//
// The device cannot display anything meaningful for a bare hash, so blind signing has to be
// enabled in the Canton app settings, otherwise the device rejects the request with 0x6A80.
// Prefer SignTopologyTransactions for party onboarding, which lets the device derive and
// display what it signs.
func (w *cantonWallet) SignHash(path accounts.DerivationPath, hash []byte) ([]byte, error) {
	if len(hash) != 32 && len(hash) != 34 {
		return nil, fmt.Errorf("ledger: invalid hash length %d, want 32 or 34", len(hash))
	}

	w.mux.Lock()
	defer w.mux.Unlock()
	if w.device == nil {
		return nil, accounts.ErrWalletClosed
	}

	signature, _, err := w.ledgerSign(P1_SIGN_HASH, path, [][]byte{hash}, nil)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

// SignPreparedTransaction signs a prepared transaction obtained from PrepareSubmission.
//
// The transaction is split into the components the device expects (see
// https://github.com/LedgerHQ/app-canton/blob/develop/doc/SPLIT_TRANSACTION.md) and streamed to
// the device, which recomputes the transaction hash itself instead of trusting a host supplied
// one. The device clear-signs the transaction if it recognizes the template, and otherwise falls
// back to blind signing.
// Currently supported transaction types are:
//
//	Splice.Api.Token.TransferInstructionV1.TransferFactory_Transfer
//	Splice.ExternalPartyAmuletRules.ExternalPartyAmuletRules_CreateTransferCommand
//	Splice.Wallet.TransferPreapproval.TransferPreapprovalProposal
//
// See: https://github.com/LedgerHQ/app-canton/blob/develop/doc/CLEARSIGN.md#supported-transaction-types-lines-113-123
func (w *cantonWallet) SignPreparedTransaction(path accounts.DerivationPath, transaction *interactivepb.PreparedTransaction) ([]byte, error) {
	messages, err := splitPreparedTransaction(transaction)
	if err != nil {
		return nil, err
	}

	w.mux.Lock()
	defer w.mux.Unlock()
	if w.device == nil {
		return nil, accounts.ErrWalletClosed
	}

	signature, _, err := w.ledgerSign(P1_SIGN_PREPARED_TRANSACTION, path, messages, nil)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

// SignTopologyTransactions signs the onboarding topology transactions returned by
// GenerateExternalPartyTopology. The device hashes every transaction itself and derives the
// multi-hash from them, so the returned signature is the multi-hash signature.
//
// The optional challenge (16 byte nonce + 8 byte timestamp) additionally yields an attestation
// signature over multi-hash || challenge.
//
// The app only clear-signs onboardings onto the Canton Network validators it ships an allow list
// for: it requires a namespace delegation, a party to key mapping and a party to participant
// mapping hosting the party on exactly two known Ledger/Kiln validators. Onboarding onto any
// other participant is rejected with an SW_TOPOLOGY_* status word, so for local or private
// synchronizers use SignHash on the multi-hash instead.
func (w *cantonWallet) SignTopologyTransactions(path accounts.DerivationPath, transactions [][]byte, challenge []byte) ([]byte, []byte, error) {
	if len(transactions) == 0 {
		return nil, nil, errors.New("ledger: at least one topology transaction is required")
	}
	if len(challenge) != 0 && len(challenge) != challengeLength {
		return nil, nil, fmt.Errorf("ledger: invalid challenge length %d, want %d", len(challenge), challengeLength)
	}

	w.mux.Lock()
	defer w.mux.Unlock()
	if w.device == nil {
		return nil, nil, accounts.ErrWalletClosed
	}

	signature, challengeSignature, err := w.ledgerSign(P1_SIGN_UNTYPED_VERSIONED_MESSAGE, path, transactions, challenge)
	if err != nil {
		return nil, nil, err
	}

	return signature, challengeSignature, nil
}

// Low-level device communication

func (w *cantonWallet) ledgerExchange(ins cantonInstruction, p1 cantonParam1, p2 cantonParam2, data []byte) ([]byte, error) {
	if len(data) > maxAPDUDataLength {
		return nil, fmt.Errorf("ledger: APDU data too long: %d > %d", len(data), maxAPDUDataLength)
	}

	// Construct the message payload, possibly split into multiple chunks
	apdu := make([]byte, 2, 7+len(data))

	binary.BigEndian.PutUint16(apdu, uint16(5+len(data)))                               //nolint:gosec // guarded by the length check above
	apdu = append(apdu, []byte{CLA, byte(ins), byte(p1), byte(p2), byte(len(data))}...) //nolint:gosec // guarded by the length check above
	apdu = append(apdu, data...)

	// Stream all the chunks to the device
	header := []byte{0x01, 0x01, 0x05, 0x00, 0x00} // Channel ID and command tag appended
	chunk := make([]byte, 64)
	space := len(chunk) - len(header)

	for i := 0; len(apdu) > 0; i++ {
		// Construct the new message to stream
		chunk = append(chunk[:0], header...)
		binary.BigEndian.PutUint16(chunk[3:], uint16(i))

		if len(apdu) > space {
			chunk = append(chunk, apdu[:space]...) //nolint:makezero // header is preallocated
			apdu = apdu[space:]
		} else {
			chunk = append(chunk, apdu...) //nolint:makezero // header is preallocated
			apdu = nil
		}
		// Send over to the device
		w.log.Trace("Data chunk sent to the Ledger", "chunk", hexutil.Bytes(chunk))
		if _, err := w.device.Write(chunk); err != nil {
			return nil, err
		}
	}

	// Stream the reply back from the wallet in 64 byte chunks
	var reply []byte
	chunk = chunk[:64] // Yeah, we surely have enough space
	for {
		// Read the next chunk from the Ledger wallet
		if _, err := io.ReadFull(w.device, chunk); err != nil {
			return nil, err
		}
		w.log.Trace("Data chunk received from the Ledger", "chunk", hexutil.Bytes(chunk))

		// Make sure the transport header matches
		if chunk[0] != 0x01 || chunk[1] != 0x01 || chunk[2] != 0x05 {
			return nil, errLedgerReplyInvalidHeader
		}
		// If it's the first chunk, retrieve the total message length
		var payload []byte

		if chunk[3] == 0x00 && chunk[4] == 0x00 {
			reply = make([]byte, 0, int(binary.BigEndian.Uint16(chunk[5:7])))
			payload = chunk[7:]
		} else {
			payload = chunk[5:]
		}
		// Append to the reply and stop when filled up
		if left := cap(reply) - len(reply); left > len(payload) {
			reply = append(reply, payload...)
		} else {
			reply = append(reply, payload[:left]...)
			break
		}
	}

	// The reply is terminated by a 2 byte status word which must be checked before the
	// response data is interpreted, otherwise device side errors are silently swallowed.
	if len(reply) < 2 {
		return nil, errors.New("ledger: reply too short to contain a status word")
	}
	if sw := binary.BigEndian.Uint16(reply[len(reply)-2:]); sw != swOK {
		return nil, LedgerStatusError(sw)
	}

	return reply[:len(reply)-2], nil
}

// ledgerSendChunked splits payload across as many APDUs as needed and sets the P2 flags
// according to the Canton app protocol. isFinal marks the payload as the last message of the
// whole signing sequence, which is the only APDU that carries response data.
func (w *cantonWallet) ledgerSendChunked(ins cantonInstruction, p1 cantonParam1, payload []byte, isFinal bool) ([]byte, error) {
	var reply []byte

	for first := true; first || len(payload) > 0; first = false {
		size := min(len(payload), maxAPDUDataLength)
		chunk := payload[:size]
		payload = payload[size:]

		p2 := P2_MORE
		if len(payload) == 0 {
			if isFinal {
				p2 = P2_MESSAGE_END
			} else {
				p2 = P2_MORE | P2_MESSAGE_END
			}
		}

		var err error
		if reply, err = w.ledgerExchange(ins, p1, p2, chunk); err != nil {
			return nil, err
		}
	}

	return reply, nil
}

// encodeDerivationPath serializes a BIP32 derivation path as expected by the Canton app:
// a single length byte followed by the big endian encoded path components.
func encodeDerivationPath(derivationPath accounts.DerivationPath) []byte {
	path := make([]byte, 1+4*len(derivationPath))
	path[0] = byte(len(derivationPath)) //nolint:gosec
	for i, component := range derivationPath {
		binary.BigEndian.PutUint32(path[1+4*i:], component)
	}

	return path
}

// ledgerSign runs a full SIGN_TX sequence: the derivation path (plus the optional attestation
// challenge) first, then every message, with the last one carrying the signature response.
func (w *cantonWallet) ledgerSign(p1 cantonParam1, derivationPath accounts.DerivationPath, messages [][]byte, challenge []byte) ([]byte, []byte, error) {
	pathData := encodeDerivationPath(derivationPath)
	pathData = append(pathData, challenge...)

	// Ensure the hub doesn't refresh while we're in the middle of a signing operation
	// See (hub *Hub) refreshWallets() for details
	w.hub.commsLock.Lock()
	w.hub.commsPend++
	w.hub.commsLock.Unlock()
	defer func() {
		w.hub.commsLock.Lock()
		w.hub.commsPend--
		w.hub.commsLock.Unlock()
	}()

	// The first APDU carries the path and must always announce that more data follows.
	if _, err := w.ledgerExchange(INS_SIGN_TX, p1, P2_FIRST|P2_MORE, pathData); err != nil {
		return nil, nil, err
	}

	var reply []byte
	for i, message := range messages {
		var err error
		if reply, err = w.ledgerSendChunked(INS_SIGN_TX, p1, message, i == len(messages)-1); err != nil {
			return nil, nil, err
		}
	}

	return parseSignatureResponse(reply, len(challenge) > 0)
}

// parseSignatureResponse decodes the signature reply of the Canton app, which is
// [sig_len][signature][v] or, when an attestation challenge was supplied,
// [sig_len][signature][v][sig_len][challenge signature].
func parseSignatureResponse(reply []byte, withChallenge bool) ([]byte, []byte, error) {
	readSignature := func(buf []byte) ([]byte, []byte, error) {
		if len(buf) < 1 || len(buf) < 1+int(buf[0]) {
			return nil, nil, errors.New("ledger: reply lacks signature entry")
		}
		signature := buf[1 : 1+int(buf[0])]
		if len(signature) != ed25519.SignatureSize {
			return nil, nil, fmt.Errorf("ledger: unexpected signature length %d, want %d", len(signature), ed25519.SignatureSize)
		}

		return bytes.Clone(signature), buf[1+int(buf[0]):], nil
	}

	signature, rest, err := readSignature(reply)
	if err != nil {
		return nil, nil, err
	}
	if len(rest) < 1 {
		return nil, nil, errors.New("ledger: reply lacks recovery byte")
	}
	rest = rest[1:] // Skip the recovery byte, it carries no information for Ed25519

	if !withChallenge {
		return signature, nil, nil
	}

	challengeSignature, _, err := readSignature(rest)
	if err != nil {
		return nil, nil, fmt.Errorf("ledger: failed to read challenge signature: %w", err)
	}

	return signature, challengeSignature, nil
}
