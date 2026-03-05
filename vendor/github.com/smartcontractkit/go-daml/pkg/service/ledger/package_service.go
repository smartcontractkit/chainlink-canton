package ledger

import (
	"context"

	"google.golang.org/grpc"

	v2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
	"github.com/smartcontractkit/go-daml/pkg/model"
)

type PackageService interface {
	ListPackages(ctx context.Context, req *model.ListPackagesRequest) (*model.ListPackagesResponse, error)
	GetPackage(ctx context.Context, req *model.GetPackageRequest) (*model.GetPackageResponse, error)
	GetPackageStatus(ctx context.Context, req *model.GetPackageStatusRequest) (*model.GetPackageStatusResponse, error)
}

type packageService struct {
	client v2.PackageServiceClient
}

func NewPackageServiceClient(conn *grpc.ClientConn) *packageService {
	client := v2.NewPackageServiceClient(conn)
	return &packageService{
		client: client,
	}
}

func (c *packageService) ListPackages(ctx context.Context, req *model.ListPackagesRequest) (*model.ListPackagesResponse, error) {
	protoReq := &v2.ListPackagesRequest{}

	resp, err := c.client.ListPackages(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return &model.ListPackagesResponse{
		PackageIDs: resp.PackageIds,
	}, nil
}

func (c *packageService) GetPackage(ctx context.Context, req *model.GetPackageRequest) (*model.GetPackageResponse, error) {
	protoReq := &v2.GetPackageRequest{
		PackageId: req.PackageID,
	}

	resp, err := c.client.GetPackage(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return &model.GetPackageResponse{
		ArchivePayload: resp.ArchivePayload,
		HashFunction:   hashFunctionFromProto(resp.HashFunction),
		Hash:           resp.Hash,
	}, nil
}

func (c *packageService) GetPackageStatus(ctx context.Context, req *model.GetPackageStatusRequest) (*model.GetPackageStatusResponse, error) {
	protoReq := &v2.GetPackageStatusRequest{
		PackageId: req.PackageID,
	}

	resp, err := c.client.GetPackageStatus(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return &model.GetPackageStatusResponse{
		PackageStatus: packageStatusFromProto(resp.PackageStatus),
	}, nil
}

func hashFunctionFromProto(hf v2.HashFunction) model.HashFunction {
	switch hf {
	case v2.HashFunction_HASH_FUNCTION_SHA256:
		return model.HashFunctionSHA256
	default:
		return model.HashFunctionSHA256
	}
}

func packageStatusFromProto(ps v2.PackageStatus) model.PackageStatus {
	switch ps {
	case v2.PackageStatus_PACKAGE_STATUS_REGISTERED:
		return model.PackageStatusRegistered
	default:
		return model.PackageStatusUnknown
	}
}
