{{ define "iprange" -}}
package generated

import "net/netip"

var (
  {{ .Variable }} = []AddrRangeCountry{
    {{ range $e := .Entries }}
    {
      RangeStart: netip.MustParseAddr("{{ $e.RangeStart }}"),
      RangeEnd: netip.MustParseAddr("{{ $e.RangeEnd }}"),
      CountryCode: "{{ $e.CountryCode }}",
    },
    {{- end }}
  }
)
{{ end }}
