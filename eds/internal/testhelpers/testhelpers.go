package testhelpers

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	oapiCommon "github.com/smartcontractkit/chainlink-canton/openapi/gen/eds/common"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func MakeHashedAddress(hex string) oapiCommon.RawOrHashedAddress {
	var addr oapiCommon.RawOrHashedAddress
	_ = addr.FromInstanceAddress(hex)
	return addr
}

func MakeRawAddress(raw string) oapiCommon.RawOrHashedAddress {
	var addr oapiCommon.RawOrHashedAddress
	_ = addr.FromRawInstanceAddress(raw)
	return addr
}

func MakeOversizedRequest(size int) io.Reader {
	return strings.NewReader(`{"` + strings.Repeat("a", size-2))
}
