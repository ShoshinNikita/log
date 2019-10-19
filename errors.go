package clog

import (
	"fmt"
	"time"
)

// Error prints error
// Output pattern: (?time) [ERR] (?file:line) error
func (l Logger) Error(v ...interface{}) {
	print := func() (int, error) {
		return fmt.Fprintln(l.buff, v...)
	}

	l.error(print)
}

// Errorf prints error
// Output pattern: (?time) [ERR] (?file:line) error
func (l Logger) Errorf(format string, v ...interface{}) {
	print := func() (int, error) {
		return fmt.Fprintf(l.buff, format, v...)
	}

	l.error(print)
}

// error is an internal function for printing error messages
// Output pattern: (?time) [ERR] (?file:line) error
func (l Logger) error(print messagePrintFunction) {
	if !l.shouldPrint(LevelError) {
		return
	}

	now := time.Now()

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.buff.Reset()

	l.writeIntoBuffer(l.getTime(now))
	l.writeIntoBuffer(l.getErrPrefix())
	l.writeIntoBuffer(l.getCaller())

	print()

	l.output.Write(l.buff.Bytes())
}
