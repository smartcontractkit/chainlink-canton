package client

import "testing"

func TestLogicalSynchronizerID(t *testing.T) {
	t.Parallel()

	physical := "global-domain::1220f22a8b8f2d813c25b9a684dc4dd52b532a0174d8e73a13cdf2baabfff7518337::34-0"
	logical := "global-domain::1220f22a8b8f2d813c25b9a684dc4dd52b532a0174d8e73a13cdf2baabfff7518337"

	if got := LogicalSynchronizerID(physical); got != logical {
		t.Fatalf("physical: got %q, want %q", got, logical)
	}
	if got := LogicalSynchronizerID(logical); got != logical {
		t.Fatalf("logical passthrough: got %q, want %q", got, logical)
	}
	if got := LogicalSynchronizerID("global"); got != "global" {
		t.Fatalf("alias passthrough: got %q", got)
	}
}
