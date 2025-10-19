package errors

import (
	"fmt"
	"net"
	"os"
	"syscall"

	"github.com/cockroachdb/errors"
)

func PanicRecovered(recover any) error {

	if r := recover; r != nil {

		var err error
		switch e := r.(type) {
		case error:
			err = e
		case string:
			err = errors.New(e)
		default:
			err = fmt.Errorf("panic: %v", e)
		}

		err = errors.WithStackDepth(err, 2)
		return err
	}

	return nil
}

func RequestRecovered(err error) error {

	switch {
	case os.IsTimeout(err):
		return errors.WithDetail(err, "request timeout")
	case errors.Is(err, syscall.ECONNREFUSED):
		return errors.WithDetail(err, "connection refused")
	case errors.Is(err, syscall.ECONNRESET):
		return errors.WithDetail(err, "connection reset")
	default:
		var dnsErr *net.DNSError
		if errors.As(err, &dnsErr) {
			return errors.WithDetail(err, "DNS error")
		}
		return errors.WithDetail(err, "network error")
	}
}
