package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"sync"
)

func server(group *sync.WaitGroup) {
	// listen on a port
	defer group.Done()
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		fmt.Println(err)
		return
	}
	for {
		// accept a connection
		c, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", ln.Addr(), err)
			fmt.Println(err)
		}
		// handle the connection
		go handleServerConnection(c)
	}
}

func handleServerConnection(c net.Conn) {
	//receive the message
	var msg string
	err := gob.NewDecoder(c).Decode(&msg)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("received message:", msg)

	}

	defer func() {
		i := 0
		err = c.Close()
	TryAgain:
		if err != nil {
			err = c.Close()
			i += 1
			if i < 5 {
				goto TryAgain
			} else {
				fmt.Println(err)
				fmt.Println("Couldn't close Server connection")

			}
		}
	}()
}

func client(group *sync.WaitGroup) {
	// connect to the server
	defer group.Done()
	c, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		i := 0
		err = c.Close()
	TryAgain:
		if err != nil {
			err = c.Close()
			i += 1
			if i < 5 {
				goto TryAgain
			} else {
				fmt.Println(err)
				fmt.Println("Couldn't close connection")
			}
		}
	}()
	// send the message
	fmt.Print("What shall we send? ")
	var msg string
	reader := bufio.NewScanner(os.Stdin)
	_ = reader.Scan()   // Scan single line.
	msg = reader.Text() // get text from line.
	if err != nil {
		fmt.Println("Couldn't read from console.")
		fmt.Println(err)
		return
	}
	msg += "\n"
	fmt.Println("We are sending:", msg)
	err = gob.NewEncoder(c).Encode(msg)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go server(&wg)
	wg.Add(1)
	go client(&wg)
	var input string
	wg.Wait()
	fmt.Print("End? ")
	_, err := fmt.Scanln(&input)
	if err != nil {
		return
	}
}
