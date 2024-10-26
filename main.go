package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
)

var outPort int

func init() {
	var err error

	outPort, err = checkExistingProcess()
	if err == nil && isPortListening(outPort) {
		fmt.Printf("Process is already running on port %d\n", outPort)
		return
	}

	// No existing process found or not listening, start a new one
	outPort, err := getFreePort()
	if err != nil {
		panic("Can't get free port")
	}

	args := os.Args
	bg := os.Getenv("NOHUP") == "1"

	// Check if there are additional arguments
	if len(args) > 1 {
		if bg {
			var out bytes.Buffer
			var stderr bytes.Buffer
			cmdd := generateBgCmd(args[1], args[2:]...)
			cmd := exec.Command("bash", "-c", cmdd)
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			cmd.Env = os.Environ()
			cmd.Env = append(cmd.Env, fmt.Sprintf("PORT=%d", outPort))
			err = cmd.Run()
			if err != nil {
				fmt.Printf("Error starting command: %v\n", err)
				os.Exit(1)
			}
			pid, err := strconv.Atoi(out.String())
			if nil != err {
				fmt.Printf("Error starting command: %v\n", err)
				os.Exit(1)
			}
			if pid == 0 {
				fmt.Printf("invalid PID of 0")
				os.Exit(1)
			}
			writePidPortFile(pid, outPort)
			fmt.Printf("Started process %s with PID %d in background\n", args[1], pid)
		} else {
			cmd := exec.Command(args[1], args[2:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Env = os.Environ()
			cmd.Env = append(cmd.Env, fmt.Sprintf("PORT=%d", outPort))

			// Start the specified command
			err := cmd.Start()
			if err != nil {
				fmt.Printf("Error starting command: %v\n", err)
				os.Exit(1)
			}
			pid := cmd.Process.Pid
			fmt.Printf("Started process %s with PID %d\n", args[1], pid)
		}

	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to port 8080 if no PORT env variable is set
	}
	err := startProxy("0.0.0.0:" + port)
	if err != nil {
		panic(err)
	}
}

func startProxy(address string) error {
	proxy := Proxy{
		DialTarget: fmt.Sprintf("localhost:%d", outPort),
	}
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()
	fmt.Printf("Proxy server listening on %s\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go proxy.handleConnection(conn)
	}
}
