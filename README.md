# expresso
> A library designed for hosting Minecraft: Java Edition listeners (currently for 1.18.0).

## Features
- [X] Hosting listeners.
- [X] All handshake, status, and login state packets.
- [X] Login verification.
- [X] Compression and encryption.
- [X] Chunk column reading/writing.
- [ ] Dialing/connecting to listeners.
- [ ] All play state packets.

## Example
You can find a basic example in main.go. The example sends a chunk column for `(0, 0)` to every connection which has
the block at `(0, 1, 0)` and `(1, 3, 0)` set to the block state of `10`, which is `minecraft:dirt`.

## Disclaimer
Do not expect anything completely working right now! Currently, there's only enough to get the player spawned in the
world, and for chunk data to be sent to the client. There is also no support for connecting to listeners at the moment,
however it is planned.

## Credits
These projects helped me design expresso and gave the general idea of how to build a protocol library for Minecraft.
Many thanks to all the authors and contributors of these projects!

### [wiki.vg](https://wiki.vg/Protocol)
An absolute godsend for any project interesting the Java Edition protocol. Contains a lot of useful information for
getting on the right track, and documents the entire protocol, while still being mostly up to date.

### [go-mc](https://github.com/Tnze/go-mc)
Many parts of expresso are based off of go-mc, such as the BitStorage implementation or certain parts of the
reader/writer. I would like to thank the authors of go-mc for their work, and for making it possible to write
this library.

### [gophertunnel](https://github.com/Sandertv/gophertunnel)
gophertunnel helped me with the general design of packets and reader/writers, as well as the implementation for NBT.
If you're interested in the Bedrock protocol, I would definitely recommend using gophertunnel.

### [MCProtocolLib](https://github.com/GeyserMC/MCProtocolLib)
Much of the chunk implementation was inspired from this project. It's a pretty big library and much more established 
and complete compared to this implementation. I would recommend using it if you're interested in utilizing the protocol 
in Java and are looking for something more complete.
