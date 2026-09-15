package chainlink_canton

import (
	_ "embed"
	"os"
	"strings"
)

//go:embed .SPLICE-VERSION
var spliceVersionFile string

const SpliceVersionEnv = "SPLICE_VERSION"

// spliceVersion is parsed from .SPLICE-VERSION, a dotenv-style file that can also be
// passed directly to `docker compose --env-file` to keep the version in a single place.
var spliceVersion = parseSpliceVersion()

func parseSpliceVersion() string {
	for line := range strings.SplitSeq(spliceVersionFile, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if value, found := strings.CutPrefix(line, SpliceVersionEnv+"="); found {
			return strings.Trim(strings.TrimSpace(value), `"'`)
		}
	}

	return ""
}

func GetSpliceVersion() string {
	if os.Getenv(SpliceVersionEnv) != "" {
		return os.Getenv(SpliceVersionEnv)
	}

	return spliceVersion
}
