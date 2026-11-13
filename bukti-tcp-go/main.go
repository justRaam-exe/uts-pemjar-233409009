package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: go run ./bukti-tcp-go <host:port>\n")
		os.Exit(1)
	}

	target := os.Args[1]
	if !strings.Contains(target, ":") {
		target = target + ":8080"
	}

	conn, err := net.Dial("tcp", target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to %s: %v\n", target, err)
		os.Exit(1)
	}
	defer conn.Close()

	httpRequest := "GET /healthz HTTP/1.1\r\n" + "Host: " + target + "\r\n" + "Connection: close\r\n" + "\r\n"

	_, err = conn.Write([]byte(httpRequest))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to send TCP data: %v\n", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read response: %v\n", err)
		os.Exit(1)
	}
}
