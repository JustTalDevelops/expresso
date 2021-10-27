package main

import (
	"fmt"
	"github.com/justtaldevelops/expresso/expresso"
)

func main() {
	l, err := expresso.Listen(":25565")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn *expresso.Connection) {
	fmt.Println("Made it to the play phase!")
	for {
		fmt.Println("Reading packet...")
		pk, err := conn.ReadPacket()
		if err != nil {
			break
		}
		fmt.Printf("%+v\n", pk)
	}
}
