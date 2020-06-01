package config

import (
	"errors"
)

type Runtime struct {
	Image string
	Commands []string
	FileName string
}

var runtimes = map[string]Runtime {
	"java8": Runtime {
		Image: "openjdk:8-alpine",
		Commands: []string{"javac submission.java && java Submission"},
		FileName: "submission.java",
	},
	"golang": Runtime {
		Image: "golang:alpine",
		Commands: []string{"go", "run", "/tmp/submission.go"},
		FileName: "submission.go",
	},
}

func GetRuntime(name string) (Runtime, error) {
	val, found := runtimes[name]
	
	if !found {
		return Runtime{}, errors.New("Could not find a matching runtime")
	}

	return val, nil
}