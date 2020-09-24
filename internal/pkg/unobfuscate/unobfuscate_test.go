package unobfuscate

import (
	"testing"
)

func TestUnobfuscate(t *testing.T) {
	test := "testATandrewDOTcmuDOTedu"
	expected := "test@andrew.cmu.edu"
	actual := Unobfuscate(test)

	if actual != expected {
		t.Errorf("Expected %v to translate to %v but instead got %v!", test, expected, actual)
	}
}
