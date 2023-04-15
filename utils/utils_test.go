package utils

import "testing"

func TestShortString(t *testing.T) {

	testCases := []struct {
		input          string
		max            int
		expectedOutput string
	}{
		{input: "Hello World", max: 20, expectedOutput: "Hello World"},
		{input: "0123456789abcde zxvjlk asd lkzjxc", max: 10, expectedOutput: "0123456..."},
		{input: "0123456789", max: 10, expectedOutput: "0123456789"},
	}

	for _, testCase := range testCases {
		actualOutput := ShortString(testCase.input, testCase.max)
		if actualOutput != testCase.expectedOutput {
			t.Errorf("ShortString(%v, %v) = %v; want %v", testCase.input, testCase.max, actualOutput, testCase.expectedOutput)
		}
	}
}
