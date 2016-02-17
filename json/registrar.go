package main

import (
	"encoding/json"
	// "errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	// cfg "github.com/elastic/beats/filebeat/config"
	// "github.com/elastic/beats/filebeat/input"
	// . "github.com/elastic/beats/filebeat/input"
	// "github.com/elastic/beats/libbeat/logp"
)

// Registrar should only have one entry, which
// is the offset into the unified2 file
// currently being tailed (if any)
type Registrar struct {
	registryFile string    // path to the registry file
	State        FileState // unified2 file name and offset
	sync.Mutex             // lock and unlock during writes
}

type FileState struct {
	Source string `json:"source,omitempty"`
	Offset int64  `json:"offset,omitempty"`
}

func main() {
	// registry (bookmark) processing:
	// set path to registry file
	//   which should be in the same folder as unified2 files
	// create or load registry file - where were we?
	// if r.State.Source (unified2 file) exists:
	//   use r.State.Offset
	//   continue tailing/indexing r.State.Source
	// else
	//   WriteRegistry for the current file (if any)

	// test access denied to create file in folder:
	// r, err := NewRegistrar("/var/log/apache/.u2")
	// test access allowed to create file in folder:
	r, err := NewRegistrar("/var/log/u2beat/.u2")
	if err != nil {
		fmt.Printf("NewRegistrar: err=%v\n", err)
		// any error here is critical for unifiedbeat,
		// so we must exit immediately:
		os.Exit(1)
	}
	fmt.Printf("NewRegistrar result: r=%T=%#v\n", r, r)

	// test load/read state from registry file, which may or may not exist yet:
	r.LoadState()
	fmt.Printf("\n>>> LoadState:\n\tr=%T=%#v\n\n", r, r)

	// test writing state to registry file:
	r.State.Source = "snort.log.1"
	r.State.Offset = 987
	err = r.WriteRegistry()
	if err != nil {
		fmt.Printf("WriteRegistry result: err=%v\n", err)
	}
	fmt.Printf("WriteRegistry: r=%T=%#v\n", r, r)

	// test reading state from registry file:
	r.State.Source = ""
	r.State.Offset = 0
	r.LoadState()
	fmt.Printf("LoadState:\n\tr=%T=%#v\n", r, r)
}

func NewRegistrar(registryFile string) (*Registrar, error) {
	fmt.Printf("NewRegistrar: registryFile=%#v\n", registryFile)
	r := &Registrar{
		registryFile: registryFile,
	}

	// Ensure we have access to write registryFile, of course,
	// this could still fail in later calls to LoadState or WriteRegistry.
	// There is no perfect solution as files and permissions are just a mess,
	// but the issue should be resolved during deployment.
	// The big assumption is that unifiedbeat will not be run on Windows.
	testfile := r.registryFile + ".access.test"
	file, err := os.Create(testfile)
	// err = errors.New("Create failed!")
	if err != nil {
		fmt.Printf("NewRegistrar: test 'create file' access was denied to path for registry file: '%v'\n", r.registryFile)
		return nil, err
	}
	err = file.Close()
	// err = errors.New("Close failed!")
	if err != nil {
		// really? we lost access after Create, really?
		fmt.Printf("NewRegistrar: test 'close file' access was denied to path for registry file: '%v'\n", r.registryFile)
		return nil, err
	}
	err = os.Remove(testfile)
	// err = errors.New("Remove failed!")
	if err != nil {
		// really? we lost access after Create and Close, really?
		fmt.Printf("NewRegistrar: test 'remove file' access was denied to path for registry file: '%v'\n", r.registryFile)
		return nil, err
	}

	// Set absolute path to the registryFile
	absPath, err := filepath.Abs(r.registryFile)
	// err = errors.New("filepath.Abs failed!")
	if err != nil {
		fmt.Printf("NewRegistrar: failed to set the absolute path for registry file: '%s'\n", r.registryFile)
		return nil, err
	}
	r.registryFile = absPath
	fmt.Printf("NewRegistrar: registry file set to: %s\n", r.registryFile)

	return r, err
}

// LoadState fetches the previous reading state from the RegistryFile
// The default file is ".unifiedbeat" file which is stored in the same path as the unified2 files
func (r *Registrar) LoadState() {
	if existing, e := os.Open(r.registryFile); e == nil {
		defer existing.Close()
		fmt.Printf("Loading registrar data from %s\n", r.registryFile)
		decoder := json.NewDecoder(existing)
		decoder.Decode(&r.State)
	}
}

// WriteRegistry writes the new json registry info to a file.
func (r *Registrar) WriteRegistry() error {
	// logp.Debug("registrar", "Write registry file: %s", r.registryFile)
	fmt.Printf("WriteRegistry file: %s\n", r.registryFile)

	r.Lock()
	defer r.Unlock()
	// can't truncate a file that does not exist:
	_, err := os.Stat(r.registryFile)
	// fmt.Printf("WriteRegistry: os.Stat: err=%v\n", err)
	// return err
	if os.IsExist(err) {
		err := os.Truncate(r.registryFile, 0) // the file must not be open on Windows!
		if err != nil {
			fmt.Printf("WriteRegistry: os.Truncate: err=%v\n", err)
			return err
		}
	}
	// if marshal json, or writefile fail then most likely
	// unifiedbeat does not have access to the registryFile's path
	jsonState, err := json.Marshal(r.State)
	if err != nil {
		fmt.Printf("WriteRegistry: json.Marshal: err=%v\n", err)
		return err
	}
	// https://golang.org/pkg/io/ioutil/#WriteFile
	//   If the file does not exist, WriteFile creates it with permissions perm;
	//   otherwise WriteFile truncates it before writing.
	err = ioutil.WriteFile(r.registryFile, jsonState, 0644)
	if err != nil {
		fmt.Printf("WriteRegistry: ioutil.WriteFile: err=%v\n", err)
		return err
	}

	fmt.Printf("Registry file updated: file: %v offset: %v.\n", r.State.Source, r.State.Offset)

	return nil
}
