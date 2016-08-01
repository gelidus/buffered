package buffered

import (
  "io"
  "encoding/binary"
)

// UInt32ToBytesConverter is a function prototype that
// should be able to transfer given uint32 number into
// a set of bytes
type UInt32ToBytesConverter func([]byte, uint32)

//BytesToUInt32Converter is a function prototype that
// should be able to transfer given bytes slice into
// a uint32 number
type BytesToUInt32Converter func([]byte) uint32

// Pack returns a slice of bytes prefixed with its length.
// Length is converted into bytes by UInt32ToBytesConverter.
// This converter is able to consume interface of standard
// binary converters.
func Pack(bytes []byte, converter UInt32ToBytesConverter) []byte {
  var length []byte = make([]byte, 4)
  converter(length, uint32(len(bytes)))

  return append(length, bytes...)
}

// PackLE returns a slice of bytes prefixed with its length.
// It uses LittleEndian context of standard binary conversion library
func PackLE(bytes []byte) []byte {
  return Pack(bytes, binary.LittleEndian.PutUint32)
}

// PackBE returns a slice of bytes prefixed with its length.
// Ut uses BigEndian context of standard binary conversion library
func PackBE(bytes []byte) []byte {
  return Pack(bytes, binary.BigEndian.PutUint32)
}

// Unpack returns a slice of read bytes from the given reader.
// It uses given converter to get the total length of data. Only
// errors caused by Read operations on reader can be returned
func Unpack(reader io.Reader, converter BytesToUInt32Converter) ([]byte, error) {
  var lengthBytes []byte = make([]byte, 4)
  _, err := reader.Read(lengthBytes)
  if err != nil {
    return []byte{}, err
  }

  length := converter(lengthBytes)
  var contentBytes []byte = make([]byte, length)
  _, err = reader.Read(contentBytes)

  return contentBytes, err
}

// UnpackLE uses Unpack with predefined little endian uint32
// conversion
func UnpackLE(reader io.Reader) ([]byte, error) {
  return Unpack(reader, binary.LittleEndian.Uint32)
}

// UnpackBE uses Unpack with predefined big endian uint32
// conversion
func UnpackBE(reader io.Reader) ([]byte, error) {
  return Unpack(reader, binary.BigEndian.Uint32)
}