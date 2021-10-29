package main

import (
	"fmt"
	"github.com/justtaldevelops/expresso/expresso"
	"github.com/justtaldevelops/expresso/expresso/protocol"
	"github.com/justtaldevelops/expresso/expresso/protocol/packet"
	"math"
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
		Worlds:       []string{"minecraft:world"},
		World:        "minecraft:world",
		HashedSeed:   100,
		ViewDistance: 16,
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

	// Initialize empty chunks.
	emptyChunks := make([]*protocol.Chunk, 256)
	for i := 0; i < 256; i++ {
		emptyChunks[i] = protocol.NewEmptyChunk()
	}

	// Send empty columns in a radius of sixteen.
	radius := int32(16)
	for x := -radius; x <= radius; x++ {
		for z := -radius; z <= radius; z++ {
			// Make sure we're in bounds.
			distance := math.Sqrt(float64(x*x) + float64(z*z))
			chunkDistance := int32(math.Round(distance))
			if chunkDistance > radius {
				// The column was outside the chunk radius.
				continue
			}

			// Write the column data for this specific chunk.
			err = conn.WritePacket(&packet.ChunkData{Column: protocol.Column{
				X: x, Z: z,
				Chunks:     emptyChunks,
				Tiles:      []map[string]interface{}{},
				HeightMaps: map[string]interface{}{},
				Biomes:     emptyBiomeData,
			}})
		}
	}

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
				panic(err)
			}
			fmt.Printf("%+v\n", pk)
		}
	}()
}
