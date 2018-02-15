package pkg

import (
	"strings"

	"launchpad.net/xmlpath"
)

func KrDictID(lang, entryId, lexicalUnit string) string {
	return strings.Join([]string{"krdict", lang, entryId, lexicalUnit}, ":")
}

func XpathString(node *xmlpath.Node, xpath string) string {
	path := xmlpath.MustCompile(xpath)

	if value, ok := path.String(node); ok {
		return strings.TrimSpace(value)
	}

	return ""
}
