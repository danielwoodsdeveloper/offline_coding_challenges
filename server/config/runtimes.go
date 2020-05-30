package config

type runtime struct {
	image string
	command string
}

runtimes := map[string]runtime {
	"java8": runtime {
		image: "openjdk:8-alpine",
		command: "javac Hello.java && java Hello"
	},
	"golang": runtime {
		image: "golang:alpine",
		command: "run main.go"
	}
}

func runtime GetRuntime(string) {
	return runtimes[string]
}