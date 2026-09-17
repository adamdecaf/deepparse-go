package deepparsego

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestModelsWrite(t *testing.T) {
	req := searchRequest{
		{Raw: "350 rue des Lilas Ouest Quebec city Quebec G1L 1B6"},
		{Raw: "2325 Rue de l'Université, Québec, QC G1V 0A6"},
	}
	bs, err := json.Marshal(req)
	require.NoError(t, err)

	expected := strings.TrimSpace(`[{"raw":"350 rue des Lilas Ouest Quebec city Quebec G1L 1B6"},{"raw":"2325 Rue de l'Université, Québec, QC G1V 0A6"}]`)

	require.Equal(t, expected, string(bs))
}

func TestModelsRead(t *testing.T) {
	bs, err := os.ReadFile(filepath.Join("testdata", "sample-response.json"))
	require.NoError(t, err)

	var resp searchResponse
	err = json.Unmarshal(bs, &resp)
	require.NoError(t, err)

	require.Equal(t, "bpemb_attention", resp.ModelType)
	require.Equal(t, "cfb190902476376573591c0ec6f91ece", resp.Version)
	require.Len(t, resp.ParsedAddresses, 2)

	firstRaw, first := onlyEntry(t, resp.ParsedAddresses[0])
	require.Equal(t, "350 rue des Lilas Ouest Quebec city Quebec G1L 1B6", firstRaw)
	require.Equal(t, ParsedAddress{
		StreetNumber: "350",
		StreetName:   "rue des lilas ouest",
		Municipality: "quebec city",
		Province:     "quebec",
		PostalCode:   "g1l 1b6",
	}, first)

	secondRaw, second := onlyEntry(t, resp.ParsedAddresses[1])
	require.Equal(t, "2325 Rue de l'Université, Québec, QC G1V 0A6", secondRaw)
	require.Equal(t, ParsedAddress{
		StreetNumber: "2325",
		StreetName:   "rue de l'université",
		Municipality: "québec",
		Province:     "qc",
		PostalCode:   "g1v 0a6",
	}, second)
}

func onlyEntry(t *testing.T, item map[string]ParsedAddress) (string, ParsedAddress) {
	t.Helper()
	require.Len(t, item, 1)
	for raw, addr := range item {
		return raw, addr
	}
	return "", ParsedAddress{}
}
