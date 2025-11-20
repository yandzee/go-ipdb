//go:generate go run ../cmd/generate.go
package generated

import (
	"github.com/yandzee/go-ipdb/internal/types"
)

type AddrRangeCountry struct {
	RangeStart  types.Uint128
	RangeEnd    types.Uint128
	CountryCode string
}
