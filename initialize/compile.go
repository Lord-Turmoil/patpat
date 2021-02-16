package initialize

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func CompileJava(command string, args ...string) {
	subProcess := exec.Command(command, args...)
	stdin, err := subProcess.StdinPipe()
	if err != nil {
		panic(err) // replace with logger, or anything you want
	}
	defer stdin.Close() // the doc says subProcess. Wait will close it, but I'm not sure, so I kept this line
	subProcess.Stdout = os.Stdout
	subProcess.Stderr = os.Stderr
	fmt.Println("START COMPILE")
	if err = subProcess.Start(); err != nil { // Use start, not run
		panic(err) // replace with logger, or anything you want
	}
	subProcess.Wait()
	fmt.Println("END COMPILE")
}

func RunCommand(name string, args ...string) (exitCode int) {
	// log.Println("run command:", name, args)
	var outBuf, errBuf bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	stdout := outBuf.String()
	stderr := errBuf.String()

	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			// This will happen (in OSX) if `name` is not available in $PATH,
			// in this situation, exit code could not be get, and stderr will be
			// empty string very likely, so we use the default fail code, and format err
			// to string and set to stderr
			log.Printf("Could not get exit code for failed program: %v, %v", name, args)
			exitCode = -100 // defaultFailedCode
			if stderr == "" {
				stderr = err.Error()
			}
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}
	fmt.Printf("Compile result: (stdout: %v) (stderr: %v) (exitCode: %v)\n", stdout, strings.TrimRight(stderr, "\r\n"), exitCode)
	return exitCode
}
