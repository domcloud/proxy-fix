package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
)

var outPort int

func init() {
	var err error
	outPort, err = getFreePort()
	if err != nil {
		panic("Can't get free port")
	}
	// Check if there are additional arguments
	if len(os.Args) > 1 {
		cmd := exec.Command(os.Args[1], os.Args[2:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = append(cmd.Env, fmt.Sprintf("PORT=%d", outPort))

		// Start the specified command
		err := cmd.Start()
		if err != nil {
			fmt.Printf("Error starting command: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Started process %s with PID %d\n", os.Args[1], cmd.Process.Pid)

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
