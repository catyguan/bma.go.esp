package logger

import (
	"fmt"
	"io"
	"os"
)

var stdout io.Writer = os.Stdout

// This is the standard writer that prints to standard output.
type ConsoleLogWriter chan *LogRecord

// This creates a new ConsoleLogWriter
func NewConsoleLogWriter(logBufferLength int) ConsoleLogWriter {
	records := make(ConsoleLogWriter, logBufferLength)
	go records.run(stdout)
	return records
}

func (w ConsoleLogWriter) String() string {
	return "ConsoleLogWriter"
}

func (w ConsoleLogWriter) run(out io.Writer) {
	var timestr string
	var timestrAt int64

	for rec := range w {
		if at := rec.Created.UnixNano() / 1e9; at != timestrAt {
			timestr, timestrAt = rec.Created.Format(timeFormatString), at
		}
		fmt.Fprint(out, timestr, " [", rec.Tag, "] ", levelPrintStrings[rec.Level], " ", rec.Message, "\n")
	}
}

// This is the ConsoleLogWriter's output method.  This will block if the output
// buffer is full.
func (w ConsoleLogWriter) LogWrite(rec *LogRecord) {
	defer func() {
		recover()
	}()
	w <- rec
}

// Close stops the logger from sending messages to standard output.  Attempts to
// send log messages to this logger after a Close have undefined behavior.
func (w ConsoleLogWriter) Close() {
	close(w)
}
