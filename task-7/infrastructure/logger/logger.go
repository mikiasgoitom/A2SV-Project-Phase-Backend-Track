package logger

import (
	"log"
	"os"
)

var env string

type StdLogger struct{}

func init() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
}

func NewLogger() *StdLogger {
	return &StdLogger{}
}

func (l *StdLogger) Info(msg string) {
	log.Printf("[INFO] %s\n", msg)
}

func (l *StdLogger) Warn(msg string) {
	log.Printf("[WARN] %s\n", msg)
}

func (l *StdLogger) Debug(msg string) {
	if env == "development" {
		log.Printf("[DEBUG] %s\n", msg)
	}
}

func (l *StdLogger) Error(msg string) {
	log.Printf("[ERROR] %s\n", msg)
}

func (l *StdLogger) Secure(tag string, msg string) {
	if env == "development" {
		log.Printf("[SECURE][%s] %s", tag, msg)
	} else {
		log.Printf("[SECURE][%s] ***REDACTED***", tag)
	}
}
