package auth

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type TokenProvider func() (string, error)

type BearerTokenAuth struct {
	token    string
	provider TokenProvider
}

func NewBearerToken(token string) *BearerTokenAuth {
	return &BearerTokenAuth{
		token: token,
	}
}

func NewBearerTokenProvider(provider TokenProvider) *BearerTokenAuth {
	return &BearerTokenAuth{
		provider: provider,
	}
}

func (b *BearerTokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	token, err := b.getToken()
	if err != nil {
		return nil, err
	}

	if token == "" {
		return nil, nil
	}

	return map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", token),
	}, nil
}

func (b *BearerTokenAuth) RequireTransportSecurity() bool {
	return false
}

func (b *BearerTokenAuth) UnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		token, err := b.getToken()
		if err != nil {
			return err
		}

		if token != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", token))
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (b *BearerTokenAuth) StreamInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		token, err := b.getToken()
		if err != nil {
			return nil, err
		}

		if token != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", token))
		}

		return streamer(ctx, desc, cc, method, opts...)
	}
}

func (b *BearerTokenAuth) getToken() (string, error) {
	if b.provider != nil {
		return b.provider()
	}
	return b.token, nil
}
