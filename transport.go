package marionette_client

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type Transporter interface {
	MessageID() int
	Connect(host string, port int) error
	Close() error
	Send(command string, values interface{}) (*Response, error)
	Receive() ([]byte, error)
}

type MarionetteTransport struct {
	ApplicationType    string
	MarionetteProtocol int32
	messageID          int
	conn               net.Conn
	de                 DecoderEncoder
}

type Response struct {
	MessageID   int32
	Size        int32
	Value       string
	DriverError *DriverError
}

func connDefaultTimeout() time.Time {
	return time.Now().Add(time.Minute * 5)
}

func (t *MarionetteTransport) MessageID() int {
	return t.messageID
}

func (t *MarionetteTransport) Connect(host string, port int) error {
	if t.conn != nil {
		return errors.New("a connection is already established. please disconnect before connecting")
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
	r, err := t.Receive()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(r), &t)
	if err != nil {
		return err
	}

	d, err := NewDecoderEncoder(t.MarionetteProtocol)
	if err != nil {
		return err
	}

	t.de = d

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

func (t *MarionetteTransport) Send(command string, values interface{}) (*Response, error) {
	t.messageID = t.messageID + 1 // next message ID
	buf, err := t.de.Encode(t, command, values)
	if err != nil {
		return nil, err
	}

	_, err = write(t.conn, buf)
	if err != nil {
		return nil, err
	}

	rBuf, err := t.Receive()
	if err != nil {
		return nil, err
	}

	//Debug only
	if RunningInDebugMode {
		if len(buf) >= 512 {
			log.Println(string(buf)[0:512] + " - END - " + string(buf)[len(buf)-512:])
		} else {
			log.Println(string(buf))
		}
	}
	//Debug only end

	data := &Response{}
	err = t.de.Decode(rBuf, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func write(c net.Conn, b []byte) (int, error) {
	err := c.SetDeadline(connDefaultTimeout())
	if err != nil {
		return 0, err
	}

	return c.Write(b)
}

func (t *MarionetteTransport) Receive() ([]byte, error) {
	err := t.conn.SetDeadline(connDefaultTimeout())
	if err != nil {
		return nil, err
	}

	return read(t.conn)
}

// ReadFull reads exactly len(buf) bytes from r into buf.
// It returns the number of bytes copied and an error if fewer bytes were read.
// The error is EOF only if no bytes were read.
// If an EOF happens after reading some but not all the bytes,
// ReadFull returns ErrUnexpectedEOF.
// On return, n == len(buf) if and only if err == nil.
func read(c net.Conn) ([]byte, error) {
	var msgSize, err = messageLength(c)
	if err != nil {
		return nil, err
	}

	msgBuf := make([]byte, msgSize)
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
func messageLength(c net.Conn) (int, error) {
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
