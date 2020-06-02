package main

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/danielwoodsdeveloper/offline_coding_challenges/server/config"
	"github.com/danielwoodsdeveloper/offline_coding_challenges/server/tests"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
	"github.com/sony/sonyflake"
)

type Submission struct {
	Runtime string `json:"runtime"`
	Code []string `json:"code"`
}

type TestCaseResponse struct {
	Number int `json:"number"`
	Pass bool `json:"pass"`
	Output string `json:"output"`
}

type SubmissionResponse struct {
	Pass bool `json:"pass"`
	TestCases []TestCaseResponse `json:"test-cases"`
}

var sf *sonyflake.Sonyflake

func Run(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := context.Background()

	// Decode request
	var sub Submission
	err := json.NewDecoder(r.Body).Decode(&sub)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch the test details
	num, err := strconv.Atoi(vars["test_number"])
	if err != nil {
		panic(err)
	}

	test, err := tests.GetTestByNumber(num)
	if err != nil {
		panic(err)
	}

	// Setup CLI
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// Fetch the runtime details
	runtime, err := config.GetRuntime(sub.Runtime)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	var tcResponses []TestCaseResponse

	for _, tc := range test.TestCases {
		wg.Add(1)

		go func(tc tests.TestCase) {
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

			// Create the container, mounting a volume to share the temp file
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

			// Start the container...
			err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
			if err != nil {
				panic(err)
			}

			// ...and wait for it to run, capturing all outputs...
			statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
			select {
			case err := <-errCh:
				if err != nil {
					panic(err)
				}
			case <-statusCh:
			}

			// ...then read it all...
			logReader, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
			if err != nil {
				panic(err)
			}
			defer logReader.Close()

			// ...and transform it into a string, stripping the first 8 header bytes
			stripped := make([]byte, 8)
			logReader.Read(stripped)
			logContent, _ := ioutil.ReadAll(logReader)

			err = os.RemoveAll("./" + strconv.Itoa(int(id)))
			if err != nil {
				panic(err)
			}

			tcr := TestCaseResponse{}
			tcr.Number = tc.Number
			tcr.Pass = strings.TrimSuffix(string(logContent), "\n") == strings.Join(tc.ExpectedOutput, "\n")
			tcr.Output = strings.TrimSuffix(string(logContent), "\n")

			tcResponses = append(tcResponses, tcr)

			wg.Done()
		}(tc)
	}

	wg.Wait()

	overallPass := true
	for _, tcr := range tcResponses {
		if tcr.Pass == false {
			overallPass = false
			break
		}
	}

	// Create our HTTP response
	res := SubmissionResponse{}
	res.Pass = overallPass
	res.TestCases = tcResponses

	json, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(json)
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
	
	router.HandleFunc("/tests/{test_number}/submission", Run).Methods("POST")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	
	http.ListenAndServe(":8080", router)
}