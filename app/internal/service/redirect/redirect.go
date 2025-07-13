package redirect

import (
	"catcher/app/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
)

var errTimeout = errors.New("timeout exceeded")
var errSending = errors.New("sending error")

type Sender interface {
	Send() (*models.EventID, error)
}

type Report struct {
	models.AppContext
	ID   string
	Data []byte
}

func NewReport(id string, data []byte, appCtx models.AppContext) Report {
	return Report{
		ID:         id,
		Data:       data,
		AppContext: appCtx,
	}
}

func Send(sender Sender, timeout int) (*models.EventID, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	type result struct {
		eventID *models.EventID
		err     error
	}

	ch := make(chan result)

	go func() {

		defer func() {
			if r := recover(); r != nil {
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
				ch <- result{nil, err}
			}
		}()

		eventID, err := sender.Send()
		ch <- result{eventID, err}

	}()

	var res result

	select {
	case res = <-ch:
		switch {
		case res.err != nil:
			return nil, res.err
		case res.eventID.IsEmpty():
			return nil, errSending
		}
	case <-ctx.Done():
		return nil, errTimeout

	}

	return res.eventID, nil
}
