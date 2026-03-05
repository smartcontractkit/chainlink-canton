package common

// https://github.com/tronprotocol/java-tron/blob/c136a26f8140b17f2c05df06fb5efb1bb47d3baa/protocol/src/main/protos/api/api.proto#L1041
const (
	ResponseCodeSuccess                      = "SUCCESS"
	ResponseCodeSigError                     = "SIGERROR" // error in signature
	ResponseCodeContractValidateError        = "CONTRACT_VALIDATE_ERROR"
	ResponseCodeContractExeError             = "CONTRACT_EXE_ERROR"
	ResponseCodeBandwidthError               = "BANDWITH_ERROR"
	ResponseCodeDupTransactionError          = "DUP_TRANSACTION_ERROR"
	ResponseCodeTaposError                   = "TAPOS_ERROR"
	ResponseCodeTooBigTransactionError       = "TOO_BIG_TRANSACTION_ERROR"
	ResponseCodeTransactionExpirationError   = "TRANSACTION_EXPIRATION_ERROR"
	ResponseCodeServerBusy                   = "SERVER_BUSY"
	ResponseCodeNoConnection                 = "NO_CONNECTION"
	ResponseCodeNotEnoughEffectiveConnection = "NOT_ENOUGH_EFFECTIVE_CONNECTION"
	ResponseCodeBlockUnsolidified            = "BLOCK_UNSOLIDIFIED"
	ResponseCodeOtherError                   = "OTHER_ERROR"
)
