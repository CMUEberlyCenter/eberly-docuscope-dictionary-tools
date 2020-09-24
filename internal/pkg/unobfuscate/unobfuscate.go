package unobfuscate

import ("regexp")

/**
 * Quick and dirty function for undoing email obfuscation.
 */
func Unobfuscate(email string) string {
	reAT := regexp.MustCompile(`AT`)
	reDOT := regexp.MustCompile(`DOT`)
	out := reAT.ReplaceAllString(email, "@")
	out = reDOT.ReplaceAllString(out, ".")
	return out
}

func obfuscate(email string) string {
	reAT := regexp.MustCompile(`@`)
	reDOT := regexp.MustCompile(`\.`)
	out := reAT.ReplaceAllString(email, "AT")
	out = reDOT.ReplaceAllString(out, "DOT")
	return out
}
