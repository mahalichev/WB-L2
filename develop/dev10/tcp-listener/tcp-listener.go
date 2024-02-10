package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:3010")
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}
	defer listener.Close()

	go func() {
		for {
			if conn, err := listener.Accept(); err == nil {
				go Handle(conn)
			} else {
				fmt.Fprint(os.Stderr, err)
				continue
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	signal.Stop(c)
	close(c)
}

func Handle(connection net.Conn) {
	defer connection.Close()
	reader := bufio.NewReader(connection)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Fprint(os.Stderr, err)
			}
			break
		}
		fmt.Print("client: ", string(message))
		connection.Write([]byte("server: " + message))
	}
}
