package testhelpers

import (
	"context"
	"fmt"

	participantv30 "github.com/digital-asset/dazl-client/v8/go/api/com/digitalasset/canton/admin/participant/v30"
)

func UploadDARstoMultipleParticipants(ctx context.Context, dars [][]byte, participants ...Participant) ([]string, error) {
	darData := make([]*participantv30.UploadDarRequest_UploadDarData, 0, len(dars))
	for _, dar := range dars {
		darData = append(darData, &participantv30.UploadDarRequest_UploadDarData{
			Bytes: dar,
		})
	}

	var packageIDs []string
	for _, participant := range participants {
		res, err := participant.PackageServiceClient.UploadDar(ctx, &participantv30.UploadDarRequest{
			Dars:               darData,
			VetAllPackages:     true,
			SynchronizeVetting: true,
		})
		if err != nil {
			return nil, fmt.Errorf("uploadDAR to participant %q failed: %w", participant.Name, err)
		}
		packageIDs = append(packageIDs, res.GetDarIds()...)
	}

	return packageIDs, nil
}
