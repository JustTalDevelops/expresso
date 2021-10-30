package packet

// Direction represents what direction a packet is going.
type Direction struct {
	direction
}

// DirectionServer is used when the packet is meant for servers.
func DirectionServer() Direction {
	return Direction{direction: directionServer}
}

// DirectionClient is used when the packet is meant for clients.
func DirectionClient() Direction {
	return Direction{direction: directionClient}
}

type direction byte

const (
	directionServer direction = iota
	directionClient
)
