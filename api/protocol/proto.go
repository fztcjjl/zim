package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// |-----------|---------------|---------|---------|---------|---------------|
// | HeaderLen | ClientVersion |   Cmd   |   Seq   | BodyLen |     Body      |
// |-----------|---------------|---------|---------|---------|---------------|
// |  4 bytes  |    4 bytes    | 4 bytes | 4 bytes | 4 bytes | BodyLen bytes |
// |-----------|---------------|---------|---------|---------|---------------|
// |                        16 bytes                         |
// |---------------------------------------------------------|

const (
	HeaderLen = 20
)

type Proto struct {
	HeaderLen     uint32
	ClientVersion uint32
	Cmd           uint32
	Seq           uint32
	BodyLen       uint32
	Body          []byte
}

func (p *Proto) Read(data []byte) (err error) {
	if len(data) < HeaderLen {
		err = fmt.Errorf("packet error")
		return
	}
	buf := &bytes.Buffer{}
	buf.Write(data)
	if err = binary.Read(buf, binary.BigEndian, &p.HeaderLen); err != nil {
		return
	}
	if err = binary.Read(buf, binary.BigEndian, &p.ClientVersion); err != nil {
		return
	}

	if err = binary.Read(buf, binary.BigEndian, &p.Cmd); err != nil {
		return
	}

	if err = binary.Read(buf, binary.BigEndian, &p.Seq); err != nil {
		return
	}
	if err = binary.Read(buf, binary.BigEndian, &p.BodyLen); err != nil {
		return
	}
	if p.BodyLen > 0 {
		body := make([]byte, p.BodyLen)
		if _, err = buf.Read(body); err != nil {
			return
		}
		p.Body = body
	}

	return
}

func (p *Proto) Write(buf *bytes.Buffer) (err error) {
	p.HeaderLen = HeaderLen
	p.ClientVersion = 1
	p.BodyLen = uint32(len(p.Body))
	if err = binary.Write(buf, binary.BigEndian, p.HeaderLen); err != nil {
		return
	}

	if err = binary.Write(buf, binary.BigEndian, p.ClientVersion); err != nil {
		return
	}

	if err = binary.Write(buf, binary.BigEndian, p.Cmd); err != nil {
		return
	}

	if err = binary.Write(buf, binary.BigEndian, p.Seq); err != nil {
		return
	}

	if err = binary.Write(buf, binary.BigEndian, p.BodyLen); err != nil {
		return
	}

	if p.Body != nil {
		_, err = buf.Write(p.Body)
	}

	return
}
