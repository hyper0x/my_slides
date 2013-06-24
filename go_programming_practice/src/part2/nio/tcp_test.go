package nio

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	DELIM byte = '\t'
)

func TestMainFuncs(t *testing.T) {
	serverAddr := "127.0.0.1:8080"
	var listener PubSubListener = NewTcpListener()
	t.Logf("Start Listening at address %s ...", serverAddr)
	err := listener.Init(serverAddr)
	if err != nil {
		t.Errorf("Listener: Init Error: %s\n", err)
		return
	}
	requestHandler := func(conn net.Conn) {
		for {
			content, err := ReadFromTcp(conn, DELIM)
			if err != nil {
				t.Errorf("Listener: Read Error: %s\n", err)
			} else {
				t.Logf("Listener: Received content: '%s'\n", content)
				content = strings.TrimSpace(content)
				if strings.HasSuffix(content, "q!") {
					t.Log("Listener: Quit!")
					break
				}
				resp := generateTestContent(fmt.Sprintf("Resp: %s", content))
				n, err := WriteToTcp(conn, resp)
				if err != nil {
					t.Errorf("Listener: Echo Error: %s\n", err)
				} else {
					t.Logf("Listener: Send response: '%s' (n=%d)\n", resp, n)
				}
			}
		}
	}
	err = listener.Listen(requestHandler)
	if err != nil {
		t.Errorf("Listener: Error: %s\n", err)
		t.FailNow()
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		multiSend("127.0.0.1:8080", "S1", 2, (2 * time.Second), t)
	}()
	go func() {
		defer wg.Done()
		multiSend("127.0.0.1:8080", "S2", 1, (2 * time.Second), t)
	}()
	wg.Wait()
	listener.Close()
}

func multiSend(remoteAddr string, clientName string, number int, timeout time.Duration, t *testing.T) {
	sender := NewTcpSender()
	t.Logf("Initializing sender (%s) (remote address: %s, timeout: %d) ...", clientName, remoteAddr, timeout)
	err := sender.Init(remoteAddr, timeout)
	if err != nil {
		t.Errorf("%s: Init Error: %s\n", clientName, err)
		return
	}
	if number <= 0 {
		number = 5
	}
	for i := 0; i < number; i++ {
		content := generateTestContent(fmt.Sprintf("%s-%d", clientName, i))
		t.Logf("%s: Send content: '%s'\n", clientName, content)
		err := sender.Send(content)
		if err != nil {
			t.Errorf("%s: Send Error: %s\n", clientName, err)
		}
		respChan := sender.Receive(DELIM)
		var resp PubSubContent
		timeoutChan := time.After(1 * time.Second)
		select {
		case resp = <-respChan:
		case <-timeoutChan:
			break
		}
		if err = resp.Err(); err != nil {
			t.Errorf("Sender: Receive Error: %s\n", err)
		} else {
			respContent := resp.Content()
			t.Logf("%s: Received response: '%s'\n", clientName, respContent)
		}
	}
	content := generateTestContent(fmt.Sprintf("%s-q!", clientName))
	t.Logf("%s: Send content: '%s'\n", clientName, content)
	err = sender.Send(content)
	if err != nil {
		t.Errorf("%s: Send Error: %s\n", clientName, err)
	}
	sender.Close()
}

func generateTestContent(content string) string {
	var respBuffer bytes.Buffer
	respBuffer.WriteString(strings.TrimSpace(content))
	respBuffer.WriteByte(DELIM)
	return respBuffer.String()
}
