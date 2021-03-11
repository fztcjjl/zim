package main

import (
	"encoding/binary"
	"io"
	"log"
)

/// +-------------------+------------------+------------------+
/// | prependable bytes |  readable bytes  |  writable bytes  |
/// |                   |     (CONTENT)    |                  |
/// +-------------------+------------------+------------------+
/// |                   |                  |                  |
/// 0      <=      readerIndex   <=   writerIndex    <=     size

const (
	kCheapPrepend = 8
	kInitialSize  = 1024
)

type Buffer struct {
	buf []byte
	r   int
	w   int
}

func NewBuffer(size int) *Buffer {
	b := new(Buffer)
	b.r = kCheapPrepend
	b.w = kCheapPrepend
	if size == 0 {
		size = kInitialSize
	}
	b.buf = make([]byte, size)

	return b
}

func (b *Buffer) Readable() int {
	return b.w - b.r
}

func (b *Buffer) Writeable() int {
	return cap(b.buf) - b.w
}

func (b *Buffer) Prependable() int {
	return b.r
}

func (b *Buffer) Peek() []byte {
	return b.buf[b.r:b.w]
}

func (b *Buffer) Retrieve(n int) {
	if n > b.Readable() {
		log.Fatal("Retrieve")
	}
	if n < b.Readable() {
		b.r += n
	} else {
		b.RetrieveAll()
	}
}

func (b *Buffer) retrieveUint64() {
	b.Retrieve(8)
}

func (b *Buffer) retrieveUint32() {
	b.Retrieve(4)
}

func (b *Buffer) retrieveUint16() {
	b.Retrieve(2)
}

func (b *Buffer) retrieveUint8() {
	b.Retrieve(1)
}

func (b *Buffer) RetrieveAll() {
	b.r = kCheapPrepend
	b.w = kCheapPrepend
}

func (b *Buffer) RetrieveAllAsString() string {
	result := string(b.Peek())
	b.RetrieveAll()
	return result
}

func (b *Buffer) Append(data []byte) {
	b.EnsureWritable(len(data))
	copy(b.buf[b.w:], data)
	b.HasWritten(len(data))

}

func (b *Buffer) EnsureWritable(n int) {
	if b.Writeable() < n {
		b.makeSpace(n)
	}

	if b.Writeable() < n {
		log.Fatal("EnsureWritable")
	}
}

func (b *Buffer) HasWritten(n int) {
	if n > b.Writeable() {
		log.Fatal("HasWritten")
	}
	b.w += n
}

func (b *Buffer) Unwrite(n int) {
	if n > b.Readable() {
		log.Fatal("Unwrite")
	}
	b.w -= n
}

func (b *Buffer) WriteUint64(x uint64) {
	b.EnsureWritable(2)
	buf := make([]byte, 2)
	binary.BigEndian.PutUint64(buf, x)
	b.Append(buf)
}

func (b *Buffer) WriteUint32(x uint32) {
	b.EnsureWritable(2)
	buf := make([]byte, 2)
	binary.BigEndian.PutUint32(buf, x)
	b.Append(buf)
}

func (b *Buffer) WriteUnt16(x uint16) {
	b.EnsureWritable(2)
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, x)
	b.Append(buf)
}

func (b *Buffer) WriteUint8(x uint8) {
	b.EnsureWritable(1)
	b.Append([]byte{byte(x)})
}

func (b *Buffer) ReadUint64() uint64 {
	result := b.PeekUint64()
	b.retrieveUint64()
	return result
}

func (b *Buffer) ReadUint32() uint32 {
	result := b.PeekUint32()
	b.retrieveUint32()
	return result
}

func (b *Buffer) ReadUint16() uint16 {
	result := b.PeekUint16()
	b.retrieveUint16()
	return result
}

func (b *Buffer) ReadUint8() uint8 {
	result := b.PeekUint8()
	b.retrieveUint8()
	return result
}

func (b *Buffer) PeekUint64() uint64 {
	if b.Readable() < 8 {
		log.Fatal("PeekUint64")
	}

	buf := b.Peek()
	return binary.BigEndian.Uint64(buf[:8])
}

func (b *Buffer) PeekUint32() uint32 {
	if b.Readable() < 4 {
		log.Fatal("PeekUint32")
	}

	buf := b.Peek()
	return binary.BigEndian.Uint32(buf[:4])
}

func (b *Buffer) PeekUint16() uint16 {
	if b.Readable() < 2 {
		log.Fatal("PeekUint16")
	}
	buf := b.Peek()
	return binary.BigEndian.Uint16(buf[:2])
}

func (b *Buffer) PeekUint8() uint8 {
	if b.Readable() < 1 {
		log.Fatal("PeekUint8")
	}
	buf := b.Peek()

	return uint8(buf[0])
}

func (b *Buffer) PrependUint64(x uint64) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, x)
	b.Prepend(buf)
}

func (b *Buffer) PrependUint32(x uint32) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, x)
	b.Prepend(buf)
}

func (b *Buffer) PrependUint16(x uint16) {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, x)
	b.Prepend(buf)
}

func (b *Buffer) PrependUint8(x uint8) {
	b.Prepend([]byte{byte(x)})
}

func (b *Buffer) Prepend(data []byte) {
	log.Fatal(len(data) > b.Prependable())
	b.r -= len(data)
	copy(b.buf[:b.r], data)
}

func (b *Buffer) Shrink(reserve int) {
	readable := b.Readable()
	buf := make([]byte, kCheapPrepend+readable+reserve)
	copy(buf[:kCheapPrepend], b.buf[b.r:b.w])
	b.buf = buf
	b.r = kCheapPrepend
	b.w = b.r + readable
}

func (b *Buffer) Capacity() int {
	return cap(b.buf)
}

func (b *Buffer) Read(reader io.Reader) (int, error) {
	n, err := reader.Read(b.buf[b.w:])
	if err != nil {
		return n, err
	}
	b.w += n

	return n, nil
}

func (b *Buffer) makeSpace(n int) {
	if b.Writeable()+b.Prependable() < n+kCheapPrepend {
		buf := make([]byte, b.w+n)
		copy(buf, b.buf)
		b.buf = buf
	} else {
		if kCheapPrepend < b.r {
			log.Fatal("makeSpace error")
		}
		readable := b.Readable()

		copy(b.buf[:kCheapPrepend], b.buf[b.r:b.w])
		b.r = kCheapPrepend
		b.w = b.r + readable
	}
}
