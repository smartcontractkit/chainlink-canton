package fullnode

import (
	"net/http"

	"github.com/fbsobreira/gotron-sdk/pkg/http/soliditynode"
)

type Client struct {
	*soliditynode.Client
}

func NewClient(baseURL string, client *http.Client) *Client {
	return &Client{
		Client: soliditynode.NewClient(baseURL, client),
	}
}
