package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type SimpleTelnetClient struct {
	Address string
	Timeout time.Duration
	In      io.ReadCloser
	Out     io.Writer
	Conn    net.Conn
}

func (s *SimpleTelnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", s.Address, s.Timeout)
	if err != nil {
		return fmt.Errorf("Connect: %w", err)
	}
	s.Conn = conn

	return nil
}

func (s *SimpleTelnetClient) Close() error {
	if s.Conn != nil {
		if err := s.Conn.Close(); err != nil {
			return fmt.Errorf("Close: %w", err)
		}
	}
	return nil
}

func (s *SimpleTelnetClient) Send() error {
	// s.Out.Write([]byte("I'm telnet client\n"))
	if _, err := io.Copy(s.Conn, s.In); err != nil {
		return fmt.Errorf("Send: %w", err)
	}

	// s.Out.Write([]byte("I will be back!\n"))
	return nil
}

func (s *SimpleTelnetClient) Receive() error {
	if _, err := io.Copy(s.Out, s.Conn); err != nil {
		return fmt.Errorf("Receive: %w", err)
	}

	// s.Out.Write([]byte("Bye-bye\n"))
	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	// Place your code here.
	return &SimpleTelnetClient{
		Address: address,
		Timeout: timeout,
		In:      in,
		Out:     out,
	}
}
