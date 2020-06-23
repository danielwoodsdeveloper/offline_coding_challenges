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

	// Beginner > Most Letters
	Test {
		Number: 3,
		Group: "Beginner",
		Title: "Most Letters",
		Description: `<p>Write some code that will sort the input arguments by
		the number of letters in those arguments. Print these arguments in descending
		order to the console, each on a new line.</p>
		<p>This can take any number of arguments.<p>`,
		TestCases: []TestCase {
			TestCase {
				Number: 1,
				Inputs: []string { "Homer", "Bart", "Moe", "Hibbert" },
				ExpectedOutput: []string { "Hibbert", "Homer", "Bart", "Moe" },
			},
			TestCase {
				Number: 2,
				Inputs: []string { "AB123456789", "ABC123", "ABC123DEF" },
				ExpectedOutput: []string { "ABC123DEF", "ABC123", "AB123456789" },
			},
		},
	},

	// Intermediate > Caesar Solver
	Test {
		Number: 4,
		Group: "Intermediate",
		Title: "Caesar Solver",
		Description: `<p>Wrtie some code that will take 2 arguments: firstly a positive
		integer <i>N</i> which is the cipher key, and a ciphertext string. The code
		should solve the Caesar cipher using <i>N</i>, and print the plaintext to the
		console.</p>`,
		TestCases: []TestCase {
			TestCase {
				Number: 1,
				Inputs: []string { "7", "\"Olssv, dvysk!\"" },
				ExpectedOutput: []string { "Hello, world!" },
			},
			TestCase {
				Number: 2,
				Inputs: []string { "3", "\"Wkh zdb wr jhw vwduwhg lv wr txlw wdonlqj dqg ehjlq grlqj.\"" },
				ExpectedOutput: []string { "The way to get started is to quit talking and begin doing." },
			},
			TestCase {
				Number: 3,
				Inputs: []string { "15", "\"Axut xh lwpi wpeetch lwtc ndj'gt qjhn bpzxcv diwtg eapch.\"" },
				ExpectedOutput: []string { "Life is what happens when you're busy making other plans." },
			},
			TestCase {
				Number: 4,
				Inputs: []string { "26", "\"Try not to become a man of success. Rather become a man of value.\"" },
				ExpectedOutput: []string { "Try not to become a man of success. Rather become a man of value." },
			},
		},
	},

	// Intermediate > Candy Shop
	Test {
		Number: 5,
		Group: "Intermediate",
		Title: "Candy Shop",
		Description: `<p>You have some spare change, and you want to spend it at
		a lolly shop. You want to make sure you get the best value for money, and
		you want to try as many different lollies as you can.</p>
		<p><b>arg1</b> will be the amount of money you have in cents, as an integer,
		for example 125.</p>
		<p>The following arguments will come in pairs, the first being the lolly name
		as a string, the second the lolly's price as an integer in cents. There can be any
		number of these. For example, <b>arg2</b> could be "Milk Bottles" and <b>arg3</b>
		35.</p>
		<p>The output of this code will be two integers, the first being the number of
		lolly types you can buy, and the second the amount of money you have leftover.</p>`,
		TestCases: []TestCase {
			TestCase {
				Number: 1,
				Inputs: []string { "100", "\"Gummi Bears\"", "75", "\"Clouds\"", "35", "\"Milk Duds\"", "50", "\"Freckles\"", "15", "\"Lollipop\"", "85", "\"Jelly Babies\"", "15", "\"Sour Worms\"", "80", "\"Gummi Snakes\"", "25", "\"Minties\"", "45", "\"Red Frogs\"", "50", "\"Green Frogs\"", "5" },
				ExpectedOutput: []string { "5", "5" },
			},
			TestCase {
				Number: 2,
				Inputs: []string { "120", "\"Gummi Bears\"", "75", "\"Clouds\"", "35", "\"Milk Duds\"", "50", "\"Freckles\"", "15", "\"Lollipop\"", "85", "\"Jelly Babies\"", "15", "\"Sour Worms\"", "80", "\"Gummi Snakes\"", "25", "\"Minties\"", "45", "\"Red Frogs\"", "50", "\"Green Frogs\"", "5" },
				ExpectedOutput: []string { "5", "25" },
			},
			TestCase {
				Number: 3,
				Inputs: []string { "250", "\"Gummi Bears\"", "75", "\"Clouds\"", "35", "\"Milk Duds\"", "50", "\"Freckles\"", "15", "\"Lollipop\"", "85", "\"Jelly Babies\"", "15", "\"Sour Worms\"", "80", "\"Gummi Snakes\"", "25", "\"Minties\"", "45", "\"Red Frogs\"", "50", "\"Green Frogs\"", "5" },
				ExpectedOutput: []string { "8", "10" },
			},
			TestCase {
				Number: 4,
				Inputs: []string { "500", "\"Gummi Bears\"", "75", "\"Clouds\"", "35", "\"Milk Duds\"", "50", "\"Freckles\"", "15", "\"Lollipop\"", "85", "\"Jelly Babies\"", "15", "\"Sour Worms\"", "80", "\"Gummi Snakes\"", "25", "\"Minties\"", "45", "\"Red Frogs\"", "50", "\"Green Frogs\"", "5" },
				ExpectedOutput: []string { "11", "20" },
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