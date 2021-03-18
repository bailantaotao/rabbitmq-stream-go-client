package stream

import (
	"bufio"
	"encoding/binary"
	"io"
)

func ReadUShort(readerStream io.Reader) uint16 {
	var res uint16
	_ = binary.Read(readerStream, binary.BigEndian, &res)
	return res
}

func ReadUInt(readerStream io.Reader) uint32 {
	var res uint32
	_ = binary.Read(readerStream, binary.BigEndian, &res)
	return res
}

func PeekByte(readerStream *bufio.Reader) uint8 {
	res, _ := readerStream.Peek(1)
	return res[0]
}

func ReadInt(readerStream io.Reader) uint {
	var res uint
	_ = binary.Read(readerStream, binary.BigEndian, &res)
	return res
}

func ReadInt64(readerStream io.Reader) int64 {
	var res int64
	_ = binary.Read(readerStream, binary.BigEndian, &res)
	return res
}

func ReadByte(readerStream io.Reader) uint8 {
	var res uint8
	_ = binary.Read(readerStream, binary.BigEndian, &res)
	return res
}

func ReadString(readerStream io.Reader) string {
	lenString := ReadUShort(readerStream)
	buff := make([]byte, lenString)
	_ = binary.Read(readerStream, binary.BigEndian, &buff)
	return string(buff)
}

func ReadUint8Array(readerStream io.Reader, size uint32) []uint8 {
	var res = make([]uint8, size)
	_ = binary.Read(readerStream, binary.BigEndian, &res)
	return res
}
