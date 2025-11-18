//go:generate go run ../cmd/generate.go
package generated

import (
	"net/netip"
)

type AddrRangeCountry struct {
	RangeStart  netip.Addr
	RangeEnd    netip.Addr
	CountryCode string
}
