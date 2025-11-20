{{ define "iprange" -}}
package generated

var (
  {{ .Variable }} = unpackEntries(`{{ .Entries }}`)
)
{{ end }}
