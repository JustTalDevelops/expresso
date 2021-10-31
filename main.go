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
	// Join the game client side.
	err := conn.WritePacket(&packet.JoinGame{
		GameMode:         1,
		PreviousGameMode: 1,
		Worlds:           []string{"minecraft:world"},
		World:            "minecraft:world",
		HashedSeed:       100,
		ViewDistance:     16,
	})
	if err != nil {
		panic(err)
	}

	// Set the player's position and rotation.
	err = conn.WritePacket(&packet.ServerPlayerPositionRotation{})
	if err != nil {
		panic(err)
	}

	// Initialize empty biome data.
	emptyBiomeData := make([]int32, 1024)
	for i := 0; i < 1024; i++ {
		emptyBiomeData[i] = 1
	}

	chunk := protocol.NewEmptyChunk()
	err = chunk.Set(0, 0, 0, 10)
	if err != nil {
		panic(err)
	}

	// Write the column data for this specific chunk.
	err = conn.WritePacket(&packet.ChunkData{Column: protocol.Column{
		X: 0, Z: 0,
		Chunks: []*protocol.Chunk{
			nil, nil, nil, nil, nil, nil, nil, chunk, nil, chunk, nil, nil, nil, chunk,
			nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
			nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
			nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
			nil, nil, nil, nil, nil,
		},
		Tiles:      []map[string]interface{}{},
		HeightMaps: map[string]interface{}{},
		Biomes:     emptyBiomeData,
	}})

	// Write the column data for this specific chunk x2.
	err = conn.WritePacket(&packet.ChunkData{Column: protocol.Column{
		X: 1, Z: 1,
		Chunks: []*protocol.Chunk{
			chunk, chunk, chunk, chunk, chunk, chunk, chunk, chunk, chunk, chunk,
			chunk, chunk, chunk, chunk, chunk, chunk, nil, nil, nil, nil, nil, nil,
			nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
			nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
			nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		},
		Tiles:      []map[string]interface{}{},
		HeightMaps: map[string]interface{}{},
		Biomes:     emptyBiomeData,
	}})

	// Update the player's viewing column.
	err = conn.WritePacket(&packet.UpdateViewPosition{})
	if err != nil {
		panic(err)
	}

	// The client should now be spawned in.
	fmt.Println("We should be spawned in!")
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
