package main

import (
	"fmt"
	"github.com/justtaldevelops/expresso/expresso"
	"github.com/justtaldevelops/expresso/expresso/protocol"
	"github.com/justtaldevelops/expresso/expresso/protocol/packet"
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
	emptyBiome := make([]int32, 1024)
	for i := 0; i < 1024; i++ {
		emptyBiome[i] = 1
	}

	fmt.Println("join game")
	err := conn.WritePacket(&packet.JoinGame{
		Worlds:       []string{"minecraft:world"},
		World:        "minecraft:world",
		HashedSeed:   100,
		ViewDistance: 16,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("position rotation")
	err = conn.WritePacket(&packet.ServerPlayerPositionRotation{})
	if err != nil {
		panic(err)
	}
	fmt.Println("chunk data")
	err = conn.WritePacket(&packet.ChunkData{Column: protocol.Column{
		X: 0,
		Z: 0,
		Chunks: []*protocol.Chunk{
			protocol.NewEmptyChunk(),
		},
		TileEntities: []map[string]interface{}{},
		HeightMaps:   map[string]interface{}{},
		Biomes:       emptyBiome,
	}})
	fmt.Println("update view position")
	err = conn.WritePacket(&packet.UpdateViewPosition{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Made it to the play phase!")
	go func() {
		for {
			fmt.Println("Reading packet...")
			pk, err := conn.ReadPacket()
			if err != nil {
				break
			}
			fmt.Printf("%+v\n", pk)
		}
	}()
}
