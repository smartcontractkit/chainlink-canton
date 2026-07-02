package replication

// ── RecordLedgerOffset ───────────────────────────────────────────────────────

// RecordLedgerOffsetInput is the input to [RecordLedgerOffsetOp].
type RecordLedgerOffsetInput struct {
	ParticipantID  string `json:"participant_id"`
	SynchronizerID string `json:"synchronizer_id"`
}

// RecordLedgerOffsetOutput is the output of [RecordLedgerOffsetOp].
type RecordLedgerOffsetOutput struct {
	ParticipantID    string `json:"participant_id"`
	LedgerOffset     int64  `json:"ledger_offset"`
	TimestampSeconds int64  `json:"timestamp_seconds"`
}

// ── RecordTargetLedgerOffset ─────────────────────────────────────────────────

// RecordTargetLedgerOffsetInput is the input to [RecordTargetLedgerOffsetOp].
type RecordTargetLedgerOffsetInput struct {
	ParticipantID  string `json:"participant_id"`
	SynchronizerID string `json:"synchronizer_id"`
}

// RecordTargetLedgerOffsetOutput is the output of [RecordTargetLedgerOffsetOp].
type RecordTargetLedgerOffsetOutput struct {
	ParticipantID    string `json:"participant_id"`
	LedgerOffset     int64  `json:"ledger_offset"`
	TimestampSeconds int64  `json:"timestamp_seconds"`
}

// ── DisconnectSynchronizer ──────────────────────────────────────────────────

// DisconnectSynchronizerInput is the input to [DisconnectSynchronizerOp].
type DisconnectSynchronizerInput struct {
	ParticipantID     string `json:"participant_id"`
	SynchronizerAlias string `json:"synchronizer_alias"`
}

// DisconnectSynchronizerOutput is the output of [DisconnectSynchronizerOp].
type DisconnectSynchronizerOutput struct {
	ParticipantID string `json:"participant_id"`
	Disconnected  bool   `json:"disconnected"`
}

// ReconnectSynchronizerInput is the input to [ReconnectSynchronizerOp].
type ReconnectSynchronizerInput struct {
	ParticipantID     string `json:"participant_id"`
	SynchronizerAlias string `json:"synchronizer_alias"`
}

// ReconnectSynchronizerOutput is the output of [ReconnectSynchronizerOp].
type ReconnectSynchronizerOutput struct {
	ParticipantID string `json:"participant_id"`
	Reconnected   bool   `json:"reconnected"`
}

// ── ExportAcs ────────────────────────────────────────────────────────────────

// ExportAcsInput is the input to [ExportAcsOp].
type ExportAcsInput struct {
	ParticipantID    string   `json:"participant_id"`
	PartyIDs         []string `json:"party_ids"`
	SynchronizerID   string   `json:"synchronizer_id"`
	LedgerOffset     int64    `json:"ledger_offset"`
	TimestampSeconds int64    `json:"timestamp_seconds"`
}

// ExportAcsOutput is the output of [ExportAcsOp].
type ExportAcsOutput struct {
	ParticipantID  string `json:"participant_id"`
	AcsSnapshotB64 string `json:"acs_snapshot_b64"`
	SHA256         string `json:"sha256"`
	SizeBytes      int    `json:"size_bytes"`
}

// ── ImportAcs ────────────────────────────────────────────────────────────────

// ImportAcsInput is the input to [ImportAcsOp].
type ImportAcsInput struct {
	ParticipantID     string `json:"participant_id"`
	SynchronizerID    string `json:"synchronizer_id"`
	SynchronizerAlias string `json:"synchronizer_alias"`
	AcsSnapshotB64    string `json:"acs_snapshot_b64"`
	ExpectedSHA256    string `json:"expected_sha256,omitempty"`
}

// ImportAcsOutput is the output of [ImportAcsOp].
type ImportAcsOutput struct {
	ParticipantID string `json:"participant_id"`
	Imported      bool   `json:"imported"`
}

// ── ClearOnboardingFlag ──────────────────────────────────────────────────────

// ClearOnboardingFlagInput is the input to [ClearOnboardingFlagOp].
type ClearOnboardingFlagInput struct {
	ParticipantID        string `json:"participant_id"`
	PartyID              string `json:"party_id"`
	SynchronizerID       string `json:"synchronizer_id"`
	BeginOffsetExclusive int64  `json:"begin_offset_exclusive"`
}

// ClearOnboardingFlagOutput is the output of [ClearOnboardingFlagOp].
type ClearOnboardingFlagOutput struct {
	ParticipantID string `json:"participant_id"`
	Onboarded     bool   `json:"onboarded"`
}
