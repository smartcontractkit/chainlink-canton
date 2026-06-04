package tokenpool

import (
	"testing"
)

func TestValidateTokenDecimals(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		decimals int64
		wantErr  bool
	}{
		{name: "min", decimals: 0},
		{name: "evm typical", decimals: 18},
		{name: "max", decimals: 37},
		{name: "negative", decimals: -1, wantErr: true},
		{name: "above max", decimals: 100, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateTokenDecimals(tt.decimals)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateTokenDecimals(%d) error = %v, wantErr %v", tt.decimals, err, tt.wantErr)
			}
		})
	}
}
