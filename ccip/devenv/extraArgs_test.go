package devenv

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
)

func TestEncodeGenericExtraArgsV3(t *testing.T) {
	t.Parallel()

	args := &GenericExtraArgsV3{
		GasLimit:           500_000,
		BlockConfirmations: 12,
		CCVs:               [][]byte{hexutil.MustDecode("0xbf9f84ade4dfe17de8e6ccfcb543f7fcddfced84"), hexutil.MustDecode("0xdcaf44b3cca7a8a50feaae70d3b7aaec2f040d09")},
		CCVArgs:            [][]byte{hexutil.MustDecode("0x6578656375746f7220617267756d656e742031"), hexutil.MustDecode("0x746869732069732061206d756368206c6f6e67657220617267756d656e7420666f72206578656375746f72206e756d6265722074776f")},
		Executor:           hexutil.MustDecode("0xb3db8bbdbfee9b6cd2368f0dcbbf8314ec25fc0d"),
		ExecutorArgs:       []byte{},
		TokenReceiver:      hexutil.MustDecode("0x9b618e642add81cf490b4bdfc2cfc163dcf8e0a7"),
		TokenArgs:          hexutil.MustDecode("0x68656c6c6f20776f726c64"),
	}
	encoded, err := EncodeGenericExtraArgsV3(args)
	require.NoError(t, err)

	// Values taken from Daml ExtraArgsCodec.testExtraArgsV3Codec test
	expected := hexutil.MustDecode("0xa69dd4aa0007a120000c0214bf9f84ade4dfe17de8e6ccfcb543f7fcddfced8400136578656375746f7220617267756d656e74203114dcaf44b3cca7a8a50feaae70d3b7aaec2f040d090036746869732069732061206d756368206c6f6e67657220617267756d656e7420666f72206578656375746f72206e756d6265722074776f14b3db8bbdbfee9b6cd2368f0dcbbf8314ec25fc0d0000149b618e642add81cf490b4bdfc2cfc163dcf8e0a7000b68656c6c6f20776f726c64")
	require.Equal(t, expected, encoded)
}
