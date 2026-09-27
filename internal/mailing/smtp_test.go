package mailing

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/serge-masiutin/go-template/internal/config"
)

func TestSMTPDeliveryThroughGoMail(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	received := make(chan string, 1)
	failures := make(chan error, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			failures <- err
			return
		}
		defer conn.Close()
		conn.SetDeadline(time.Now().Add(5 * time.Second))
		session := textproto.NewConn(conn)
		if err := session.PrintfLine("220 fixture ESMTP"); err != nil {
			failures <- err
			return
		}
		for {
			line, err := session.ReadLine()
			if err != nil {
				failures <- err
				return
			}
			switch {
			case strings.HasPrefix(line, "EHLO"), strings.HasPrefix(line, "HELO"), strings.HasPrefix(line, "MAIL FROM:"), strings.HasPrefix(line, "RCPT TO:"), line == "RSET", line == "NOOP":
				err = session.PrintfLine("250 OK")
			case line == "DATA":
				if err = session.PrintfLine("354 send message"); err != nil {
					failures <- err
					return
				}
				body, readErr := session.ReadDotBytes()
				if readErr != nil {
					failures <- readErr
					return
				}
				received <- string(body)
				err = session.PrintfLine("250 accepted")
			case line == "QUIT":
				session.PrintfLine("221 bye")
				return
			default:
				failures <- fmt.Errorf("unexpected SMTP command %q", line)
				return
			}
			if err != nil {
				failures <- err
				return
			}
		}
	}()
	_, port, _ := net.SplitHostPort(listener.Addr().String())
	number, _ := strconv.Atoi(port)
	sender, err := New(config.Mail{Host: "127.0.0.1", Port: number, From: "starter@example.test", TLS: "none", Timeout: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	if err := sender.Send(context.Background(), "recipient@example.test", "Your notes", "A private note."); err != nil {
		t.Fatal(err)
	}
	select {
	case message := <-received:
		for _, expected := range []string{"recipient@example.test", "Your notes", "A private note."} {
			if !strings.Contains(message, expected) {
				t.Fatalf("message lacks %q", expected)
			}
		}
	case err := <-failures:
		t.Fatal(err)
	case <-time.After(2 * time.Second):
		t.Fatal("no message received")
	}
}

func TestSMTPCancellationClosesStalledGreeting(t *testing.T) {
	for _, policy := range []string{"none", "starttls", "tls"} {
		t.Run(policy, func(t *testing.T) {
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			defer listener.Close()
			accepted := make(chan struct{})
			closed := make(chan error, 1)
			go func() {
				connection, err := listener.Accept()
				if err != nil {
					closed <- err
					return
				}
				defer connection.Close()
				connection.SetDeadline(time.Now().Add(2 * time.Second))
				close(accepted)
				_, err = io.Copy(io.Discard, connection)
				closed <- err
			}()
			_, port, _ := net.SplitHostPort(listener.Addr().String())
			number, _ := strconv.Atoi(port)
			sender, err := New(config.Mail{Host: "127.0.0.1", Port: number, From: "starter@example.test", TLS: policy, Timeout: time.Minute})
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			sent := make(chan error, 1)
			go func() { sent <- sender.Send(ctx, "recipient@example.test", "Test", "Synthetic content") }()
			select {
			case <-accepted:
			case <-time.After(time.Second):
				t.Fatal("SMTP not connected")
			}
			cancel()
			select {
			case err := <-sent:
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("cancellation lost: %v", err)
				}
			case <-time.After(time.Second):
				t.Fatal("canceled SMTP delivery did not stop")
			}
			if err := <-closed; err != nil {
				t.Fatalf("connection was not closed: %v", err)
			}
		})
	}
}
