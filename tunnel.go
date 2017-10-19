package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"golang.org/x/net/proxy"
	"golang.org/x/net/websocket"
)

var (
	target = flag.String("target", "", "The target host:port to tunnel to")
	port       = flag.Int("port", 8080, "The local port to listen on")
)

func iocopy(dst io.Writer, src io.Reader, c chan error) {
	_, err := io.Copy(dst, src)
	c <- err
}

func handleConnection(wsConfig *websocket.Config, conn net.Conn) {
	defer conn.Close()

	tcp, err := proxy.FromEnvironment().Dial("tcp", wsConfig.Location.Host)
	if err != nil {
		log.Print("proxy.FromEnvironment().Dial(): ", err)
		return
	}

	ws, err := websocket.NewClient(wsConfig, tcp)
	if err != nil {
		log.Print("websocket.NewClient(): ", err)
		return
	}
	defer ws.Close()

	c := make(chan error, 2)
	go iocopy(ws, conn, c)
	go iocopy(conn, ws, c)

	for i := 0; i < 2; i++ {
		if err := <-c; err != nil {
			fmt.Print("io.Copy(): ", err)
			return
		}
	}
}

func main() {
	flag.Parse()

	config, err := websocket.NewConfig(*target, "http://localhost/")
	if err != nil {
		panic(err)
	}

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print("ln.Accept(): ", err)
			continue
		}
		go handleConnection(config, conn)
	}
}