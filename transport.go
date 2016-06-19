package marionette_client

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"strconv"
	"time"
)

type Transporter interface {
	Connect(host string, port int) error
	Close() error
	Send(command string, values interface{}) (*response, error)
	Receive() ([]byte, error)
}

type MarionetteTransport struct {
	ApplicationType    string
	MarionetteProtocol int32
	messageID          int
	conn               net.Conn
}

type response struct {
	MessageID   int32
	Size        int32
	Value       string
	DriverError *DriverError
}

func (t *MarionetteTransport) Connect(host string, port int) error {
	if t.conn != nil {
		return errors.New("A Connection is already established. please disconnect before connecting.")
	}

	if host == "" {
		host = "127.0.0.1"
	}

	if port == 0 {
		port = 2828
	}

	hostname := host + ":" + strconv.Itoa(port)
	c, err := net.Dial("tcp", hostname)
	if err != nil {
		return err
	}

	t.conn = c
	t.conn.SetDeadline(time.Now().Add(time.Minute * 5)) // default read and write time out
	r, err := t.Receive()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(r), &t)
	if err != nil {
		return err
	}

	return nil
}

func (t *MarionetteTransport) Close() error {
	err := t.conn.Close()
	if err != nil {
		return err
	}

	t.conn = nil
	return err
}

func (t *MarionetteTransport) Send(command string, values interface{}) (*response, error) {
	t.messageID = t.messageID + 1
	buf, err := t.transformToCommand(command, values)
	if err != nil {
		return nil, err
	}

	_, err = write(t.conn, buf)
	if err != nil {
		return nil, err
	}

	// get response to sent command.
	return t.transformToResponse(t.Receive())
}

func write(c net.Conn, b []byte) (int, error) {
	return c.Write(b)
}

func (t *MarionetteTransport) Receive() ([]byte, error) {
	return read(t.conn)
}

func read(c net.Conn) ([]byte, error) {
	var msgSize, err = getMessageLength(c)
	if err != nil {
		return nil, err
	}

	msgBuf := make([]byte, msgSize)
	//_, err = c.Read(msgBuf)
	// ReadFull reads exactly len(buf) bytes from r into buf.
	// It returns the number of bytes copied and an error if fewer bytes were read.
	// The error is EOF only if no bytes were read.
	// If an EOF happens after reading some but not all the bytes,
	// ReadFull returns ErrUnexpectedEOF.
	// On return, n == len(buf) if and only if err == nil.
	_, err = io.ReadFull(c, msgBuf)
	if err != nil {
		return nil, err
	}

	return msgBuf, nil
}

// Reads from the connection byte by byte until the message length is found, according to
// marionette's protocol.
// the protocol say's that message length is the first part for the message until ":" is found.
// this signals the next bytes as the message
func getMessageLength(c net.Conn) (int, error) {
	var byteSize = make([]byte, 0)
	tmp := make([]byte, 1)
	for {
		_, err := c.Read(tmp)
		if err != nil {
			if err != io.EOF {
				return 0, err
			}
		}

		if string(tmp) != ":" {
			byteSize = append(byteSize, tmp...)
			continue
		}

		// the message length
		intSize, err := strconv.Atoi(string(byteSize))
		if err != nil {
			return 0, err
		}

		return intSize, err
	}
}

func (t *MarionetteTransport) transformToCommand(command string, values interface{}) (bytes []byte, err error) {
	var size int
	if t.MarionetteProtocol == MARIONETTE_PROTOCOL_V2 {
		bytes, err = makeProto2Command(command, values)
	} else if t.MarionetteProtocol == MARIONETTE_PROTOCOL_V3 {
		bytes, err = makeProto3Command(t.messageID, command, values)
	} else {
		return nil, errors.New("Marionete Protocol version not supported.")
	}

	if err != nil {
		return nil, err
	}

	size = len(bytes)
	return []byte(strconv.Itoa(size) + ":" + string(bytes)), nil
}

func (t *MarionetteTransport) transformToResponse(buf []byte, err error) (*response, error) {
	if err != nil {
		return nil, err
	}

	if t.MarionetteProtocol == MARIONETTE_PROTOCOL_V2 {
		return makeProto2Response(buf)
	} else if t.MarionetteProtocol == MARIONETTE_PROTOCOL_V3 {
		return makeProto3Response(buf)
	}

	return nil, errors.New("Unable to decode Protocol version for message decoding.")
}
