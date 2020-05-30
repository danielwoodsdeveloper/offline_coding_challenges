package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtimes"
	"strings"
	"time"

	"github.com/danielwoodsdeveloper/offline_coding_challenges/server/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
)

func Run(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	out, err := cli.ImagePull(ctx, "alpine", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	defer out.Close()

	io.Copy(os.Stdout, out)
}

func main() {
	router := mux.NewRouter()
	
    router.HandleFunc("/run/{language}", Run)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	
    http.ListenAndServe(":8080", router)
}