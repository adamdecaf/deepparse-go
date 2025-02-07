package deepparsego

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client interface {
	ParseAddresses(ctx context.Context, model Model, addresses []string) (SearchResponse, error)
}

func NewClient(httpClient *http.Client, baseAddress string) Client {
	return &client{
		httpClient:  httpClient,
		baseAddress: baseAddress,
	}
}

type client struct {
	httpClient  *http.Client
	baseAddress string
}

type Model string

const (
	ModelFastText          Model = "fasttext"
	ModelFastTextAttention Model = "fasttext-attention"
	ModelFastTextLight     Model = "fasttext-light"
	ModelBPEmb             Model = "bpemb"
	ModelBPEmbAttention    Model = "bpemb-attention"
)

func (c *client) ParseAddresses(ctx context.Context, model Model, addresses []string) (SearchResponse, error) {
	var out SearchResponse

	var body searchRequest
	for _, addr := range addresses {
		body = append(body, rawAddress{Raw: addr})
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return out, fmt.Errorf("encoding request addresses: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseAddress+"/parse/"+string(model), &buf)
	if err != nil {
		return out, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return out, fmt.Errorf("parsing addresses: %w", err)
	}
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	var wrapper searchResponse
	err = json.NewDecoder(resp.Body).Decode(&wrapper)
	if err != nil {
		return out, fmt.Errorf("reading parsed addresses response: %w", err)
	}

	out.Model = Model(wrapper.ModelType)
	out.Version = wrapper.Version

	for _, addr := range wrapper.ParsedAddresses {
		out.Addresses = append(out.Addresses, addr)
	}

	return out, nil
}

// [
//
//	{\"raw\": \"350 rue des Lilas Ouest Quebec city Quebec G1L 1B6\"},
//	{\"raw\": \"2325 Rue de l'Université, Québec, QC G1V 0A6\"}
//
// ]
type rawAddress struct {
	Raw string `json:"raw"`
}
type searchRequest []rawAddress

// SearchResponse is the model returned from parsing addresses
type SearchResponse struct {
	Model     Model
	Addresses []ParsedAddress
	Version   string
}

// ParsedAddress is the fields of a parsed address
type ParsedAddress struct {
	StreetNumber    string `json:"StreetNumber"`
	StreetName      string `json:"StreetName"`
	Unit            string `json:"Unit"`
	Municipality    string `json:"Municipality"`
	Province        string `json:"Province"`
	PostalCode      string `json:"PostalCode"`
	Orientation     string `json:"Orientation"`
	GeneralDelivery string `json:"GeneralDelivery"`
}

type searchResponse struct {
	ModelType       string                   `json:"model_type"`
	ParsedAddresses map[string]ParsedAddress `json:"parsed_addresses"`
	Version         string                   `json:"version"`
}
