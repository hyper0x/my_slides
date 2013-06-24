package nio

import (
	"bufio"
	"go_lib/logging"
	"net"
)

var logger logging.Logger = logging.GetSimpleLogger()

func Logger() logging.Logger {
	return logger
}

func ReadFromTcp(conn net.Conn, delim byte) (string, error) {
	reader := bufio.NewReader(conn)
	content, err := reader.ReadString(delim)
	if err != nil {
		return "", err
	} else {
		return content, nil
	}
}

func WriteToTcp(conn net.Conn, content string) (int, error) {
	writer := bufio.NewWriter(conn)
	number, err := writer.WriteString(content)
	if err == nil {
		err = writer.Flush()
	}
	return number, err
}
