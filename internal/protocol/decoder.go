package protocol

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"io"
)

// LigoloDecoder is the struct containing the decoded Envelope and the reader
type LigoloDecoder struct {
	reader   io.Reader
	Envelope Envelope
}

// NewDecoder decode Ligolo-ng packets
func NewDecoder(reader io.Reader) LigoloDecoder {
	return LigoloDecoder{reader: reader}
}

// Decode read content from the reader and fill the Envelope
func (d *LigoloDecoder) Decode() error {
	if err := binary.Read(d.reader, binary.LittleEndian, &d.Envelope.Type); err != nil {
		return err
	}

	if err := binary.Read(d.reader, binary.LittleEndian, &d.Envelope.Size); err != nil {
		return err
	}

	payload := make([]byte, d.Envelope.Size)

	if _, err := d.reader.Read(payload); err != nil {
		return err
	}

	gobdecoder := gob.NewDecoder(bytes.NewReader(payload))

	// Kind of dirty, but it's the only way I found to satisfy gob
	switch d.Envelope.Type {
	case MessageInfoRequest:
		p := InfoRequestPacket{}
		if err := gobdecoder.Decode(&p); err != nil {
			panic(err)
		}
		d.Envelope.Payload = p
	case MessageInfoReply:
		p := InfoReplyPacket{}
		if err := gobdecoder.Decode(&p); err != nil {
			panic(err)
		}
		d.Envelope.Payload = p
	case MessageConnectRequest:
		p := ConnectRequestPacket{}
		if err := gobdecoder.Decode(&p); err != nil {
			panic(err)
		}
		d.Envelope.Payload = p
	case MessageConnectResponse:
		p := ConnectResponsePacket{}
		if err := gobdecoder.Decode(&p); err != nil {
			panic(err)
		}
		d.Envelope.Payload = p
	case MessageHostPingRequest:
		p := HostPingRequestPacket{}
		if err := gobdecoder.Decode(&p); err != nil {
			panic(err)
		}
		d.Envelope.Payload = p
	case MessageHostPingResponse:
		p := HostPingResponsePacket{}
		if err := gobdecoder.Decode(&p); err != nil {
			panic(err)
		}
		d.Envelope.Payload = p
	case MessageRedirectorRequest:
		p := RedirectorRequestPacket{}
		if err := gobdecoder.Decode(&p); err != nil {
			panic(err)
		}
		d.Envelope.Payload = p
	case MessageRedirectorResponse:
		p := RedirectorResponsePacket{}
		if err := gobdecoder.Decode(&p); err != nil {
			panic(err)
		}
		d.Envelope.Payload = p
	case MessageRedirectorCloseRequest:
		p := RedirectorCloseRequestPacket{}
		if err := gobdecoder.Decode(&p); err != nil {
			panic(err)
		}
		d.Envelope.Payload = p
	case MessageRedirectorCloseResponse:
		p := RedirectorCloseResponsePacket{}
		if err := gobdecoder.Decode(&p); err != nil {
			panic(err)
		}
		d.Envelope.Payload = p
	case MessageDisconnectRequest:
		p := DisconnectRequestPacket{}
		if err := gobdecoder.Decode(&p); err != nil {
			panic(err)
		}
		d.Envelope.Payload = p
	case MessageDisconnectResponse:
		p := DisconnectResponsePacket{}
		if err := gobdecoder.Decode(&p); err != nil {
			panic(err)
		}
		d.Envelope.Payload = p
	default:
		return errors.New("invalid message type")
	}

	return nil
}
