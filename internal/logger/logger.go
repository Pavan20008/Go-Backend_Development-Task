package logger

import (
	"go.uber.org/zap"
)

// New builds a configured Uber Zap logger. A production JSON logger is used
// outside of development, and a human-friendly console logger otherwise.
func New(environment string) (*zap.Logger, error) {
	if environment == "production" {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}
