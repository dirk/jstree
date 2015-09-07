package jstree

import "bytes"
import "fmt"
import "io/ioutil"
import "log"
import "os/exec"
import "path"
import "path/filepath"
import "syscall"
import "runtime"

import "github.com/bitly/go-simplejson"

func GetAcornPath() string {
	// Current file of this call
	_, file, _, _ := runtime.Caller(0)

	dir := filepath.Dir(file)
	// Ascend to the root of this package
	dir = filepath.Join(dir, "..", "..")

	return path.Join(dir, "deps", "acorn-2.4.0", "bin", "acorn")
}

func ParseFile(file string) (*Program, error) {
	var status syscall.WaitStatus
	var exitCode int
	var err error

	acornPath := GetAcornPath()

	cmd := exec.Command(acornPath, "--ecma6", "--module", file)

	stdoutBuffer := &bytes.Buffer{}
	stderrBuffer := &bytes.Buffer{}

	cmd.Stdout = stdoutBuffer
	cmd.Stderr = stderrBuffer

	err = cmd.Run()
	if err != nil {
		exitError, _ := err.(*exec.ExitError)
		status   = exitError.Sys().(syscall.WaitStatus)
		exitCode = status.ExitStatus()
	} else {
		exitCode = -1
	}

	stdoutBytes, err := ioutil.ReadAll(stdoutBuffer)
	if err != nil { log.Fatal(err) }

	stderrBytes, err := ioutil.ReadAll(stderrBuffer)
	if err != nil { log.Fatal(err) }

	if exitCode > -1 {
		if len(stdoutBytes) > 0 { fmt.Println(string(stdoutBytes)) }
		if len(stderrBytes) > 0 { fmt.Println(string(stderrBytes)) }

		return nil, fmt.Errorf("Exited with code: %d; %s", exitCode, string(stderrBytes))
	}

	json, err := simplejson.NewJson(stdoutBytes)
	if err != nil { return nil, err }

	return ParseProgram(json)
}
