package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const MAX_RETRY = 20
const WAIT_RETRY = time.Second * 1

var invalidHeaderRegex = regexp.MustCompile("[^a-zA-Z0-9-]+")

func getFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}

func filterInvalidHeaders(headers http.Header) {
	for header := range headers {
		if invalidHeaderRegex.MatchString(header) {
			delete(headers, header)
		}
	}
}

func getPidFile() string {
	return filepath.Join(os.Getenv("HOME"), "tmp", "app.pid")
}

func checkExistingProcess() (port int, pid int, err error) {
	data, err := os.ReadFile(getPidFile())
	if err != nil {
		return 0, 0, err
	}

	fields := strings.Split(string(data), ":")
	if len(fields) != 2 {
		return 0, 0, fmt.Errorf("invalid pid/port file format")
	}

	pid, err = strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid PID format: %v", err)
	}

	port, err = strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid port format: %v", err)
	}

	if processExists(pid) {
		return port, pid, nil
	}
	return 0, 0, fmt.Errorf("no running process found")
}

func processExists(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

func processKill(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = process.Signal(syscall.SIGKILL)
	return err == nil
}

func isPortListening(port int) bool {
	retries := 0
retry:
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		if retries < MAX_RETRY {
			retries += 1
			fmt.Printf("Retrying to connect for %d-th time\n", retries)
			time.Sleep(WAIT_RETRY)
			goto retry
		}
		fmt.Printf("Error connecting to destination: %v\n", err)
		return false
	}
	conn.Close()
	return true
}

func writePidPortFile(pid, port int) {
	err := os.WriteFile(getPidFile(), []byte(fmt.Sprintf("%d:%d", pid, port)), 0644)
	if err != nil {
		fmt.Printf("Failed to write PID/port file: %v\n", err)
	}
}

func generateBgCmd(name string, arg ...string) string {
	return fmt.Sprintf(
		"nohup %s %s < /dev/null &>$HOME/tmp/app.log & echo -n $! | awk '/[0-9]+$/{ printf $0 }'",
		name,
		strings.Join(arg, " "),
	)
}
