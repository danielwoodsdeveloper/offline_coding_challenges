package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/danielwoodsdeveloper/offline_coding_challenges/server/runtimes"
	"github.com/danielwoodsdeveloper/offline_coding_challenges/server/tests"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/icza/gox/stringsx"
	"github.com/sony/sonyflake"
)

type SubmissionRequest struct {
	Runtime string `json:"runtime"`
	Code []string `json:"code"`
}

type TestCaseResponse struct {
	Number int `json:"number"`
	Pass bool `json:"pass"`
	Inputs []string `json:"inputs"`
	Outputs []string `json:"outputs"`
	Expected []string `json:"expected"`
}

type SubmissionResponse struct {
	Pass bool `json:"pass"`
	TestCases []TestCaseResponse `json:"testcases"`
}

type RuntimeDetailResponse struct {
	Name string `json:"name"`
	Display string `json:"display"`
	Image string `json:"image"`
	Installed bool `json:"installed"`
	Template string `json:"template"`
}

type TestDetailResponse struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Number int `json:"number"`
}

var sf *sonyflake.Sonyflake
var cli *client.Client

// Handles code submission (POST). Fetches the inputs from the relevant
// test case, spins up a container relative to the selected runtime
// passes in the code and runs it, fetching the output
func TestSubmission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Decode request
	var sub SubmissionRequest
	err := json.NewDecoder(r.Body).Decode(&sub)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Pull number from query params...
	num, err := strconv.Atoi(vars["test_number"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ...and fetch the details of the corresponding test
	test, err := tests.GetTestByNumber(num)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch the runtime details
	runtime, err := runtimes.GetRuntime(sub.Runtime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// If image is not pulled, return error
	if !checkImageIsPulled(runtime.Image) {
		return
	}

	// Run each test case in a new Go routine
	var wg sync.WaitGroup
	var tcResponses []TestCaseResponse
	for _, tc := range test.TestCases {
		wg.Add(1)

		go func(tc tests.TestCase) {
			// Get a UID
			id, err := sf.NextID()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Request code to byte array
			content := []byte(strings.Join(sub.Code, "\n"))

			// Code into temp file
			dir, err := filepath.Abs(filepath.Dir("."))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Convert UID to string
			strID := strconv.FormatUint(id, 10)

			os.Mkdir(dir + "/temp/" + strID, 0755)
			file, err := os.Create(dir + "/temp/" + strID + "/" + runtime.FileName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			file.Write(content)

			// Setup commands for container
			cmds := []string{}
			for _, str := range runtime.Commands {
				// Substitute in runtime specific variables
				str = strings.ReplaceAll(str, "{INPUTS}", strings.Join(tc.Inputs, " "))
				str = strings.ReplaceAll(str, "{LOCATION}", strID)

				cmds = append(cmds, str)
			}

			// Create the container, mounting a volume to share the temp file
			resp, err := cli.ContainerCreate(context.Background(), &container.Config {
				Image: runtime.Image,
				Cmd: cmds,
			}, &container.HostConfig {
				Mounts: []mount.Mount {
					{
						Type: mount.TypeVolume,
						Source: "cc-data",
						Target: "/tmp",	
					},
				},
			}, nil, nil, "")
			if err != nil {
				fmt.Println(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Start the container...
			err = cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// ...and wait for it to run, capturing all outputs...
			statusCh, errCh := cli.ContainerWait(context.Background(), resp.ID, container.WaitConditionNotRunning)
			select {
			case err := <-errCh:
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			case <-statusCh:
			}

			// ...then read all the outputs...
			logReader, err := cli.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer logReader.Close()

			// ...and transform it into a string
			logContent, _ := ioutil.ReadAll(logReader)

			// Clean up container
			err = cli.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Clean up temp files
			err = os.RemoveAll("./" + strconv.Itoa(int(id)))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Remove trailing empty log entries
			output := strings.TrimSuffix(string(logContent), "\n")

			// Remove non-Ascii chars from output
			splitOut := strings.Split(output, "\n")
			for i := range splitOut {
				splitOut[i] = stringsx.Clean(splitOut[i])
			}

			// Build response
			tcr := TestCaseResponse{}
			tcr.Number = tc.Number
			tcr.Pass = strings.Join(splitOut, "\n") == strings.Join(tc.ExpectedOutput, "\n")
			tcr.Inputs = tc.Inputs
			tcr.Outputs = splitOut
			tcr.Expected = tc.ExpectedOutput

			tcResponses = append(tcResponses, tcr)

			// Go routine is done
			wg.Done()
		}(tc)
	}

	// Wait for all Go routines to finish
	wg.Wait()

	// If all test cases have passed, overall pass
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(json)
}

// Handles GET request for the details of a specific runtime.
// Runtime is specified by the "name" query param.
func GetRuntimeDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Fetch the runtime details
	runtime, err := runtimes.GetRuntime(vars["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := RuntimeDetailResponse{}
	res.Name = runtime.Name
	res.Display = runtime.Display
	res.Image = runtime.Image
	res.Installed = checkImageIsPulled(runtime.Image)
	res.Template = strings.Join(runtime.Template, "\n")

	json, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(json)
}

// Handles a GET request for the details of all runtimes
func GetAllRuntimeDetails(w http.ResponseWriter, r *http.Request) {
	// Get all runtimes
	res := []RuntimeDetailResponse{}
	for _, runtime := range runtimes.GetAllRuntimes() {
		det := RuntimeDetailResponse{}
		det.Name = runtime.Name
		det.Display = runtime.Display
		det.Image = runtime.Image
		det.Installed = checkImageIsPulled(runtime.Image)
		det.Template = strings.Join(runtime.Template, "\n")

		res = append(res, det)
	}

	json, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(json)
}

// Handles a POST request to install a specific runtime. Runtime
// is specified by the "name" query param. To install the runtime,
// the relevant Docker image will be pulled
func InstallRuntime(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Fetch the runtime details
	runtime, err := runtimes.GetRuntime(vars["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Pull Docker image
	reader, err := cli.ImagePull(context.Background(), runtime.Image, types.ImagePullOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.Copy(os.Stdout, reader)

	// Build HTTP response (details of runtime)
	res := RuntimeDetailResponse{}
	res.Name = runtime.Name
	res.Display = runtime.Display
	res.Image = runtime.Image
	res.Installed = checkImageIsPulled(runtime.Image)
	res.Template = strings.Join(runtime.Template, "\n")

	json, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(json)
}

// Handles a GET request to fetch the details of a specific
// test. Test is specified by the "test_number" query param
func GetTestDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Pull number from query params...
	num, err := strconv.Atoi(vars["test_number"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ...and fetch the details of the corresponding test
	test, err := tests.GetTestByNumber(num)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Build HTTP response
	res := TestDetailResponse{}
	res.Title = test.Title
	res.Description = test.Description
	res.Number = test.Number

	json, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(json)
}

// Handles a GET request to fetch the details of all tests
func GetAllTestDetails(w http.ResponseWriter, r *http.Request) {
	// Get all tests
	res := []TestDetailResponse{}
	for _, test := range tests.GetAllTests() {
		det := TestDetailResponse{}
		det.Title = test.Title
		det.Description = test.Description
		det.Number = test.Number

		res = append(res, det)
	}

	json, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(json)
}

// Checks if Docker image has been pulled on the host
func checkImageIsPulled(name string) bool {
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return false
	}

	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == name {
				return true
			}
		}
	}

	return false
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
	
	// Handle routes
	router.HandleFunc("/tests/{test_number}/submission", TestSubmission).Methods("POST")
	router.HandleFunc("/tests/{test_number}", GetTestDetails).Methods("GET")
	router.HandleFunc("/tests", GetAllTestDetails).Methods("GET")
	router.HandleFunc("/runtimes/{name}", GetRuntimeDetails).Methods("GET")
	router.HandleFunc("/runtimes", GetAllRuntimeDetails).Methods("GET")
	router.HandleFunc("/runtimes/{name}/install", InstallRuntime).Methods("POST")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	// Setup Docker Client
	var err error
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	origins := handlers.AllowedOrigins([]string{"*"})
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Origin", "Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	http.ListenAndServe(":8080", handlers.CORS(origins, headers, methods)(router))
}