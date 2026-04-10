package testhelpers

import (
	"context"
	"fmt"

	participantv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/admin/participant/v30"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/canton"
	"google.golang.org/grpc/status"
)

func UploadDARstoMultipleParticipants(ctx context.Context, dars [][]byte, participants ...canton.Participant) ([]string, error) {
	var packageIDs []string
	for _, participant := range participants {
		for i, dar := range dars {
			res, err := participant.AdminServices.Package.UploadDar(ctx, &participantv30.UploadDarRequest{
				Dars: []*participantv30.UploadDarRequest_UploadDarData{
					{
						Bytes:                 dar,
						Description:           nil,
						ExpectedMainPackageId: nil,
					},
				},
				VetAllPackages:     true,
				SynchronizeVetting: true,
			})
			if err != nil {
				// Upload dars one-by-one and print error details to be able to debug
				s, ok := status.FromError(err)
				if ok {
					fmt.Println("gRPC error details:", s.Details())
					fmt.Println("gRPC error message:", s.Message())
					fmt.Println("gRPC error code:", s.Code())
					fmt.Println("gRPC error err:", s.Err())
				}

				return nil, fmt.Errorf("uploading dar #%d to participant %q failed: %w", i, participant.Name, err)
			}
			packageIDs = append(packageIDs, res.GetDarIds()...)
		}
	}

	return packageIDs, nil
}
