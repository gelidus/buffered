## Buffered - a byte packer (that maybe works)

Buffered is a small library that implements length prefix protocol for
byte arrays sent via network. It Prefixes given content by it's length,
so the byte array can be constructed back if fragmented by the network.

**Features:**
- [x] 100% test coverage! (of that 20 lines of code)
- [x] documentation
- [x] it works
- [ ] maybe

## Installation

```bash
go get github.com/gelidus/buffered
```

## Simple Usage

```go
bytez := []byte{'a', 'b', 'c', 'd'}
// Pack, PackLE(little endian) or PackBE(big endian) can be used
packed := buffered.PackLE(bytez)

// we need io.Reader for the unpacks to work, so we create new buffer
unpacked, err := buffered.UnpackLE(bytes.NewBuffer(packed))
```

## Usage with TCP Socket

```go
conn, _ := net.Dial("tcp", "whatever.address:port")

// read from the socket
// PackLE should be used by the sender
// this will wait until new packet arrives
packet := buffered.UnpackLE(conn)
// use your packet here (without prefixed length)

// write to the socket
conn.Write(buffered.PackLE(myBytes))
```

## Defining length Encoder and Decoder (lol wat?)

If little and big endianity is too mainstream to you, Buffered has your
back! You can simply define Length decoders and encoders for `buffered.Pack`
and `buffered.Unpack` functions.

```go
func MyUInt32Encoder(dest []byte, length uint32) {
  var leBytes []byte = make([]byte, 4)
  binary.LittleEndian.PutUint32(leBytes, length)
  
  // converts 0 1 2 3
  // into     1 3 2 0
  dest[0] = leBytes[1]
  dest[1] = leBytes[3]
  dest[2] = leBytes[2]
  dest[3] = leBytes[0]
}

func MyUInt32Decoder(src []byte) uint32 {
  var decoded []byte = make([]byte, 4)
  
  // converts 1 3 2 0
  // back to  0 1 2 3
  decoded[0] = src[1]
  decoded[1] = src[3]
  decoded[2] = src[2]
  decoded[3] = src[0]
  
  return binary.LittleEndian.Uint32(src)
}

// They can now be used by pack and unpack to decode length of data
buffered.Pack(myBytes, MyUInt32Encoder)
// and
buffered.Unpack(myReader, MyUInt32Decoder)
```