package ftp

import (
	"fmt"
	"io"
	"log"
)

// DefaultFTPLogger is a logger used by goftp driver to log FTP activity
type DefaultFTPLogger struct{
	Logger *log.Logger
}

var defLogPrefix = "FTP"

//New NewDefaultFTPLogger create a default Logger instance which write logs to out io.Writer
func NewDefaultFTPLogger(out io.Writer) *DefaultFTPLogger {
	return &DefaultFTPLogger{ Logger: log.New(out, "", log.LstdFlags)}
}

//Printf write the message to the log
func (l *DefaultFTPLogger) Print(sessionID string, message interface{}) {
	var s = ""
	if sessionID != "" {
		s = " "
	}

	l.Logger.Printf("%s %s%s%s", defLogPrefix, sessionID, s, message)
}

//Printf write a formatted string to the log
func (l *DefaultFTPLogger) Printf(sessionID string, format string, v ...interface{}) {
	l.Print(sessionID, fmt.Sprintf(format, v...))
}

// PrintCommand log FTP control command
func (l *DefaultFTPLogger) PrintCommand(sessionID string, command string, params string) {
	if command == "PASS" {
		l.Logger.Printf("%s %s > PASS ****",defLogPrefix, sessionID)
	} else {
		l.Logger.Printf("%s %s > %s %s",defLogPrefix, sessionID, command, params)
	}
}

// PrintResponse log response to FTP control command
func (l *DefaultFTPLogger) PrintResponse(sessionID string, code int, message string) {
	l.Logger.Printf("%s %s < %d %s",defLogPrefix, sessionID, code, message)
}
