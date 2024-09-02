package helpers

import (
	"testing"
)

const ansiCodesLength = 48

func testText(t *testing.T, testFunc func(string) string, input string, expectedLength int) {
	result := len(testFunc(input))
	if result != expectedLength {
		t.Errorf("len(%T(%q)) = %d; want %d", testFunc, input, result, expectedLength)
	}
}

func TestNormalText(t *testing.T) {
	input := "test"
	expectedLength := ansiCodesLength + len(input)
	testText(t, NormalText, input, expectedLength)
}

func TestBoldText(t *testing.T) {
	input := "test"
	expectedLength := ansiCodesLength + len(input)
	testText(t, BoldText, input, expectedLength)
}

func TestGreenText(t *testing.T) {
	input := "test"
	expectedLength := ansiCodesLength + len(input)
	testText(t, GreenText, input, expectedLength)
}
