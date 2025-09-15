package errors

import (
	"fmt"

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
