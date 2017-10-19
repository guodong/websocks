package main

import (
	"flag"
	"github.com/armon/go-socks5"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

func main() {
	flag.Parse()
	socks, err := socks5.New(&socks5.Config{})
	if err != nil {
		panic(err)
	}
	http.Handle("/", websocket.Handler(func(conn *websocket.Conn) { socks.ServeConn(conn) }))
	log.Fatal(http.ListenAndServe(*addr, nil))
}
