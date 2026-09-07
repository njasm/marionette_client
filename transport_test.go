package marionette_client

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"reflect"
	"testing"
)

func TestMarionetteTransportConnectAndClose(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = listener.Close() }()

	serverError := make(chan error, 1)
	go func() {
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			serverError <- acceptErr
			return
		}
		defer func() { _ = conn.Close() }()
		greeting := `{"applicationType":"gecko","marionetteProtocol":3}`
		if _, writeErr := fmt.Fprintf(conn, "%d:%s", len(greeting), greeting); writeErr != nil {
			serverError <- writeErr
			return
		}
		_, copyErr := io.Copy(io.Discard, conn)
		serverError <- copyErr
	}()

	address := listener.Addr().(*net.TCPAddr)
	transport := &MarionetteTransport{}
	if err = transport.Connect(address.IP.String(), address.Port); err != nil {
		t.Fatal(err)
	}
	if transport.ApplicationType != "gecko" || transport.MarionetteProtocol != MarionetteProtocolV3 || transport.de == nil {
		t.Fatalf("unexpected handshake state: %#v", transport)
	}
	if err = transport.Connect(address.IP.String(), address.Port); err == nil {
		t.Fatal("expected an active-connection error")
	}
	if err = transport.Close(); err != nil {
		t.Fatal(err)
	}
	if transport.conn != nil {
		t.Fatal("connection was not cleared")
	}
	if err = <-serverError; err != nil {
		t.Fatal(err)
	}
}

func TestMarionetteTransportSend(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer func() { _ = clientConn.Close() }()
	defer func() { _ = serverConn.Close() }()

	serverError := make(chan error, 1)
	go func() {
		request, err := read(serverConn)
		if err != nil {
			serverError <- err
			return
		}
		var message []any
		if err = json.Unmarshal(request, &message); err != nil {
			serverError <- err
			return
		}
		expected := []any{float64(0), float64(1), "WebDriver:Test", map[string]any{"value": "sent"}}
		if !reflect.DeepEqual(message, expected) {
			serverError <- fmt.Errorf("unexpected request: %#v", message)
			return
		}
		response := `[1,1,null,{"value":"received"}]`
		_, err = fmt.Fprintf(serverConn, "%d:%s", len(response), response)
		serverError <- err
	}()

	transport := &MarionetteTransport{conn: clientConn, de: ProtoV3DecoderEncoder{}}
	response, err := transport.Send("WebDriver:Test", map[string]any{"value": "sent"})
	if err != nil {
		t.Fatal(err)
	}
	if response.MessageID != 1 || response.Value != `{"value":"received"}` || transport.MessageID() != 1 {
		t.Fatalf("unexpected response: %#v", response)
	}
	if err = <-serverError; err != nil {
		t.Fatal(err)
	}
}

func TestMarionetteFramingErrors(t *testing.T) {
	t.Run("invalid length", func(t *testing.T) {
		conn, peer := net.Pipe()
		defer func() { _ = conn.Close() }()
		defer func() { _ = peer.Close() }()
		go func() { _, _ = io.WriteString(peer, "invalid:") }()
		if _, err := read(conn); err == nil {
			t.Fatal("expected invalid message length to fail")
		}
	})

	t.Run("short body", func(t *testing.T) {
		conn, peer := net.Pipe()
		defer func() { _ = conn.Close() }()
		go func() {
			_, _ = io.WriteString(peer, "5:abc")
			_ = peer.Close()
		}()
		if _, err := read(conn); err != io.ErrUnexpectedEOF {
			t.Fatalf("expected unexpected EOF, got %v", err)
		}
	})

	t.Run("missing length", func(t *testing.T) {
		conn, peer := net.Pipe()
		defer func() { _ = conn.Close() }()
		_ = peer.Close()
		if _, err := messageLength(conn); err != io.EOF {
			t.Fatalf("expected EOF, got %v", err)
		}
	})
}

func TestWriteFramesRawBytes(t *testing.T) {
	conn, peer := net.Pipe()
	defer func() { _ = conn.Close() }()
	defer func() { _ = peer.Close() }()
	written := make(chan error, 1)
	go func() {
		_, err := write(conn, []byte("payload"))
		written <- err
	}()
	buf := make([]byte, len("payload"))
	if _, err := io.ReadFull(peer, buf); err != nil {
		t.Fatal(err)
	}
	if string(buf) != "payload" {
		t.Fatalf("unexpected bytes %q", buf)
	}
	if err := <-written; err != nil {
		t.Fatal(err)
	}
}
