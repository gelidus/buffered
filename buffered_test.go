package buffered

import (
  "testing"
  "encoding/binary"
  "github.com/stretchr/testify/assert"
  "bytes"
  "net"
  "io"
)

var TestArrays [][]byte = [][]byte{
  {'a', 'b', 'c', 'd'},
  {15, 16, 17, 18},
  {'a', 'b', 'c', 'd', 'e', 'f'},
}

func TestPack(t *testing.T) {
  for _, arr := range TestArrays {
    res := Pack(arr, binary.LittleEndian.PutUint32)
    assert.Equal(t, len(arr) + 4, len(res))
    assert.Equal(t, byte(len(arr)), res[0])
  }
}

func TestPackBE(t *testing.T) {
  for _, arr := range TestArrays {
    res := PackBE(arr)
    assert.Equal(t, len(arr) + 4, len(res))
    assert.Equal(t, byte(len(arr)), res[3])
  }
}

func TestPackLE(t *testing.T) {
  for _, arr := range TestArrays {
    res := PackLE(arr)
    assert.Equal(t, len(arr) + 4, len(res))
    assert.Equal(t, byte(len(arr)), res[0])
  }
}

func TestUnpack(t *testing.T) {
  for _, arr := range TestArrays {
    packed := Pack(arr, binary.LittleEndian.PutUint32)
    unpacked, err := Unpack(bytes.NewBuffer(packed), binary.LittleEndian.Uint32)
    assert.Equal(t, arr, unpacked)
    assert.Equal(t, nil, err)
  }
}

func TestUnpackBE(t *testing.T) {
  for _, arr := range TestArrays {
    packed := PackLE(arr)
    unpacked, err := UnpackLE(bytes.NewBuffer(packed))
    assert.Equal(t, arr, unpacked)
    assert.Equal(t, nil, err)
  }
}

func TestUnpackLE(t *testing.T) {
  for _, arr := range TestArrays {
    packed := PackBE(arr)
    unpacked, err := UnpackBE(bytes.NewBuffer(packed))
    assert.Equal(t, arr, unpacked)
    assert.Equal(t, nil, err)
  }
}

func TestPackUnpackOnTCPSocket(t *testing.T) {
  // create the tcp listener
  listener, err := net.ListenTCP("tcp", nil)
  assert.Equal(t, nil, err)

  // start the server in new goroutine
  go func() {
    conn, err := listener.Accept()
    assert.Equal(t, nil, err)
    for _, arr := range TestArrays {
      packet := PackLE(arr)
      conn.Write(packet)
    }

    // close server socket and listener
    conn.Close()
    listener.Close()
  }()

  // connect the client
  conn, err := net.Dial("tcp", listener.Addr().String())
  assert.Equal(t, nil, err)

  // test if everything will be unpacked
  for i := 0; i < len(TestArrays); i++ {
    res, err := UnpackLE(conn)
    assert.Equal(t, nil, err)
    assert.Equal(t, TestArrays[i], res)
  }

  // EOF should be found here
  _, err = UnpackLE(conn)
  assert.Equal(t, io.EOF, err)
  // close the socket
  assert.Equal(t, nil, conn.Close())
}