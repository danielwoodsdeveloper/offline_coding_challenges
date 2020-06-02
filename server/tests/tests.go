package tests

import (
	"errors"
)

type Test struct {
	Number int `json:"number"`
	Title string `json:"title"`
	Description string `json:"description"`
	TestCases []TestCase `json:"test-cases"`
}

type TestCase struct {
	Number int `json:"number"`
	ExpectedOutput []string `json:"expected-output"`
}

var tests = []Test {
	Test {
		Number: 1,
		Title: "Hello, World!",
		Description: "The starting place of all coding journeys.",
		TestCases: []TestCase {
			TestCase {
				Number: 1,
				ExpectedOutput: []string {
					"Hello, world!",
				},
			},
			TestCase {
				Number: 2,
				ExpectedOutput: []string {
					"Hello, waldo!",
				},
			},
		},
	},
}

func GetTestByNumber(num int) (Test, error) {
	for _, test := range tests {
		if test.Number == num {
			return test, nil
		}
	}

	return Test{}, errors.New("Test not found")
}