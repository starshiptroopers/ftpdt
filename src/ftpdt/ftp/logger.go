package ftp

import (
	"fmt"
	"io"
	"log"
)

// StdLogger use an instance of this to log in a standard format
type DefaultFTPLogger struct{
	Logger *log.Logger
}

var defLogPrefix = "FTP"

func NewDefaultFTPLogger(out io.Writer) *DefaultFTPLogger {
	return &DefaultFTPLogger{ Logger: log.New(out, "", log.LstdFlags)}
}

func (l *DefaultFTPLogger) Print(sessionID string, message interface{}) {
	var s = ""
	if sessionID != "" {
		s = " "
	}

	l.Logger.Printf("%s %s%s%s", defLogPrefix, sessionID, s, message)
}

func (l *DefaultFTPLogger) Printf(sessionID string, format string, v ...interface{}) {
	l.Print(sessionID, fmt.Sprintf(format, v...))
}

func (l *DefaultFTPLogger) PrintCommand(sessionID string, command string, params string) {
	if command == "PASS" {
		l.Logger.Printf("%s %s > PASS ****",defLogPrefix, sessionID)
	} else {
		l.Logger.Printf("%s %s > %s %s",defLogPrefix, sessionID, command, params)
	}
}

func (l *DefaultFTPLogger) PrintResponse(sessionID string, code int, message string) {
	l.Logger.Printf("%s %s < %d %s",defLogPrefix, sessionID, code, message)
}
