package runtimes

import (
	"errors"
)

type Runtime struct {
	Name string
	Image string
	Display string
	Commands []string
	FileName string
	Template []string
}

var runtimes = map[string]Runtime {
	"java7": Runtime {
		Name: "java7",
		Display: "Java 7",
		Image: "openjdk:7-alpine",
		Commands: []string{"/bin/sh", "-c", "cd /tmp && javac Submission.java && java Submission {INPUTS}"},
		FileName: "Submission.java",
		Template: []string{
			"public class Submission {",
			"\tpublic static void main(String[] args) {",
			"\t\tSystem.out.println(args[0]);",
			"\t}",
			"}",
		},
	},
	"java8": Runtime {
		Name: "java8",
		Display: "Java 8",
		Image: "openjdk:8-alpine",
		Commands: []string{"/bin/sh", "-c", "cd /tmp && javac Submission.java && java Submission {INPUTS}"},
		FileName: "Submission.java",
		Template: []string{
			"public class Submission {",
			"\tpublic static void main(String[] args) {",
			"\t\tSystem.out.println(args[0]);",
			"\t}",
			"}",
		},
	},
	"golang": Runtime {
		Name: "golang",
		Display: "Go",
		Image: "golang:alpine",
		Commands: []string{"/bin/sh", "-c", "cd /tmp && go run submission.go {INPUTS}"},
		FileName: "submission.go",
		Template: []string{
			"package main",
			"",
			"import (",
			"\t\"fmt\"",
			"\t\"os\"",
			")",
			"",
			"func main() {",
			"\targs := os.Args[1:]",
			"\tfmt.Println(args[0])",
			"}",
		},
	},
	"python": Runtime {
		Name: "python",
		Display: "Python",
		Image: "python:alpine",
		Commands: []string{"/bin/sh", "-c", "cd /tmp && python submission.py {INPUTS}"},
		FileName: "submission.py",
		Template: []string{
			"import sys",
			"",
			"args = sys.argv[1:]",
			"print(args[0])",
		},
	},
	"javascript": Runtime {
		Name: "javascript",
		Display: "Node JS",
		Image: "node:alpine",
		Commands: []string{"/bin/sh", "-c", "cd /tmp && node submission.js {INPUTS}"},
		FileName: "submission.js",
		Template: []string{
			"var args = process.argv.slice(1)",
			"console.log(args[0])",
		},
	},
}

func GetRuntime(name string) (Runtime, error) {
	val, found := runtimes[name]
	
	if !found {
		return Runtime{}, errors.New("Could not find a matching runtime")
	}

	return val, nil
}

func GetAllRuntimes() ([]Runtime) {
	res := []Runtime{}
	for _, val := range runtimes {
		res = append(res, val)
	}

	return res
}