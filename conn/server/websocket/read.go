package websocket

import (
	"encoding/binary"
	"fmt"
	"github.com/panjf2000/gnet"
)

// Errors used by frame reader.
var (
	ErrHeaderLengthMSB        = fmt.Errorf("header error: the most significant bit must be 0")
	ErrHeaderLengthUnexpected = fmt.Errorf("header error: unexpected payload length bits")
)

func ReadHeader(conn gnet.Conn) (h Header, err error) {
	buf := conn.Read()
	if len(buf) < 6 {
		err = fmt.Errorf("header error: not enough")
		return
	}
	// Make slice of bytes with capacity 12 that could hold any header.
	//
	// The maximum header size is 14, but due to the 2 hop reads,
	// after first hop that reads first 2 constant bytes, we could reuse 2 bytes.
	// So 14 - 2 = 12.
	//bts := make([]byte, 2, MaxHeaderSize-2)

	// Prepare to hold first 2 bytes to choose size of next read.
	//_, err = io.ReadFull(r, bts)
	//if err != nil {
	//	return
	//}

	bts := buf[:2]

	h.Fin = bts[0]&bit0 != 0
	h.Rsv = (bts[0] & 0x70) >> 4
	h.OpCode = OpCode(bts[0] & 0x0f)

	var extra int

	if bts[1]&bit0 != 0 {
		h.Masked = true
		extra += 4
	}

	length := bts[1] & 0x7f
	switch {
	case length < 126:
		h.Length = int64(length)

	case length == 126:
		extra += 2

	case length == 127:
		extra += 8

	default:
		err = ErrHeaderLengthUnexpected
		return
	}

	if extra == 0 {
		return
	}

	if len(buf) < 2+extra {
		return
	}
	// Increase len of bts to extra bytes need to read.
	// Overwrite first 2 bytes that was read before.
	//bts = bts[:extra]
	//_, err = io.ReadFull(r, bts)
	//if err != nil {
	//	return
	//}

	bts = buf[2 : 2+extra]

	switch {
	case length == 126:
		h.Length = int64(binary.BigEndian.Uint16(bts[:2]))
		bts = bts[2:]

	case length == 127:
		if bts[0]&0x80 != 0 {
			err = ErrHeaderLengthMSB
			return
		}
		h.Length = int64(binary.BigEndian.Uint64(bts[:8]))
		bts = bts[8:]
	}

	if h.Masked {
		copy(h.Mask[:], bts)
	}

	conn.ShiftN(2 + extra)

	return
}
