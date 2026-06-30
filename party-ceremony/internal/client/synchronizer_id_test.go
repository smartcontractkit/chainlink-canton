package client

import (
	"context"
	"testing"
)

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

type stubSyncLister struct {
	syncs []SynchronizerInfo
	err   error
}

func (s stubSyncLister) ListConnectedSynchronizers(_ context.Context) ([]SynchronizerInfo, error) {
	if s.err != nil {
		return nil, s.err
	}

	return s.syncs, nil
}

func TestResolvePhysicalSynchronizerID(t *testing.T) {
	t.Parallel()

	physical := "global-domain::1220f22a8b8f2d813c25b9a684dc4dd52b532a0174d8e73a13cdf2baabfff7518337::34-0"
	logical := LogicalSynchronizerID(physical)

	lister := stubSyncLister{syncs: []SynchronizerInfo{{
		Alias:          "global",
		SynchronizerID: physical,
	}}}

	got, err := ResolvePhysicalSynchronizerID(t.Context(), lister, "global")
	if err != nil {
		t.Fatalf("resolve global alias: %v", err)
	}
	if got != physical {
		t.Fatalf("got %q, want %q", got, physical)
	}

	got, err = ResolvePhysicalSynchronizerID(t.Context(), lister, logical)
	if err != nil {
		t.Fatalf("resolve logical id: %v", err)
	}
	if got != physical {
		t.Fatalf("logical: got %q, want %q", got, physical)
	}

	got, err = ResolvePhysicalSynchronizerID(t.Context(), lister, physical)
	if err != nil {
		t.Fatalf("resolve physical passthrough: %v", err)
	}
	if got != physical {
		t.Fatalf("physical passthrough: got %q", got)
	}
}
