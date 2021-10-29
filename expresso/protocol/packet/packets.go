package packet

// packetCollection represents a collection of client-bound and server-bound packets.
type packetCollection struct {
	// clientBoundPackets holds client-bound packets for the state.
	clientBoundPackets map[int32]func() Packet
	// serverBoundPackets holds server-bound packets for the state.
	serverBoundPackets map[int32]func() Packet
}

var (
	// handshakingCollection is the packet collection for the handshaking state.
	handshakingCollection = &packetCollection{
		clientBoundPackets: map[int32]func() Packet{},
		serverBoundPackets: map[int32]func() Packet{
			0x00: func() Packet { return &Handshake{} },
		},
	}
	// statusCollection is the packet collection for the status state.
	statusCollection = &packetCollection{
		clientBoundPackets: map[int32]func() Packet{
			0x00: func() Packet { return &ServerStatusResponse{} },
			0x01: func() Packet { return &ServerStatusPong{} },
		},
		serverBoundPackets: map[int32]func() Packet{
			0x00: func() Packet { return &ClientStatusRequest{} },
			0x01: func() Packet { return &ClientStatusPing{} },
		},
	}
	// loginCollection is the packet collection for the login state.
	loginCollection = &packetCollection{
		clientBoundPackets: map[int32]func() Packet{
			0x00: func() Packet { return &LoginDisconnect{} },
			0x01: func() Packet { return &EncryptionRequest{} },
			0x02: func() Packet { return &LoginSuccess{} },
			0x03: func() Packet { return &SetCompression{} },
			// Modern code shouldn't be using login plugin request packets, which is why it isn't here.
		},
		serverBoundPackets: map[int32]func() Packet{
			0x00: func() Packet { return &LoginStart{} },
			0x01: func() Packet { return &EncryptionResponse{} },
			// Modern code shouldn't be using login plugin response packets, which is why it isn't here.
		},
	}
	// playCollection is the packet collection for the play state.
	playCollection = &packetCollection{
		// TODO: Add all play packets.
		clientBoundPackets: map[int32]func() Packet{
			0x1A: func() Packet { return &Disconnect{} },
			0x21: func() Packet { return &ServerKeepAlive{} },
			0x26: func() Packet { return &JoinGame{} },
			0x38: func() Packet { return &ServerPlayerPositionRotation{} },
			0x49: func() Packet { return &UpdateViewPosition{} },
		},
		serverBoundPackets: map[int32]func() Packet{
			0x0F: func() Packet { return &ClientKeepAlive{} },
		},
	}
)
