package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"net/netip"
	"os"
	"slices"
	"strings"
	"text/template"

	"github.com/yandzee/go-ipdb/generated"
	"github.com/yandzee/go-ipdb/internal/types"
)

type TemplateEntry struct {
	RangeStart  string
	RangeEnd    string
	CountryCode string
}

type TemplateVariables struct {
	Variable string
	Entries  string
}

func main() {
	fname := os.Getenv("CSV_FILE")
	if len(fname) == 0 {
		fmt.Fprintf(os.Stderr, "CSV_FILE env is not set\n")
		os.Exit(1)
	}

	tpl, err := template.ParseFiles("../internal/templates/iprange.go.tpl")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse iprange.go.tpl: %v\n", err)
		os.Exit(1)
	}

	v4file, err := os.OpenFile("../generated/v4.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open v4.go: %v\n", err)
		os.Exit(1)
	}

	v6file, err := os.OpenFile("../generated/v6.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open v6.go: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("About to parse %s\n", fname)

	file, err := os.Open(fname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		_ = file.Close()
	}()

	reader := bufio.NewReader(file)
	lines := bufio.NewScanner(reader)

	v4entries := []generated.AddrRangeCountry{}
	v6entries := []generated.AddrRangeCountry{}

	for i := 0; lines.Scan(); i += 1 {
		line := lines.Text()
		parts := strings.Split(line, ",")

		if len(parts) != 3 || parts[2] == "ZZ" {
			fmt.Printf("Line %d \"%s\": skip\n", i, line)
			continue
		}

		isv4, entry, err := processLine(parts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Line %d err: %v\n", i, err)
			os.Exit(1)
		}

		if isv4 {
			v4entries = append(v4entries, entry)
		} else {
			v6entries = append(v6entries, entry)
		}
	}

	if err := lines.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Lines iterator err: %v\n", err)
		os.Exit(1)
	}

	v4out := bufio.NewWriter(v4file)
	v6out := bufio.NewWriter(v6file)

	b64, err := entriesToB64(v4entries)
	if err != nil {
		fmt.Fprintf(os.Stderr, "v4 entriesToB64 err: %v\n", err)
		os.Exit(1)
	}

	_ = tpl.ExecuteTemplate(v4out, "iprange", TemplateVariables{
		Variable: "V4Entries",
		Entries:  b64,
	})

	b64, err = entriesToB64(v6entries)
	if err != nil {
		fmt.Fprintf(os.Stderr, "v6 entriesToB64 err: %v\n", err)
		os.Exit(1)
	}

	_ = tpl.ExecuteTemplate(v6out, "iprange", TemplateVariables{
		Variable: "V6Entries",
		Entries:  b64,
	})

	fmt.Printf("V4 base64 size: %d, packed entries: %d, err: %v\n", len(b64), len(v4entries), err)
	fmt.Printf("V6 base64 size: %d, packed entries: %d, err: %v\n", len(b64), len(v6entries), err)

	_ = v4out.Flush()
	_ = v6out.Flush()
}

func entriesToB64(entries []generated.AddrRangeCountry) (string, error) {
	sortRanges(&entries)

	buf := bytes.Buffer{}
	compressor, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return "", err
	}

	enc := gob.NewEncoder(compressor)

	for _, entry := range entries {
		if err := enc.Encode(entry); err != nil {
			return "", err
		}
	}

	if err := compressor.Flush(); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func sortRanges(arr *[]generated.AddrRangeCountry) {
	slices.SortFunc(*arr, func(a, b generated.AddrRangeCountry) int {
		return a.RangeStart.Compare(b.RangeStart)
	})
}

func processLine(parts []string) (bool, generated.AddrRangeCountry, error) {
	rangeStart, err := netip.ParseAddr(parts[0])
	if err != nil {
		return false, generated.AddrRangeCountry{}, err
	}

	rangeEnd, err := netip.ParseAddr(parts[1])
	if err != nil {
		return false, generated.AddrRangeCountry{}, err
	}

	addrRange := generated.AddrRangeCountry{
		RangeStart:  types.U128FromBEBytes(rangeStart.As16()),
		RangeEnd:    types.U128FromBEBytes(rangeEnd.As16()),
		CountryCode: parts[2],
	}

	return rangeStart.Is4(), addrRange, nil
}
