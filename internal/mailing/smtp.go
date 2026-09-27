package mailing

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/serge-masiutin/go-template/internal/config"
	mail "github.com/wneessen/go-mail"
)

// SMTP sends a single message. It never retries an ambiguous delivery.
type SMTP struct {
	host        string
	from        string
	timeout     time.Duration
	implicitTLS bool
	options     []mail.Option
}

func New(cfg config.Mail) (*SMTP, error) {
	options := []mail.Option{mail.WithPort(cfg.Port), mail.WithTimeout(cfg.Timeout)}
	switch cfg.TLS {
	case "starttls":
		options = append(options, mail.WithTLSPolicy(mail.TLSMandatory))
	case "tls":
		options = append(options, mail.WithSSL())
	case "none":
		options = append(options, mail.WithTLSPolicy(mail.NoTLS))
	default:
		return nil, fmt.Errorf("unsupported mail TLS policy")
	}
	if cfg.Username != "" {
		options = append(options, mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithUsername(cfg.Username), mail.WithPassword(cfg.Password))
	}
	_, err := mail.NewClient(cfg.Host, options...)
	if err != nil {
		return nil, err
	}
	return &SMTP{host: cfg.Host, from: cfg.From, timeout: cfg.Timeout, implicitTLS: cfg.TLS == "tls", options: options}, nil
}

func (s *SMTP) Send(ctx context.Context, recipient, subject, body string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	// go-mail applies socket timeouts per exchange. Own the raw connection so the
	// whole delivery, including greeting/TLS/QUIT, stops when the job is canceled.
	var connection net.Conn
	var stopClose func() bool
	defer func() {
		if stopClose != nil {
			stopClose()
		}
		if connection != nil {
			connection.Close()
		}
	}()
	dial := func(dialCtx context.Context, network, address string) (net.Conn, error) {
		raw, err := (&net.Dialer{}).DialContext(dialCtx, network, address)
		if err != nil {
			return nil, err
		}
		connection = raw
		stopClose = context.AfterFunc(ctx, func() { raw.Close() })
		if s.implicitTLS {
			encrypted := tls.Client(raw, &tls.Config{ServerName: s.host, MinVersion: tls.VersionTLS12})
			if err := encrypted.HandshakeContext(ctx); err != nil {
				return nil, err
			}
			return encrypted, nil
		}
		return raw, nil
	}
	options := append(append([]mail.Option(nil), s.options...), mail.WithDialContextFunc(dial))
	client, err := mail.NewClient(s.host, options...)
	if err != nil {
		return err
	}
	message := mail.NewMsg()
	if err := message.From(s.from); err != nil {
		return err
	}
	if err := message.To(recipient); err != nil {
		return err
	}
	message.Subject(subject)
	message.SetBodyString(mail.TypeTextPlain, body)
	return errors.Join(client.DialAndSendWithContext(ctx, message), ctx.Err())
}
