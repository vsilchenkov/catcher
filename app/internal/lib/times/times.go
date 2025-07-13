package times

import (
	"encoding/json"
	"fmt"
	"time"
)

const format = "2006-01-02T15:04:05"

// Кастомный тип для даты с поддержкой нескольких форматов
type TimeTZ struct {
	time.Time
}

func (tz *TimeTZ) UnmarshalJSON(b []byte) error {
	str := string(b)
	if len(str) < 2 {
		return fmt.Errorf("invalid date: %s", str)
	}
	str = str[1 : len(str)-1] // убираем кавычки

	layouts := []string{
		time.RFC3339,
		format,
	}

	var err error
	for _, layout := range layouts {
		var t time.Time
		t, err = time.Parse(layout, str)
		if err == nil {
			tz.Time = t
			return nil
		}
	}
	return fmt.Errorf("cannot parse date: %s", str)
}

type UnixTime time.Time

func (t *UnixTime) UnmarshalJSON(b []byte) error {
	var ts int64
	if err := json.Unmarshal(b, &ts); err != nil {
		return err
	}
	*t = UnixTime(time.Unix(ts, 0))
	return nil
}
