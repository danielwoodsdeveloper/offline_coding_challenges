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
		Description: `<p>The starting place of all coding journeys.</p>
		<p>Write some code that will take a command line argument,
		and print "<i>Hello, <b>arg1</b>!</i>" to the console.</p>
		<p>For example, if the argument is "World", it will print
		"<i>Hello, World!</i>"</p>`,
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
		Description: `<p>Write some code that will print the numbers <b>1</b>
		to <b>20</b> to the console.</p>
		<p>Where the number is divisible by 3,
		print <b>arg1</b>. Where divisible by 5, also print <b>arg2</b>.</p>`,
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

// Get a single test
func GetTestByNumber(num int) (Test, error) {
	for _, test := range tests {
		if test.Number == num {
			return test, nil
		}
	}

	return Test{}, errors.New("Test not found")
}

// Get all the tests
func GetAllTests() ([]Test) {
	return tests
}