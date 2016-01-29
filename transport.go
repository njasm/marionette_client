package marionette_client

import (
	"encoding/json"
	"errors"
	"net"
	"strconv"
    "io"
)

type transport struct {
	ApplicationType    string
	MarionetteProtocol int32
	conn               net.Conn
}

type response struct {
	Id    int32
	Size  int64
	Value string
}

func (t *transport) connect(host string, port int) (err error) {
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
	t.conn, err = net.Dial("tcp", hostname)
	if err != nil {
		return err
	}

	r, err := t.transformResponse(t.receive())
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(r.Value), &t)
	if err != nil {
		return err
	}

	return nil
}

func (t *transport) close() error {
	err := t.conn.Close()
	if err != nil {
		return err
	}

	t.conn = nil
	return err
}

func (t *transport) send(command string, values interface{}) (*response, error) {
	buf, err := t.transformCommand(command, values)
	if err != nil {
		return nil, err
	}

	_, err = write(t.conn, buf)
	if err != nil {
		return nil, err
	}

	return t.transformResponse(t.receive())
}

func (t *transport) receive() ([]byte, error) {
	return read(t.conn)
}

func read(c net.Conn) ([]byte, error) {
    var msgSize = make([]byte, 0)
    tmp := make([]byte, 1)
    for {
        _, err := c.Read(tmp)
        if err != nil {
            if err != io.EOF {
                return nil, err
            }
        }

        if string(tmp) == ":" {
            if size, err := strconv.Atoi(string(msgSize)); err == nil {
                msgBuf := make([]byte, size)
                _, err = c.Read(msgBuf)
                if err != nil {
                    return nil, err
                }

                return msgBuf, nil
            }

            return nil, err
        }

        msgSize = append(msgSize, tmp...)
    }
}

func (t *transport) transformCommand(command string, values interface{}) ([]byte, error) {
	data := make(map[string]interface{})
	if t.MarionetteProtocol == 2 {
		data["name"] = command
		data["parameters"] = values
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var size = len(bytes)
	return []byte(strconv.Itoa(size) + ":" + string(bytes)), nil
}

func write(c net.Conn, b []byte) (int, error) {
	return c.Write(b)
}

func (t *transport) transformResponse(buf []byte, err error) (*response, error) {
	if err != nil {
		return nil, err
	}

	stringBuf := string(buf)
    totalMessageLength := len(buf)

	if totalMessageLength != len(stringBuf) {
		return nil, errors.New("Total Message Length does not match with actual message length")
	}

	return &response{Size: int64(totalMessageLength), Value: stringBuf}, nil
}
