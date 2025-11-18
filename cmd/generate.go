package main

import (
	"bufio"
	"fmt"
	"html/template"
	"net/netip"
	"os"
	"slices"
	"strings"

	"github.com/yandzee/go-ipdb/generated"
	"github.com/yandzee/go-ipdb/internal/types"
)

type Stash struct {
	v4 []generated.AddrRangeCountry
	v6 []generated.AddrRangeCountry
}

type TemplateEntry struct {
	RangeStart  string
	RangeEnd    string
	CountryCode string
}

type TemplateVariables struct {
	Variable string
	Entries  []TemplateEntry
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
	stash := &Stash{}

	for i := 0; lines.Scan(); i += 1 {
		line := lines.Text()
		parts := strings.Split(line, ",")

		if len(parts) != 3 || parts[2] == "ZZ" {
			fmt.Printf("Line %d \"%s\": skip\n", i, line)
			continue
		}

		if err := processLine(stash, parts); err != nil {
			fmt.Fprintf(os.Stderr, "Line %d err: %v\n", i, err)
			os.Exit(1)
		}
	}

	if err := lines.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Lines iterator err: %v\n", err)
		os.Exit(1)
	}

	sortRanges(&stash.v4)
	sortRanges(&stash.v6)

	v4vars := TemplateVariables{
		Variable: "V4Ranges",
	}

	v6vars := TemplateVariables{
		Variable: "V6Ranges",
	}

	for _, r := range stash.v4 {
		v4vars.Entries = append(v4vars.Entries, TemplateEntry{
			RangeStart:  r.RangeStart.String(),
			RangeEnd:    r.RangeEnd.String(),
			CountryCode: r.CountryCode,
		})
	}

	for _, r := range stash.v6 {
		v6vars.Entries = append(v6vars.Entries, TemplateEntry{
			RangeStart:  r.RangeStart.String(),
			RangeEnd:    r.RangeEnd.String(),
			CountryCode: r.CountryCode,
		})
	}

	v4out := bufio.NewWriter(v4file)
	v6out := bufio.NewWriter(v6file)

	defer v4out.Flush()
	defer v6out.Flush()

	_ = tpl.ExecuteTemplate(v4out, "iprange", v4vars)
	_ = tpl.ExecuteTemplate(v6out, "iprange", v6vars)
}

func sortRanges(arr *[]generated.AddrRangeCountry) {
	slices.SortFunc(*arr, func(a, b generated.AddrRangeCountry) int {
		anum := types.U128FromBEBytes(a.RangeStart.As16())
		bnum := types.U128FromBEBytes(b.RangeStart.As16())

		return anum.Compare(bnum)
	})
}

func processLine(stash *Stash, parts []string) error {
	rangeStart, err := netip.ParseAddr(parts[0])
	if err != nil {
		return err
	}

	rangeEnd, err := netip.ParseAddr(parts[1])
	if err != nil {
		return err
	}

	addrRange := generated.AddrRangeCountry{
		RangeStart:  rangeStart,
		RangeEnd:    rangeEnd,
		CountryCode: parts[2],
	}

	switch {
	case rangeStart.Is4():
		stash.v4 = append(stash.v4, addrRange)
	case rangeStart.Is6():
		stash.v6 = append(stash.v6, addrRange)
	}

	return nil
}
