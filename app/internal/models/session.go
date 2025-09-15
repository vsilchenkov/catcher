package models

import (
	"catcher/app/internal/config"
	"fmt"
	"time"
)

type Session struct {
	Sid            string    `json:"sid" binding:"required"`
	Did            string    `json:"did" binding:"required"`
	Started        time.Time `json:"started"`
	Release        string    `json:"release" binding:"required"`
	Environment    string    `json:"environment,omitempty" binding:"required"`
	ErrorsEventer  int
	ErrorsReporter int
}

func (s Session) Key(prj config.Project) string {
	op := "session"
	return fmt.Sprintf("%s:%s:%s:%s:%s", prj.Id, op, s.Environment, s.Did, s.Sid)
}
