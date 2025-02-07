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

	expected := searchResponse{
		ModelType: "bpemb_attention",
		ParsedAddresses: map[string]ParsedAddress{
			"2325 Rue de l'Université, Québec, QC G1V 0A6": {
				StreetNumber:    "2325",
				StreetName:      "rue de l'université",
				Unit:            "",
				Municipality:    "québec",
				Province:        "qc",
				PostalCode:      "g1v 0a6",
				Orientation:     "",
				GeneralDelivery: "",
			},
			"350 rue des Lilas Ouest Quebec city Quebec G1L 1B6": {
				StreetNumber:    "350",
				StreetName:      "rue des lilas ouest",
				Unit:            "",
				Municipality:    "quebec city",
				Province:        "quebec",
				PostalCode:      "g1l 1b6",
				Orientation:     "",
				GeneralDelivery: "",
			},
		},
		Version: "cfb190902476376573591c0ec6f91ece",
	}

	require.Equal(t, expected, resp)
}
