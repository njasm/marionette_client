package marionette_client

import (
    "net"
    "errors"
    "strings"
    "strconv"
    "encoding/json"
)

type transport struct {
    applicationType         string
    marionetteProtocol      int8
    conn                    net.Conn
}

type response []byte

func (t *transport) connect(host string, port int) (err error) {
    if t.conn != nil {
        errors.New("A Connection is already established. please disconnect before connecting.")
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

func (t *transport) send(command string, values interface{}) (response, error) {
    buf, err := t.transformCommand(command, values)
    if err != nil {
        return nil, err
    }

    _, err = write(t.conn, buf)
    if err != nil {
        return 0, err
    }

    return t.receive()
}

func (t *transport) receive() (response, error) {
    return read(t.conn, 4096)
}

func read(c net.Conn, size int) (response, error) {
    if size <= 0 {
        size = 4096
    }

    readBuf := make([]byte, size)
    _, err := c.Read(readBuf)
    if err != nil {
        return nil, err
    }

    return readBuf, nil
}

func (t *transport) transformCommand(command string, values interface{}) ([]byte, error) {
    var data map[string]interface{}
    if t.marionetteProtocol == "2" {
        data["name"] = command
        data["parameters"] = values
    }

    bytes, err := json.Marshal(data);
    if err != nil {
        return nil, err
    }

    var size = len(bytes)
    return []byte(strconv.Itoa(size) + ":" + string(bytes))
}

func write(c net.Conn, b []byte) (int, error) {
    return c.Write(b)
}