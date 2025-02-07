package deepparsego

import (
	"context"
)

func NewMockClient() *MockClient {
	return &MockClient{}
}

type MockClient struct {
	Err error

	ParsedAddresses []ParsedAddress
}

var _ Client = (&MockClient{})

func (c *MockClient) ParseAddresses(ctx context.Context, model Model, addresses []string) (SearchResponse, error) {
	if c.Err != nil {
		var out SearchResponse
		return out, c.Err
	}

	return SearchResponse{
		Model:     model,
		Addresses: c.ParsedAddresses,
		Version:   "mock",
	}, nil
}
