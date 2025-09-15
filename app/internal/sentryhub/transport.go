package sentryhub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/getsentry/sentry-go"
)

// EnvelopeHeader — заголовок "envelope"
type EnvelopeHeader struct {
    SentAt string `json:"sent_at"`
}

// ItemHeader — заголовок "item"
type ItemHeader struct {
    Type   string `json:"type"`
    Length int    `json:"length"`
}

func (h Hub) SendRequest(itemType string, payload []byte) error {

	// 1. Парсим DSN
	dsn, err := sentry.NewDsn(h.prj.Sentry.Dsn)
	if err != nil {
		return errors.WithMessagef(err, "invalid DSN: %s", dsn)
	}

	endpoint := dsn.GetAPIURL().String()
	publicKey := dsn.GetPublicKey()
	secretKey := dsn.GetSecretKey()

	// 2. Формируем envelope header (struct → json)
	envelopeHeader, err := getEnvelopeHeader()
	if err != nil {
		return errors.WithMessage(err, "invalid envelope header")
	}

	// 3. Формируем item header (struct → json)
	itemHeader, err := getitemHeader(itemType, payload)
	if err != nil {
		return errors.WithMessage(err, "invalid item header")
	}

	envelope := []byte(envelopeHeader + itemHeader + string(payload) + "\n")

	// 4. Сборка auth header
	authHeader := fmt.Sprintf(
		"Sentry sentry_key=%s,sentry_version=7,sentry_client=gosession/1.0",
		publicKey,
	)
	if secretKey != "" {
		authHeader += fmt.Sprintf(",sentry_secret=%s", secretKey)
	}

	// 5. HTTP-запрос
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(envelope))
	if err != nil {
		return err
	}
	req.Header.Set("X-Sentry-Auth", authHeader)
	req.Header.Set("Content-Type", "application/x-sentry-envelope")

	ctx, cancel := context.WithTimeout(req.Context(), time.Duration(h.prj.Sentry.IimeOut)*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("sentry respond with status: %d", resp.StatusCode)
	}
	return nil
}

func getEnvelopeHeader() (string, error) {

	now := time.Now().UTC()

	envelopeHeader := EnvelopeHeader{
		SentAt: now.Format(time.RFC3339Nano),
	}
	b, err := json.Marshal(envelopeHeader)
	if err != nil {
		return "", err
	}
	return string(b) + "\n", nil
}

func getitemHeader(itemType string, payload []byte) (string, error) {

	itemHeader := ItemHeader{
		Type:   itemType,
		Length: len(payload),
	}
	b, err := json.Marshal(itemHeader)
	if err != nil {
		return "", err
	}
	return string(b) + "\n", nil

}
