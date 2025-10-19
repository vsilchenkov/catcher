package sentryhub

import (
	"encoding/json"

	"github.com/cockroachdb/errors"
)

const typeSession = "session"

// SentrySession — тело самой сессии
// https://develop.sentry.dev/sdk/telemetry/sessions/
type SentrySession struct {
	Sid     string       `json:"sid"`
	Did     string       `json:"did"`
	Init    bool         `json:"init"`
	Started string       `json:"started"`
	Status  string       `json:"status"`
	Errors  int          `json:"errors"`
	Attrs   SessionAttrs `json:"attrs"`
}

// SessionAttrs — вложенный релиз/окружение
type SessionAttrs struct {
	Release     string `json:"release"`
	Environment string `json:"environment,omitempty"`
}

func (h Hub) StartSession(session SentrySession) error {

	const op = "sentryhub.StartSession"
	return h.starEndSession(session, op)

}

func (h Hub) EndSession(session SentrySession) error {

	const op = "sentryhub.EndSession"
	return h.starEndSession(session, op)

}

func (h Hub) starEndSession(session SentrySession, op string) error {

	payload, err := json.Marshal(session)
	if err != nil {
		return errors.WithMessage(err, op)
	}

	if err := h.SendRequest(typeSession, payload); err != nil {
		return errors.WithMessage(err, op)
	}

	return nil
}
