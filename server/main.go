package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/danielwoodsdeveloper/offline_coding_challenges/server/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gorilla/mux"
	"github.com/sony/sonyflake"
)

type Submission struct {
	Code []string `json:"code"`
}

var sf *sonyflake.Sonyflake

func Run(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var sub Submission
	err := json.NewDecoder(r.Body).Decode(&sub)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	// Setup CLI
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// Fetch the runtime details
	runtime, err := config.GetRuntime(vars["runtime"])
	if err != nil {
		panic(err)
	}

	// Get a UID
	id, err := sf.NextID()
	if err != nil {
		panic(err)
	}

	// Request code to byte array
	content := []byte(strings.Join(sub.Code, "\n"))

	// Code into temp file
	os.Mkdir(strconv.FormatUint(id, 10), 0755)
	file, err := os.Create("./" + strconv.FormatUint(id, 10) + "/" + runtime.FileName)
	if err != nil {
		panic(err)
	}
	file.Write(content)

	dir, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		panic(err)
	}

	// Pull image
	reader, err := cli.ImagePull(ctx, runtime.Image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	// Create the container
	resp, err := cli.ContainerCreate(ctx, &container.Config {
		Image: runtime.Image,
		Cmd: runtime.Commands,
	}, &container.HostConfig {
		Mounts: []mount.Mount {
			{
				Type: mount.TypeBind,
				Source: dir + "/" + strconv.Itoa(int(id)),
				Target: "/tmp",	
			},
		},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	// Run the container and get output
	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}
	defer out.Close()

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	err = os.RemoveAll("./" + strconv.Itoa(int(id)))
	if err != nil {
		panic(err)
	}
}

func init() {
	var st sonyflake.Settings
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		os.Exit(2)
	}
}

func main() {
	router := mux.NewRouter()
	
    router.HandleFunc("/run/{runtime}", Run)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	
    http.ListenAndServe(":8080", router)
}