// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package logger

// This log writer sends output to a file
type FileLogWriter struct {
	rec  chan *LogRecord
	rot  chan bool
	file *RotateFile

	// The logging format
	format string
}

// This is the FileLogWriter's output method
func (w *FileLogWriter) LogWrite(rec *LogRecord) {
	defer func() {
		recover()
	}()
	w.rec <- rec
}

func (w *FileLogWriter) Close() {
	close(w.rec)
}

// NewFileLogWriter creates a new LogWriter which writes to the given file and
// has rotation enabled if rotate is true.
//
// If rotate is true, any time a new log file is opened, the old one is renamed
// with a .### extension to preserve it.  The various Set* methods can be used
// to configure log rotation based on lines, size, and daily.
//
// The standard log-line format is:
//   [%D %T] [%L] (%S) %M
func NewFileLogWriter(fnc RotateFilenameCreator, bufLen int, rotate bool) *FileLogWriter {
	wfile := NewRotateFile(fnc, rotate)
	if wfile == nil {
		return nil
	}
	w := &FileLogWriter{
		rec:    make(chan *LogRecord, bufLen),
		rot:    make(chan bool),
		file:   wfile,
		format: "%D %T [%S] %L %M",
	}

	// open the file for the first time
	go func() {
		defer func() {
			if w.file != nil {
				w.file.Close()
			}
		}()

		for {
			select {
			case <-w.rot:
				w.file.Rotate()
			case rec, ok := <-w.rec:
				if !ok {
					return
				}
				msg := FormatLogRecord(w.format, rec)
				if !w.file.Write(rec.Created.Day(), msg) {
					return
				}
			}
		}
	}()

	return w
}

// Request that the logs rotate
func (w *FileLogWriter) Rotate() {
	w.rot <- true
}

// Set the logging format (chainable).  Must be called before the first log
// message is written.
func (w *FileLogWriter) SetFormat(format string) *FileLogWriter {
	w.format = format
	return w
}

// Set rotate at linecount (chainable). Must be called before the first log
// message is written.
func (w *FileLogWriter) SetRotateLines(maxlines int) *FileLogWriter {
	if w.file != nil {
		w.file.SetRotateLines(maxlines)
	}
	return w
}

// Set rotate at size (chainable). Must be called before the first log message
// is written.
func (w *FileLogWriter) SetRotateSize(maxsize int) *FileLogWriter {
	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotateSize: %v\n", maxsize)
	if w.file != nil {
		w.file.SetRotateSize(maxsize)
	}
	return w
}

// Set rotate daily (chainable). Must be called before the first log message is
// written.
func (w *FileLogWriter) SetRotateDaily(daily bool) *FileLogWriter {
	if w.file != nil {
		w.file.SetRotateDaily(daily)
	}
	return w
}

// SetRotate changes whether or not the old logs are kept. (chainable) Must be
// called before the first log message is written.  If rotate is false, the
// files are overwritten; otherwise, they are rotated to another file before the
// new log is opened.
func (w *FileLogWriter) SetRotate(rotate bool) *FileLogWriter {
	if w.file != nil {
		w.file.SetRotate(rotate)
	}
	return w
}
