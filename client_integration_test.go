package deepparsego

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIntegrationParseAddresses(t *testing.T) {
	base := liveDeepparseURL(t)
	httpClient := &http.Client{Timeout: 30 * time.Second}
	cc := NewClient(httpClient, base)

	inputs := []string{
		"350 rue des Lilas Ouest Quebec city Quebec G1L 1B6",
		"2325 Rue de l'Université, Québec, QC G1V 0A6",
		"350 rue des Lilas Ouest Quebec city Quebec G1L 1B6",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := cc.ParseAddresses(ctx, ModelBPEmbAttention, inputs)
	require.NoError(t, err)
	require.Equal(t, Model("bpemb_attention"), resp.Model)
	require.NotEmpty(t, resp.Version)
	require.Len(t, resp.Addresses, len(inputs))

	require.Equal(t, inputs[0], resp.Addresses[0].Raw)
	require.Equal(t, inputs[1], resp.Addresses[1].Raw)
	require.Equal(t, inputs[2], resp.Addresses[2].Raw)

	require.Equal(t, "350", resp.Addresses[0].StreetNumber)
	require.Equal(t, "2325", resp.Addresses[1].StreetNumber)
	require.Equal(t, "350", resp.Addresses[2].StreetNumber)
	require.NotEmpty(t, resp.Addresses[0].PostalCode)
	require.NotEmpty(t, resp.Addresses[1].PostalCode)
}

func TestIntegrationParseAddresses_EmptyRejected(t *testing.T) {
	base := liveDeepparseURL(t)
	cc := NewClient(&http.Client{Timeout: 15 * time.Second}, base)

	_, err := cc.ParseAddresses(context.Background(), ModelBPEmb, nil)
	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, http.StatusUnprocessableEntity, apiErr.StatusCode)
}

func TestIntegrationParseAddresses_UnknownModelRejected(t *testing.T) {
	base := liveDeepparseURL(t)
	cc := NewClient(&http.Client{Timeout: 15 * time.Second}, base)

	_, err := cc.ParseAddresses(context.Background(), Model("not-a-model"), []string{"350 rue des Lilas Ouest Quebec city Quebec G1L 1B6"})
	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, http.StatusUnprocessableEntity, apiErr.StatusCode)
}

func BenchmarkClient(b *testing.B) {
	base := liveDeepparseURL(b)
	ctx := context.Background()

	httpClient := &http.Client{Timeout: 30 * time.Second}
	cc := NewClient(httpClient, base)

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
		ModelBPEmb,
		ModelBPEmbAttention,
	}

	for _, m := range models {
		b.Run(string(m), func(b *testing.B) {
			b.StartTimer()
			resp, err := cc.ParseAddresses(ctx, m, inputs)
			b.StopTimer()

			require.NoError(b, err)
			require.Len(b, resp.Addresses, len(inputs))
		})
	}
}

func liveDeepparseURL(t testing.TB) string {
	t.Helper()

	base := os.Getenv("DEEPPARSE_URL")
	requireServer := base != ""
	if base == "" {
		base = "http://localhost:8000"
	}

	if testing.Short() && !requireServer {
		t.Skip("skipping deepparse docker integration in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base+"/docs", nil)
	if err != nil {
		if requireServer {
			t.Fatalf("deepparse server not running at %s: %v", base, err)
		}
		t.Skipf("deepparse server not running at %s: %v", base, err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if requireServer {
			t.Fatalf("deepparse server not running at %s: %v", base, err)
		}
		t.Skipf("deepparse server not running at %s: %v", base, err)
	}
	resp.Body.Close()
	return base
}
