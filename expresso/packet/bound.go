package packet

// Bound represents which receiver the packet is bound to.
type Bound struct {
	bound
}

// BoundClient is used when the packet is meant for clients.
func BoundClient() Bound {
	return Bound{bound: boundClient}
}

// BoundServer is used when the packet is meant for servers.
func BoundServer() Bound {
	return Bound{bound: boundServer}
}

type bound byte

const (
	boundClient bound = iota
	boundServer
)
