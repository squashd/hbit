package hbit

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func NewUUID() string {
	return uuid.New().String()
}

func NewEventId(domain string) string {
	return domain + "-" + NewUUID()
}

func NewEventIdWithTimestamp(domain string) string {
	timestamp := time.Now().UTC().Format("20060102-150405")
	return fmt.Sprintf("%s-%s-%s", domain, timestamp, NewUUID())
}
