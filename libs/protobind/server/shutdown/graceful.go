package main

import (
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"syscall"
)

const pidFile = "/tmp/atomicswap.grpc.server.pid"

func main() {
	if runtime.GOOS == "windows" {
		log.Fatalln("This does not work on Windows")
	}
	dat, err := ioutil.ReadFile(pidFile)
	if err != nil {
		log.Fatalf("Cannot read %s, error: %v\n", pidFile, err)
	}
	sPid := string(dat)
	log.Printf("pid: %s\n", sPid)
	pid, err := strconv.ParseInt(sPid, 10, 64)
	if err != nil {
		log.Fatalf("Cannot parse %s, error: %v\n", string(dat), err)
	}
	process, err := os.FindProcess(int(pid))
	if err != nil {
		log.Fatalf("Cannot find process %d, error: %v\n", pid, err)
	}
	// log.Printf("Process %v\n", process)
	process.Signal(syscall.SIGTERM) // Not on windows
	log.Printf("SIGTERM: =>  %d\n", pid)
}
