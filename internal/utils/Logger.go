package utils

import (
	"log"
	"os"
	"time"
)

type AuditIntegrator interface {
	SaveAudit(event AuditEvent) error
}

type AuditEvent struct {
	Action string
	Entity string
	ID     string
	Time   time.Time
}

type AuditLogger struct {
	ch         chan AuditEvent
	logger     *log.Logger
	integrator AuditIntegrator
}

func NewAuditLogger(buffer int, integrator AuditIntegrator) *AuditLogger {
	return &AuditLogger{
		ch: make(chan AuditEvent, buffer),
		logger: log.New(
			os.Stdout,
			"[AUDIT] ",
			log.LstdFlags|log.LUTC,
		),
		integrator: integrator,
	}
}

func (l *AuditLogger) Start() {
	go func() {
		for event := range l.ch {
			l.logger.Printf(
				"action=%s entity=%s id=%s time=%s",
				event.Action,
				event.Entity,
				event.ID,
				event.Time.Format(time.RFC3339),
			)
			if l.integrator != nil {
				_ = l.integrator.SaveAudit(event)
			}
		}
	}()
}

func (l *AuditLogger) Log(event AuditEvent) {
	select {
	case l.ch <- event:
	default:

	}
}
