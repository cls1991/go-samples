package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Panic(err)
	}
	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleClientRequest(client)
	}
}

func handleClientRequest(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()
	var buf [1024]byte
	n, err := client.Read(buf[:])
	if err != nil {
		log.Println(err)
		return
	}
	var method, host, address string
	fmt.Sscanf(string(buf[:bytes.IndexByte(buf[:], '\n')]), "%s%s", &method, &host)
	hostPortUrl, err := url.Parse(host)
	if err != nil {
		log.Println(err)
		return
	}
	if hostPortUrl.Opaque == "443" {
		address = hostPortUrl.Scheme + ":443"
	} else {
		if strings.Index(hostPortUrl.Host, ":") == -1 {
			address = hostPortUrl.Host + ":80"
		} else {
			address = hostPortUrl.Host
		}
	}
	server, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return
	}
	if method == "CONNECT" {
		fmt.Fprint(client, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		server.Write(buf[:n])
	}
	go io.Copy(server, client)
	io.Copy(client, server)
}
