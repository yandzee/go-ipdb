package ipdb

import (
	"net/netip"
	"slices"

	"github.com/yandzee/go-ipdb/generated"
	"github.com/yandzee/go-ipdb/internal/types"
)

func LookupString(str string) (string, error) {
	addr, err := netip.ParseAddr(str)
	if err != nil {
		return "", err
	}

	return LookupAddr(addr), nil
}

func LookupAddr(addr netip.Addr) string {
	num := types.U128FromBEBytes(addr.As16())

	switch {
	case addr.Is4():
		return bsearch(generated.V4Entries, num)
	case addr.Is6():
		return bsearch(generated.V6Entries, num)
	}

	return ""
}

func bsearch(entries []generated.AddrRangeCountry, num types.Uint128) string {
	idx, found := slices.BinarySearchFunc(entries, num, func(e generated.AddrRangeCountry, n types.Uint128) int {
		rs, re := e.RangeStart.Compare(n), e.RangeEnd.Compare(n)

		if rs <= 0 && re >= 0 {
			return 0
		}

		return rs
	})

	if !found {
		return ""
	}

	return entries[idx].CountryCode
}
