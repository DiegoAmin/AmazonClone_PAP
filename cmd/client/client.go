package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:10000")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connected to server!")

	// Read from server and write to stdout in a goroutine
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "EXIT" {
				os.Exit(0)
			}
			fmt.Println(line)
		}
	}()

	// Read from stdin and send to server
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		fmt.Fprintln(conn, input.Text())
	}
}
