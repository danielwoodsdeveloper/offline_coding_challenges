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
	"java7": Runtime {
		Image: "openjdk:7-alpine",
		Commands: []string{"/bin/sh", "-c", "cd /tmp && javac Submission.java && java Submission"},
		FileName: "Submission.java",
	},
	"java8": Runtime {
		Image: "openjdk:8-alpine",
		Commands: []string{"/bin/sh", "-c", "cd /tmp && javac Submission.java && java Submission"},
		FileName: "Submission.java",
	},
	"golang": Runtime {
		Image: "golang:alpine",
		Commands: []string{"/bin/sh", "-c", "cd /tmp && go run submission.go"},
		FileName: "submission.go",
	},
	"python": Runtime {
		Image: "python:alpine",
		Commands: []string{"/bin/sh", "-c", "cd /tmp && python submission.py"},
		FileName: "submission.py",
	},
}

func GetRuntime(name string) (Runtime, error) {
	val, found := runtimes[name]
	
	if !found {
		return Runtime{}, errors.New("Could not find a matching runtime")
	}

	return val, nil
}