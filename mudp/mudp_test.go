package mudp

import (
	"testing"
	"time"
)

func TestMudpSendRecv(t *testing.T) {
	srv, err := NewServer(&Server{IP: "127.0.0.1", Port: 19000})
	if err != nil {
		t.Fatalf("NewServer err: %v", err)
	}

	recvCh := make(chan []byte, 1)
	srv.OnMessage = func(remote string, data []byte) {
		recvCh <- data
	}

	if err := srv.Start(); err != nil {
		t.Fatalf("Start err: %v", err)
	}
	defer srv.Stop()

	snd, err := NewSender(Sender{IP: "127.0.0.1", Port: 19000})
	if err != nil {
		t.Fatalf("NewSender err: %v", err)
	}
	defer snd.Close()

	msg := []byte("hello-mudp")
	n, err := snd.Write(msg)
	if err != nil {
		t.Fatalf("Write err: %v", err)
	}
	if n != len(msg) {
		t.Fatalf("Write len mismatch got %d want %d", n, len(msg))
	}

	select {
	case d := <-recvCh:
		if string(d) != string(msg) {
			t.Fatalf("recv data mismatch: %s", string(d))
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting for message")
	}
}
