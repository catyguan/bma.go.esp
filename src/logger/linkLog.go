package logger

import ()

type LinkLogWriter struct {
	writer1 LogWriter
	writer2 LogWriter
}

func NewLinkLogWriter(w1 LogWriter, w2 LogWriter) *LinkLogWriter {
	r := new(LinkLogWriter)
	r.writer1 = w1
	r.writer2 = w2
	return r
}

func (w LinkLogWriter) LogWrite(rec *LogRecord) {
	w.writer1.LogWrite(rec)
	w.writer2.LogWrite(rec)
}

func (w LinkLogWriter) Close() {

}
