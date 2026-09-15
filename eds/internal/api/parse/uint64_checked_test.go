package parse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUint64Checked(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      string
		want    uint64
		wantErr bool
	}{
		{
			name: "zero",
			in:   "0",
			want: 0,
		},
		{
			name: "simple integer",
			in:   "12345",
			want: 12345,
		},
		{
			name: "trailing dot",
			in:   "12345.",
			want: 12345,
		},
		{
			name: "decimal point zero",
			in:   "12345.0",
			want: 12345,
		},
		{
			name: "decimal point zeros",
			in:   "12345.000",
			want: 12345,
		},
		{
			name: "max uint64",
			in:   "18446744073709551615",
			want: ^uint64(0),
		},
		{
			name:    "overflow uint64",
			in:      "18446744073709551616",
			wantErr: true,
		},
		{
			name:    "negative",
			in:      "-1",
			wantErr: true,
		},
		{
			name:    "fractional",
			in:      "1.1",
			wantErr: true,
		},
		{
			name:    "empty",
			in:      "",
			wantErr: true,
		},
		{
			name:    "whitespace not allowed",
			in:      " 1",
			wantErr: true,
		},
		{
			name:    "not a number",
			in:      "nope",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Uint64Checked(tt.in)
			if tt.wantErr {
				require.Errorf(t, err, "expected error for input %q, got nil (value=%d)", tt.in, got)
				return
			}

			require.NoErrorf(t, err, "unexpected error for input %q: %v", tt.in, err)
			require.Equalf(t, tt.want, got, "Uint64Checked(%q)=%d, want %d", tt.in, got, tt.want)
		})
	}
}
