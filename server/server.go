package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

type Message struct {
	sender  int
	message string
}

func handleError() {
	// TODO: all
	// Deal with an error event.
	fmt.Println("Something went wrong!")
	os.Exit(1)
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	// TODO: all
	// Continuously accept a network connection from the Listener
	// and add it to the channel for handling connections.
	for {
		conn, err := ln.Accept() //accept waits for clients to Dial/block
		if err != nil {
			handleError()
		}
		conns <- conn
	}
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	// TODO: all
	// So long as this connection is alive:
	// Read in new messages as delimited by '\n's
	// Tidy up each message and add it to the messages channel,
	// recording which client it came from.
	reader := bufio.NewReader(client)

	for {
		msg, err := reader.ReadString('\n') //blocks until user inputs a message
		if err != nil {
			fmt.Printf("User %d left chat.\n", clientid)
			return
		}
		msgs <- Message{clientid, msg}
	}
	
}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	id := 0

	//TODO Create a Listener for TCP connections on the port given above.
	netLis, err := net.Listen("tcp", *portPtr)
	if err != nil {
		handleError()
	}

	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn) //a map with key int and value Conn

	//Start accepting connections
	go acceptConns(netLis, conns)
	for {
		select {
		case conn := <-conns:
			//TODO Deal with a new connection
			// - assign a client ID
			// - add the client to the clients map
			// - start to asynchronously handle messages from this client
			fmt.Printf("User %d entered chat.\n", id)
			clients[id] = conn
			go handleClient(conn, id, msgs)
			id++
		case msg := <-msgs:
			//TODO Deal with a new message
			// Send the message to all clients that aren't the sender
			for k, v := range clients {
				if k != msg.sender {
					fmt.Fprint(v, msg.message)
				}
			}
		}
	}
}
