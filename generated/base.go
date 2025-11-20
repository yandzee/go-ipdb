//go:generate go run ../cmd/generate.go
package generated

import (
	"compress/gzip"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/yandzee/go-ipdb/internal/types"
)

type AddrRangeCountry struct {
	RangeStart  types.Uint128
	RangeEnd    types.Uint128
	CountryCode string
}

func unpackEntries(b64 string) []AddrRangeCountry {
	unb64 := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64))
	// compressed, err := base64.StdEncoding.DecodeString(b64)
	// if err != nil {
	// 	panic("Failed to decode base64 of entries: " + err.Error())
	// }

	decompressor, err := gzip.NewReader(unb64)
	if err != nil {
		panic("Failed to create gzip reader: " + err.Error())
	}

	dec := gob.NewDecoder(decompressor)

	var e AddrRangeCountry
	entries := []AddrRangeCountry{}

	for {
		err := dec.Decode(&e)

		if err == nil {
			entries = append(entries, e)
			continue
		}

		// NOTE: Feels like this is a bug in gob Decoder
		if errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) {
			break
		}

		panic("Decode failed: " + err.Error())
	}

	fmt.Printf("unpacked entries: %v\n", len(entries))
	return entries
}
