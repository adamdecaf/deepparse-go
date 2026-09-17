# deepparse-go

[![Go Reference](https://pkg.go.dev/badge/github.com/adamdecaf/deepparse-go.svg)](https://pkg.go.dev/github.com/adamdecaf/deepparse-go)
[![CI](https://github.com/adamdecaf/deepparse-go/actions/workflows/ci.yml/badge.svg)](https://github.com/adamdecaf/deepparse-go/actions/workflows/ci.yml)
[![Coverage Status](https://codecov.io/gh/adamdecaf/deepparse-go/branch/master/graph/badge.svg)](https://codecov.io/gh/adamdecaf/deepparse-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/adamdecaf/deepparse-go)](https://goreportcard.com/report/github.com/adamdecaf/deepparse-go)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/adamdecaf/deepparse-go/master/LICENSE)

Go client for [GRAAL-Research/deepparse](https://github.com/GRAAL-Research/deepparse)'s [HTTP API](https://deepparse.org/api.html) (0.11.0+).

The parse endpoint returns `parsed_addresses` as a list of `{raw: parsed}` objects so duplicate inputs keep their order.

## Usage

Import the Go library and create a client to start parsing addresses. Requires a running instance of deepparse.

```go
httpClient := &http.Client{
	Timeout: 5 * time.Second,
}
cc := NewClient(httpClient, "http://localhost:8000")

ctx := context.Background()
resp, err := cc.ParseAddresses(ctx, ModelBPEmbAttention, []string{
	"350 rue des Lilas Ouest Quebec city Quebec G1L 1B6",
	"2325 Rue de l'Université, Québec, QC G1V 0A6",
})

// handle resp.Addresses
```

Supported models: `fasttext`, `fasttext-attention`, `fasttext-light`, `bpemb`, `bpemb-attention`.
The published Docker image loads the BPEmb models; FastText models are skipped at startup.

## Contributing

The integration tests require pulling the [ghcr.io/graal-research/deepparse](https://github.com/GRAAL-Research/deepparse/pkgs/container/deepparse) image (`0.11.0`).

Pull and start the container before running the full Go tests:

```
$ make docker-up   # takes a while to download and load models
$ go test ./...
ok  	github.com/adamdecaf/deepparse-go	0.773s
```

Unit tests do not need Docker:

```
$ go test -short ./...
```

## License

MIT Licensed
