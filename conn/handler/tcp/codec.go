package tcp

import (
	"bytes"
	"encoding/binary"
	"errors"
	log "github.com/fztcjjl/tiger/trpc/logger"
	"github.com/fztcjjl/zim/api/protocol"
	"github.com/panjf2000/gnet"
)

type Codec struct {
}

func (_ *Codec) Encode(c gnet.Conn, buf []byte) ([]byte, error) {
	return buf, nil
}

func (_ *Codec) Decode(c gnet.Conn) ([]byte, error) {
	if size, header := c.ReadN(protocol.HeaderLen); size == protocol.HeaderLen {
		byteBuffer := bytes.NewBuffer(header)
		var p protocol.Proto
		if err := binary.Read(byteBuffer, binary.BigEndian, &p.HeaderLen); err != nil {
			return nil, err
		}
		if err := binary.Read(byteBuffer, binary.BigEndian, &p.ClientVersion); err != nil {
			return nil, err
		}
		if err := binary.Read(byteBuffer, binary.BigEndian, &p.Cmd); err != nil {
			return nil, err
		}
		if err := binary.Read(byteBuffer, binary.BigEndian, &p.Seq); err != nil {
			return nil, err
		}
		if err := binary.Read(byteBuffer, binary.BigEndian, &p.BodyLen); err != nil {
			return nil, err
		}

		protocolLen := int(protocol.HeaderLen + p.BodyLen)
		if size, data := c.ReadN(protocolLen); size == protocolLen {
			c.ShiftN(protocolLen)
			log.Info(data)
			log.Info(len(data), size, protocolLen)
			return data, nil
		}
		return nil, errors.New("not enough payload data")
	}

	return nil, errors.New("not enough header data")
}
