package cmd

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func Test_topologyTransactionMultihash(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                 string
		topologyTransactions []string
		want                 string
	}{
		{
			name: "",
			topologyTransactions: []string{
				"0aec01080110011ae5014ae2010a4a657874313a3a313232303133313639383430313235643339646539323531316233343761663365323466373137346534366136623930333865306536643638323330363566353935613610011a550a517061727469636970616e743a3a31323230653663363161393631316663396263653262613561346431303539626131616566366331313164343262663764303238343030313961333564616235343438611002323b0a3710041a2c302a300506032b65700321005e03e77cca028d9a7ff68b505598aab9d31db17e12decb28a2c287eca69f7f9d2a0301050430011001101e",
			},
			want: "1220e498e0eed91ca03d24546ed98f3ee308a53643dec8da503ca8da52fac8c3941c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			transactions := make([][]byte, len(tt.topologyTransactions))
			for i, transaction := range tt.topologyTransactions {
				transactions[i], _ = hex.DecodeString(transaction)
			}
			want, _ := hex.DecodeString(tt.want)

			got := topologyTransactionMultihash(transactions)
			if !bytes.Equal(got, want) {
				t.Errorf("topologyTransactionMultihash() got = %v, want %v", got, want)
			}
		})
	}
}
