package acclog

import "logger"

type AccessLoggerFile struct {
	rec  chan AccLogInfo
	rot  chan bool
	file *logger.RotateFile
	cfg  map[string]string
}

func (w *AccessLoggerFile) Write(rec AccLogInfo) {
	defer func() {
		recover()
	}()
	w.rec <- rec
}

func (w *AccessLoggerFile) Close() {
	close(w.rec)
}

func NewFile(fnc logger.RotateFilenameCreator, qlen int, rotate bool, cfg map[string]string) *AccessLoggerFile {
	wfile := logger.NewRotateFile(fnc, rotate)
	if wfile == nil {
		return nil
	}
	w := &AccessLoggerFile{
		rec:  make(chan AccLogInfo, qlen),
		rot:  make(chan bool),
		file: wfile,
		cfg:  cfg,
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
				msg := rec.Message(w.cfg)
				w.file.Write(rec.TimeDay(), msg)
			}
		}
	}()

	return w
}

func (w *AccessLoggerFile) Rotate() {
	w.rot <- true
}

// Set rotate at linecount (chainable). Must be called before the first log
// message is written.
func (w *AccessLoggerFile) SetRotateLines(maxlines int) *AccessLoggerFile {
	if w.file != nil {
		w.file.SetRotateLines(maxlines)
	}
	return w
}

// Set rotate at size (chainable). Must be called before the first log message
// is written.
func (w *AccessLoggerFile) SetRotateSize(maxsize int) *AccessLoggerFile {
	//fmt.Fprintf(os.Stderr, "AccessLoggerFile.SetRotateSize: %v\n", maxsize)
	if w.file != nil {
		w.file.SetRotateSize(maxsize)
	}
	return w
}

// Set rotate daily (chainable). Must be called before the first log message is
// written.
func (w *AccessLoggerFile) SetRotateDaily(daily bool) *AccessLoggerFile {
	if w.file != nil {
		w.file.SetRotateDaily(daily)
	}
	return w
}

// SetRotate changes whether or not the old logs are kept. (chainable) Must be
// called before the first log message is written.  If rotate is false, the
// files are overwritten; otherwise, they are rotated to another file before the
// new log is opened.
func (w *AccessLoggerFile) SetRotate(rotate bool) *AccessLoggerFile {
	if w.file != nil {
		w.file.SetRotate(rotate)
	}
	return w
}
