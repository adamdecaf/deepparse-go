package deepparsego

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	cc := NewClient(httpClient, "http://localhost:8000")
	require.NotNil(t, cc)

	ctx := context.Background()
	resp, err := cc.ParseAddresses(ctx, ModelBPEmbAttention, []string{
		"350 rue des Lilas Ouest Quebec city Quebec G1L 1B6",
		"2325 Rue de l'Université, Québec, QC G1V 0A6",
	})
	require.NoError(t, err)

	require.Equal(t, Model("bpemb_attention"), resp.Model) // TODO(adam): API returns _ but requires -
	require.Len(t, resp.Addresses, 2)
	require.NotEmpty(t, resp.Version)

	expected := []ParsedAddress{
		{
			StreetNumber: "350",
			StreetName:   "rue des lilas ouest",
			Municipality: "quebec city",
			Province:     "quebec",
			PostalCode:   "g1l 1b6",
		},
		{
			StreetNumber: "2325",
			StreetName:   "rue de l'université",
			Municipality: "québec",
			Province:     "qc",
			PostalCode:   "g1v 0a6",
		},
	}
	require.ElementsMatch(t, expected, resp.Addresses)
}

func BenchmarkClient(b *testing.B) {
	ctx := context.Background()

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	cc := NewClient(httpClient, "http://localhost:8000")
	require.NotNil(b, cc)

	// same as moov-io/watchman's ./pkg/address/address_libpostal_test.go#Benchmark_ParseAddress
	inputs := []string{
		"Flat 7B, Tower 2, Ocean Financial Centre, 12 Marina Boulevard, Singapore 018982",
		"Room 1403, West Wing, Trading Complex No. 5, 47 Al Souq Street, Dubai, United Arab Emirates",
		"Office 892, Floor 8, Edificio Comercial Torres, Avenida Balboa y Calle 42, Panama City, Panama",
		"Unit 15, 3rd Floor, 123 Pyongyang Industrial Zone, Rangnang District, Pyongyang, DPRK",
		"Suite 405, Business Center Red Square, 17 Tverskaya Street, Moscow 125009, Russian Federation",
		"Warehouse 23, Port Zone B, Terminal 4, Latakia Port Complex, Latakia, Syria",
		"Office 78, Tehran Trade Tower, Block 2, Valiasr Street, Tehran 19395-4791, Iran",
		"Villa 15, Street 7, Block 4, Diplomatic Quarter, Caracas 1010, Venezuela",
		"Room 2201, Finance Plaza Building, 333 Lujiazui Ring Road, Shanghai 200120, China",
		"Suite 17, Victoria Business Park, 45 Harare Drive, Harare, Zimbabwe",
		"Office Complex Delta, Building C, Floor 5, 89 Minsk Boulevard, Minsk 220114, Belarus",
		"Unit 908, Golden Trade Center, 78 Yangon Port Road, Yangon 11181, Myanmar",
		"Floor 3, Al-Zawra Tower, Block 215, Baghdad Commercial District, Baghdad, Iraq",
		"Building 45, Industrial Zone 3, Damascus International Airport Road, Damascus, Syria",
		"Suite 301, Havana Trade Building, 67 Malecon Avenue, Havana 10400, Cuba",
		"Office 12, Floor 4, Conakry Commerce Center, Route du Niger, Conakry, Guinea",
		"Unit 55, Khartoum Business Complex, Al Gamhoria Avenue, Khartoum, Sudan",
		"Room 789, Floor 7, Trade Tower 3, Kim Il Sung Square, Pyongyang, DPRK",
		"Building 23, Floor 2, Sevastopol Maritime Complex, 45 Port Street, Sevastopol 99011",
		"Office 445, Tripoli Trade Center, Omar Al-Mukhtar Street, Tripoli, Libya",
	}

	models := []Model{
		// ModelFastText,
		// ModelFastTextAttention,
		// ModelFastTextLight,
		ModelBPEmb,
		ModelBPEmbAttention,
	}

	for _, m := range models {
		b.Run(string(m), func(b *testing.B) {
			b.StartTimer()
			resp, err := cc.ParseAddresses(ctx, m, inputs)
			b.StopTimer()

			require.NoError(b, err)
			require.Len(b, resp.Addresses, len(inputs)) // we got back the same number of addresses
		})
	}
}
