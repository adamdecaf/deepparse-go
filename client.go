package deepparsego

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client is the HTTP API for GRAAL-Research/deepparse's FastAPI parser.
type Client interface {
	ParseAddresses(ctx context.Context, model Model, addresses []string) (SearchResponse, error)
}

// NewClient returns a Client that talks to a running deepparse HTTP API.
// If httpClient is nil, a client with a 30s timeout is used.
// baseAddress is the origin only, for example "http://localhost:8000".
func NewClient(httpClient *http.Client, baseAddress string) Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &client{
		httpClient:  httpClient,
		baseAddress: strings.TrimRight(baseAddress, "/"),
	}
}

type client struct {
	httpClient  *http.Client
	baseAddress string
}

// Model is a deepparse parsing model path segment.
type Model string

const (
	ModelFastText          Model = "fasttext"
	ModelFastTextAttention Model = "fasttext-attention"
	ModelFastTextLight     Model = "fasttext-light"
	ModelBPEmb             Model = "bpemb"
	ModelBPEmbAttention    Model = "bpemb-attention"
)

// MaxAddressesPerRequest is the default deepparse REST limit (HTTP 413 above this).
const MaxAddressesPerRequest = 1024

func (c *client) ParseAddresses(ctx context.Context, model Model, addresses []string) (SearchResponse, error) {
	var out SearchResponse

	body := make(searchRequest, 0, len(addresses))
	for _, addr := range addresses {
		body = append(body, rawAddress{Raw: addr})
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return out, fmt.Errorf("encoding request addresses: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseAddress+"/parse/"+string(model), &buf)
	if err != nil {
		return out, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return out, fmt.Errorf("parsing addresses: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return out, fmt.Errorf("reading parsed addresses response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return out, newAPIError(resp.StatusCode, raw)
	}

	var wrapper searchResponse
	if err := json.Unmarshal(raw, &wrapper); err != nil {
		return out, fmt.Errorf("reading parsed addresses response: %w", err)
	}

	out.Model = Model(wrapper.ModelType)
	out.Version = wrapper.Version
	out.Addresses = make([]ParsedAddress, 0, len(wrapper.ParsedAddresses))

	for _, item := range wrapper.ParsedAddresses {
		for rawAddr, addr := range item {
			addr.Raw = rawAddr
			out.Addresses = append(out.Addresses, addr)
		}
	}

	return out, nil
}

type rawAddress struct {
	Raw string `json:"raw"`
}

type searchRequest []rawAddress

// SearchResponse is the model returned from parsing addresses.
type SearchResponse struct {
	Model     Model
	Addresses []ParsedAddress
	Version   string
}

// ParsedAddress is the fields of a parsed address.
type ParsedAddress struct {
	Raw             string `json:"-"`
	StreetNumber    string `json:"StreetNumber"`
	StreetName      string `json:"StreetName"`
	Unit            string `json:"Unit"`
	Municipality    string `json:"Municipality"`
	Province        string `json:"Province"`
	PostalCode      string `json:"PostalCode"`
	Orientation     string `json:"Orientation"`
	GeneralDelivery string `json:"GeneralDelivery"`
}

// searchResponse matches deepparse 0.11.0+: parsed_addresses is a list of
// {raw: parsed} objects so duplicates keep their order.
type searchResponse struct {
	ModelType       string                     `json:"model_type"`
	ParsedAddresses []map[string]ParsedAddress `json:"parsed_addresses"`
	Version         string                     `json:"version"`
}
