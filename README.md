# PeeringDB API - Go package

[![Go Reference](https://pkg.go.dev/badge/github.com/gmazoyer/peeringdb.svg)](https://pkg.go.dev/github.com/gmazoyer/peeringdb)
[![Go Report Card](https://goreportcard.com/badge/github.com/gmazoyer/peeringdb)](https://goreportcard.com/report/github.com/gmazoyer/peeringdb)

This is a Go package that allows developer to interact with the
[PeeringDB API](https://peeringdb.com/apidocs/) in the easiest way possible.
There are no binaries provided with this package. It can only be used as a
library.

## Installation

Install the library package with `go get github.com/gmazoyer/peeringdb`.

## Example

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gmazoyer/peeringdb"
)

func main() {
	// Create an API client (optionally with an API key)
	api := peeringdb.NewAPI(
		peeringdb.WithAPIKey("your-api-key"),
	)
	ctx := context.Background()

	// Look up a network by its ASN
	network, err := api.GetASN(ctx, 201281)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Network: %s (AS%d)\n", network.Name, network.ASN)

	// List all facilities linked to this network
	for _, netfacID := range network.NetworkFacilitySet {
		netfac, err := api.GetNetworkFacilityByID(ctx, netfacID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  Facility: %s (%s, %s)\n", netfac.Name, netfac.City, netfac.Country)
	}
}
```

More examples are available in the
[package documentation](https://pkg.go.dev/github.com/gmazoyer/peeringdb).

You can also find a real life example with the
[PeeringDB synchronization tool](https://github.com/gmazoyer/peeringdb-sync).
