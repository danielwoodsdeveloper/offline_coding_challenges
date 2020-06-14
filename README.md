# Offline Coding Challenges

Complete coding challenges. Anywhere. Anytime.

## Status

*Still being developed. Tests still being written, language support still being added.*

**Need to add unit tests to both server and client.**

## Building and Running

### Docker

**Offline Coding Challenges** is designed be run in Docker. The Dockerfile is in the root directory. To build the image (tagged as "offline-coding-challenges"):

```docker build . -t offline-coding-challengs```

To run the image we've just built:

```docker run -p 8080:8080 -v /var/run/docker.sock:/var/run/docker.sock -v cc-data:/app/temp offline-coding-challengs```

### Server (Go)

If you're working on the server, it can be run without Docker. First, make sure you've loaded all the dependencies via `go get`. Then, simply:

```go run main.go```

The server will run at port *8080*, so make sure you have it free.

### Client (React)

You'll need to make sure the server is running first at port *8080*. Without it, the client is somewhat useless.

You'll need to download all the node modules first:

```npm install && npm install react-scripts```

And then you can run the client:

```yarn start```

The dev server will run at port *3000*.

## Design

**Offline Coding Challenges** is based on Docker. The code you submit is run in temporary containers so you don't need to worry about having anything (besides Docker) installed.

The **Offline Coding Challenges** client and server is also designed to run inside a container, using *Docker On Docker*. The container it runs inside is designed to be linked to the host's Docker engine, so the Docker commands it executs are actually executed on the host. For example, when you insall a runtime, the image that is pulled is actually pulled on the host.

The server is built in Go, atop Mux. It exposes a ReST API that allows the client to interface with it. Go was picked for a few reasons:
* Its concurrency model is simple, meaning it was easy to handle multiple concurrent code submissions easily.
* Docker has an SDK built for Go, meaning the engine was easy to interface with.
* Go is quick. Very quick.
* Previous familiarity with the language.

The client is built with React. It interfaces with the server via the aforementioned ReST API. It uses Bootstrap (React Bootstrap) because it's a responsive, easy to work with framework.

## Challenges

The largest complexity came from the *Docker On Docker* design. To enable concurrency, we need to be able to pass the submitted code to several containers at the same time. This obviously required mounting a volume, and initially it seemed to make sense to use a bind mount. However, because the **Offline Coding Challenges** client/server container has been bound to the host's Docker engine, when we try to bind a directory onto our temporary code submission container, the bind lookup actually happens on the host machine, where the submitted code doesn't exist and we can't readily control the file system (nor do we want to).

Instead, we can use a named volume. Using a named volume means we can mount the same volume to the client/server container and the temporary code submission containers, and here we can readily control the file system.

## Contributing

There's still plenty of work to do, so help is more than welcome. This can be as simple as adding a new test, a new test case, or a new runtime. Or you could be improving the design or fixing bugs.

At this stage, the only process I desire is creating a feature branch and raising a PR to master with details. As (if) this grows, this process may change.