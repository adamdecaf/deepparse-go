# deepparse-go

[![GoDoc](https://godoc.org/github.com/adamdecaf/deadcheck?status.svg)](https://godoc.org/github.com/adamdecaf/deepparse-go)
[![Build Status](https://github.com/adamdecaf/deepparse-go/workflows/Go/badge.svg)](https://github.com/adamdecaf/deepparse-go/actions)
[![Coverage Status](https://codecov.io/gh/adamdecaf/deepparse-go/branch/master/graph/badge.svg)](https://codecov.io/gh/adamdecaf/deepparse-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/adamdecaf/deepparse-go)](https://goreportcard.com/report/github.com/adamdecaf/deepparse-go)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/adamdecaf/deepparse-go/master/LICENSE)


Go client for [GRAAL-Research/deepparse](https://github.com/GRAAL-Research/deepparse)'s [HTTP API](https://deepparse.org/api.html).

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

## Contributing

The tests require pulling the [ghcr.io/graal-research/deepparse](https://github.com/GRAAL-Research/deepparse/pkgs/container/deepparse) image.

Pull and start the container before running Go tests:

```
$ docker compose up -d # takes a while to download and load data
```

```
$ go test ./...
ok  	github.com/adamdecaf/deepparse-go	0.773s
```

## License

MIT Licensed
