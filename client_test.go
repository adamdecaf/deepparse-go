package deepparsego

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseAddresses(t *testing.T) {
	const (
		addr1 = "350 rue des Lilas Ouest Quebec city Quebec G1L 1B6"
		addr2 = "2325 Rue de l'Université, Québec, QC G1V 0A6"
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/parse/bpemb-attention", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var req searchRequest
		require.NoError(t, json.Unmarshal(body, &req))
		require.Equal(t, searchRequest{{Raw: addr1}, {Raw: addr2}}, req)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(`{
			"model_type": "bpemb_attention",
			"parsed_addresses": [
				{"` + addr1 + `": {"StreetNumber":"350","StreetName":"rue des lilas ouest","Municipality":"quebec city","Province":"quebec","PostalCode":"g1l 1b6"}},
				{"` + addr2 + `": {"StreetNumber":"2325","StreetName":"rue de l'université","Municipality":"québec","Province":"qc","PostalCode":"g1v 0a6"}}
			],
			"version": "test-version"
		}`))
		require.NoError(t, err)
	}))
	t.Cleanup(srv.Close)

	cc := NewClient(srv.Client(), srv.URL+"/")
	resp, err := cc.ParseAddresses(context.Background(), ModelBPEmbAttention, []string{addr1, addr2})
	require.NoError(t, err)

	require.Equal(t, Model("bpemb_attention"), resp.Model)
	require.Equal(t, "test-version", resp.Version)
	require.Equal(t, []ParsedAddress{
		{
			Raw:          addr1,
			StreetNumber: "350",
			StreetName:   "rue des lilas ouest",
			Municipality: "quebec city",
			Province:     "quebec",
			PostalCode:   "g1l 1b6",
		},
		{
			Raw:          addr2,
			StreetNumber: "2325",
			StreetName:   "rue de l'université",
			Municipality: "québec",
			Province:     "qc",
			PostalCode:   "g1v 0a6",
		},
	}, resp.Addresses)
}

func TestParseAddresses_DuplicatesKeepOrder(t *testing.T) {
	raw := "350 rue des Lilas Ouest Quebec city Quebec G1L 1B6"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/parse/bpemb", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{
			"model_type": "bpemb",
			"parsed_addresses": [
				{"` + raw + `": {"StreetNumber":"350","StreetName":"rue des lilas ouest"}},
				{"` + raw + `": {"StreetNumber":"350","StreetName":"rue des lilas ouest"}}
			],
			"version": "v"
		}`))
		require.NoError(t, err)
	}))
	t.Cleanup(srv.Close)

	resp, err := NewClient(srv.Client(), srv.URL).ParseAddresses(context.Background(), ModelBPEmb, []string{raw, raw})
	require.NoError(t, err)
	require.Len(t, resp.Addresses, 2)
	require.Equal(t, raw, resp.Addresses[0].Raw)
	require.Equal(t, raw, resp.Addresses[1].Raw)
}

func TestParseAddresses_HTTPErrors(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		body       string
		wantDetail string
	}{
		{
			name:       "empty list",
			status:     http.StatusUnprocessableEntity,
			body:       `{"detail":"Addresses parameter must not be empty"}`,
			wantDetail: "Addresses parameter must not be empty",
		},
		{
			name:       "unknown model",
			status:     http.StatusUnprocessableEntity,
			body:       `{"detail":"Parsing model not implemented, available choices: ['bpemb']"}`,
			wantDetail: "Parsing model not implemented",
		},
		{
			name:       "too many addresses",
			status:     http.StatusRequestEntityTooLarge,
			body:       `{"detail":"Too many addresses in a single request (max 1024)."}`,
			wantDetail: "Too many addresses in a single request (max 1024).",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			}))
			t.Cleanup(srv.Close)

			_, err := NewClient(srv.Client(), srv.URL).ParseAddresses(context.Background(), ModelBPEmb, []string{"an address"})
			require.Error(t, err)

			var apiErr *APIError
			require.ErrorAs(t, err, &apiErr)
			require.Equal(t, tc.status, apiErr.StatusCode)
			require.Contains(t, apiErr.Detail, tc.wantDetail)
			require.Contains(t, apiErr.Error(), "deepparse API error")
		})
	}
}

func TestParseAddresses_NilHTTPClient(t *testing.T) {
	cc := NewClient(nil, "http://127.0.0.1:1")
	require.NotNil(t, cc)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err := cc.ParseAddresses(ctx, ModelBPEmb, []string{"x"})
	require.Error(t, err)
}

func TestNewAPIError_NonJSON(t *testing.T) {
	err := newAPIError(http.StatusInternalServerError, []byte("boom"))
	require.Equal(t, 500, err.StatusCode)
	require.Equal(t, "boom", err.Detail)
}
