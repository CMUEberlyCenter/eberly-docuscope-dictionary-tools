package fix

import (
	"strings"
)

/**
 * Corrects letter case for words and wordclasses.
 * Words should be lowercase.
 * Wordclasses, indicated by ! prefix, should be uppercase.
 */
func Case(pat []string) []string {
	ret := make([]string, len(pat))
	for i, v := range pat {
		if strings.HasPrefix(v, "!") {
			ret[i] = strings.ToUpper(v)
		} else {
			ret[i] = strings.ToLower(v)
		}
	}
	return ret
}
