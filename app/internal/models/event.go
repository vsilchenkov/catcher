package models

import (
	"catcher/app/internal/lib/times"

	"github.com/getsentry/sentry-go"
)

const OpEventer = "eventer"

type EventID string

type Event struct {
	sentry.Event
	Timestamp times.UnixTime `json:"timestamp"`
	Exception Exception      `json:"exception"`
	Request   any            `json:"request"`
}

type Exception struct {
	Values []ExceptionValue `json:"values"`
}

type ExceptionValue struct {
	Type       string     `json:"type"`
	Value      string     `json:"value"`
	Stacktrace Stacktrace `json:"stacktrace"`
}

type Stacktrace struct {
	Frames []Frame `json:"frames"`
}

type Frame struct {
	Lineno      int    `json:"lineno"`
	Function    string `json:"function"`
	Filename    string `json:"filename"`
	Module      string `json:"module"`
	ModuleAbs   string `json:"module_abs"`
	ContextLine string `json:"context_line"`
	InApp       bool   `json:"in_app"`
	AbsPath     string `json:"abs_path"`
	StackStart  bool   `json:"stack_start"`
}

type SendEventResult struct {
	ID      string
	EventID *EventID
}

func (e *EventID) IsEmpty() bool {
	return e == nil || *e == ""
}

func (e EventID) String() string {
	return string(e)
}
