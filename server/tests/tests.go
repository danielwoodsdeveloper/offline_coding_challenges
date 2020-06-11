package tests

import (
	"errors"
)

type Test struct {
	Number int `json:"number"`
	Group string `json:"group"`
	Title string `json:"title"`
	Description string `json:"description"`
	TestCases []TestCase `json:"test-cases"`
}

type TestCase struct {
	Number int `json:"number"`
	Inputs []string `json:"inputs"`
	ExpectedOutput []string `json:"expected-output"`
}

var tests = []Test {
	// Beginner > Hello, World!
	Test {
		Number: 1,
		Group: "Beginner",
		Title: "Hello, World!",
		Description: "<b>The starting place of all coding journeys.</b>",
		TestCases: []TestCase {
			TestCase {
				Number: 1,
				Inputs: []string { "World" },
				ExpectedOutput: []string {
					"Hello, World!",
				},
			},
			TestCase {
				Number: 2,
				Inputs: []string { "Alfie" },
				ExpectedOutput: []string {
					"Hello, Alfie!",
				},
			},
			TestCase {
				Number: 3,
				Inputs: []string { "Turing" },
				ExpectedOutput: []string {
					"Hello, Turing!",
				},
			},
			TestCase {
				Number: 4,
				Inputs: []string { "Einstein" },
				ExpectedOutput: []string {
					"Hello, Einstein!",
				},
			},
		},
	},

	// Beginner > FizzBuzz
	Test {
		Number: 2,
		Group: "Beginner",
		Title: "FizzBuzz",
		Description: "TBD",
		TestCases: []TestCase {
			TestCase {
				Number: 1,
				Inputs: []string { "Fizz", "Buzz" },
				ExpectedOutput: []string {
					"1",
					"2",
					"Fizz",
					"4",
					"Buzz",
					"Fizz",
					"7",
					"8",
					"Fizz",
					"Buzz",
					"11",
					"Fizz",
					"13",
					"14",
					"FizzBuzz",
					"16",
					"17",
					"Fizz",
					"19",
					"Buzz",
				},
			},
			TestCase {
				Number: 2,
				Inputs: []string { "Fuzz", "Bizz" },
				ExpectedOutput: []string {
					"1",
					"2",
					"Fuzz",
					"4",
					"Bizz",
					"Fuzz",
					"7",
					"8",
					"Fuzz",
					"Bizz",
					"11",
					"Fuzz",
					"13",
					"14",
					"FuzzBizz",
					"16",
					"17",
					"Fuzz",
					"19",
					"Bizz",
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

func GetAllTests() ([]Test) {
	return tests
}