package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go-jellyfin-forwarder HOST")
		os.Exit(2)
	}

	log.SetFlags(log.Lshortfile)
	err := Run(os.Args[1])
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func Run(server string) error {
	const want = "who is JellyfinServer?"
	const port = "7359"
	conn, err := net.ListenPacket("udp", ":"+port)
	if err != nil {
		return err
	}
	defer conn.Close()

	server += ":" + port
	log.Println("listening")
	for {
		var buf [1024]byte
		n, client, err := conn.ReadFrom(buf[:])
		data := buf[:n]
		got := string(data)
		log.Println("got", n, client, err, got)
		if n == len(want) && got == want {
			log.Println("hello:", client, server)
			go Forward(conn, client, server, data)
			continue
		}
		if err != nil {
			log.Println("warning:", err)
			continue
		}
	}
}

func Forward(cconn net.PacketConn, client net.Addr, server string, data []byte) {
	err := forward(cconn, client, server, data)
	if err != nil {
		fmt.Println("error:", err)
	}
}

func forward(cconn net.PacketConn, client net.Addr, server string, data []byte) error {
	sconn, err := net.DialTimeout("udp", server, 30*time.Second)
	if err != nil {
		return err
	}
	defer sconn.Close()

	_, err = sconn.Write(data)
	if err != nil {
		return err
	}

	var buf [1024]byte
	n, err := sconn.Read(buf[:])
	data = buf[:n]
	got := string(data)
	log.Println("got2", n, err, got)
	if n > 0 {
		_, err2 := cconn.WriteTo(data, client)
		if err2 != nil {
			if err != nil {
				fmt.Println("error1:", err)
			}
			return err2
		}
	}
	return err
}
