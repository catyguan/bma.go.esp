package logger

import (
	"bytes"
	"config"
	"errors"
	"fmt"
	"strings"
	"time"
)

/****** Constants ******/

// These are the integer logging levels used by the logger
type level int

const (
	LEVEL_ALL level = iota
	LEVEL_DEBUG
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_NONE
	NOUSE
)

const (
	timeFormatString = "2006-01-02 15:04:05"
	logBufferLength  = 1024
)

// Logging level strings
var (
	levelStrings      = [...]string{"ALL", "DEBUG", "INFO", "WARN", "ERROR", "NONE", "NOUSE"}
	levelPrintStrings = [...]string{" ALL ", "DEBUG", "INFO ", "WARN ", "ERROR", "NONE ", "NOUSE"}
)

func (l level) String() string {
	if l < 0 || int(l) > len(levelStrings) {
		return "UNKNOWN"
	}
	return levelStrings[int(l)]
}

func (l level) printString() string {
	if l < 0 || int(l) > len(levelPrintStrings) {
		return "UNKNOWN"
	}
	return levelPrintStrings[int(l)]
}

func StringToLevel(s string) int {
	switch s {
	case "ALL":
		return int(LEVEL_ALL)
	case "DEBUG":
		return int(LEVEL_DEBUG)
	case "INFO":
		return int(LEVEL_INFO)
	case "WARN":
		return int(LEVEL_WARN)
	case "ERROR":
		return int(LEVEL_ERROR)
	case "NONE":
		return int(LEVEL_NONE)
	}
	return int(NOUSE)
}

/****** LogRecord ******/

// A LogRecord contains all of the pertinent information for each message
type LogRecord struct {
	Level   level     // The log level
	Created time.Time // The time at which the log message was created (nanoseconds)
	Tag     string    // The message source
	Message string    // The log message
}

/****** LogWriter ******/

// This is an interface for anything that should be able to write logs
type LogWriter interface {
	// This will be called to log a LogRecord message.
	LogWrite(rec *LogRecord)

	// This should clean up anything lingering about the LogWriter, as it is called before
	// the LogWriter is removed.  LogWrite should not be called after Close.
	Close()
}

/****** Logger ******/

// A Filter represents the log level below which no log records are written to
// the associated LogWriter.
type Filter struct {
	level  level
	writer LogWriter
}

// LoggerConfig
type LoggerConfig struct {
	root    Filter
	filters map[string]Filter
	writers map[string]LogWriter
}

var (
	logConfig LoggerConfig = LoggerConfig{Filter{LEVEL_ALL, nil}, make(map[string]Filter), make(map[string]LogWriter)}
)

// Create a new logger.
//
// DEPRECATED: Use make(Logger) instead.
func Config() *LoggerConfig {
	return &logConfig
}

func (cfg *LoggerConfig) GetLevel(tag string) level {
	if tag != "" {
		fs := logConfig.filters
		if f, ok := fs[tag]; ok {
			return f.level
		}
		return NOUSE
	}
	return logConfig.root.level
}

func newFilter(tag string, l level, w LogWriter) {
	obj := make(map[string]Filter)
	for k, v := range logConfig.filters {
		obj[k] = v
	}
	obj[tag] = Filter{l, w}
	logConfig.filters = obj
}

func newWriter(name string, w LogWriter) {
	old := logConfig.writers[name]
	if old != nil {
		panic("LogWriter '" + name + "' already exists")
	}
	obj := make(map[string]LogWriter)
	for k, v := range logConfig.writers {
		obj[k] = v
	}
	obj[name] = w
	logConfig.writers = obj
}

func (cfg *LoggerConfig) SetLevel(tag string, l level) level {
	if tag != "" {
		fs := logConfig.filters
		if f, ok := fs[tag]; ok {
			old := f.level
			f.level = l
			return old
		}
		newFilter(tag, l, nil)
		return NOUSE
	} else {
		old := logConfig.root.level
		logConfig.root.level = l
		return old
	}
}

func (cfg *LoggerConfig) HasWriter(name string) bool {
	_, ok := logConfig.writers[name]
	return ok
}

func (cfg *LoggerConfig) SetWriter(tag string, writerName string) bool {
	w, ok := logConfig.writers[writerName]
	if !ok {
		return false
	}
	if tag != "" {
		fs := logConfig.filters
		if f, ok := fs[tag]; ok {
			f.writer = w
			return true
		}
		newFilter(tag, NOUSE, w)
		return true
	} else {
		logConfig.root.writer = w
		return true
	}
}

func (cfg *LoggerConfig) NewWriter(name string, w LogWriter) bool {
	_, ok := logConfig.writers[name]
	if !ok {
		return false
	}
	newWriter(name, w)
	return true
}

func createWriter(wobj beanWriter) LogWriter {
	typ := wobj.Type
	if typ == "" {
		fmt.Println("logger writer type is nil")
		return nil
	}
	enable := !wobj.Disable
	if !enable {
		fmt.Printf("logger writer '%s' disable\n", wobj.Name)
		return nil
	}
	if typ == "console" {
		buflen := wobj.BufferSize
		if buflen == 0 {
			buflen = logBufferLength
		}
		return NewConsoleLogWriter(buflen)
	}
	if typ == "file" {
		if wobj.File == "" {
			fmt.Printf("ERROR: fileLogWriter '%s' file invalid\n", wobj.Name)
			return nil
		}
		buflen := wobj.BufferSize
		if buflen == 0 {
			buflen = logBufferLength
		}
		fnc := func(tm time.Time, num int) string {
			out := bytes.NewBuffer(make([]byte, 0))
			out.WriteString(wobj.File)
			out.WriteString("_")
			out.WriteString(tm.Format("20060102"))
			if num != 0 {
				out.WriteString(".")
				out.WriteString(fmt.Sprintf("%d", num))
			}
			out.WriteString(".log")
			return out.String()
		}
		w := NewFileLogWriter(fnc, buflen, !wobj.NoRotate)
		if wobj.Format != "" {
			w.SetFormat(wobj.Format)
		}
		w.SetRotateDaily(!wobj.NoDaily)
		if wobj.Maxlines > 0 {
			w.SetRotateLines(wobj.Maxlines)
		}
		if wobj.Maxsize > 0 {
			w.SetRotateSize(wobj.Maxsize)
		}
		return w
	}
	if typ == "link" {
		w1 := logConfig.writers[wobj.Writer1]
		w2 := logConfig.writers[wobj.Writer2]
		if w1 == nil && w2 == nil {
			fmt.Printf("ERROR: LinkLogWriter '%s' all writer invalid\n", wobj.Name)
			return nil
		}
		return NewLinkLogWriter(w1, w2)
	}
	fmt.Printf("logger unknow writer '%s' type '%s'\n", wobj.Name, typ)
	return nil
}

func createFilter(fobj beanFilter) *Filter {
	enable := !fobj.Disable
	if !enable {
		fmt.Printf("logger filter '%s' disable\n", fobj.Name)
		return nil
	}
	var w LogWriter
	wname := fobj.Writer
	if wname == "" {
		w = nil
	} else {
		w = logConfig.writers[wname]
		if w == nil {
			fmt.Printf("logger filter '%s' invalid writer '%s'\n", fobj.Name, wname)
			return nil
		}
	}
	lstr := fobj.Level
	lvl := StringToLevel(lstr)

	if w == nil {
		wname = "<none>"
	}
	fmt.Printf("logger filter '%s' => %s '%s'\n", fobj.Name, level(lvl).String(), wname)
	return &Filter{level(lvl), w}
}

type beanWriter struct {
	Name       string
	Type       string
	Disable    bool
	BufferSize int
	Format     string
	// rotate
	File     string
	NoRotate bool
	Maxlines int
	Maxsize  int
	NoDaily  bool
	// link
	Writer1 string
	Writer2 string
}

type beanFilter struct {
	Name    string
	Level   string
	Writer  string
	Disable bool
}

type beanLogger struct {
	RootLevel string
	Writer    []beanWriter
	Filter    []beanFilter
}

func (lc *LoggerConfig) InitLogger() {
	var beanLogger beanLogger
	if config.Global.GetBeanConfig("logger", &beanLogger) {
		if beanLogger.RootLevel != "" {
			l := level(StringToLevel(beanLogger.RootLevel))
			fmt.Printf("logger root level = %s\n", l)
			logConfig.root.level = l
		}
		wlist := beanLogger.Writer
		if wlist != nil {
			for _, wobj := range wlist {
				if wobj.Name == "" {
					fmt.Println("writer no Name")
					continue
				}
				w := createWriter(wobj)
				if w != nil {
					// fmt.Println("new writer =>", w)
					logConfig.writers[wobj.Name] = w
				}
			}
		}
		flist := beanLogger.Filter
		if flist != nil {
			for _, fobj := range flist {
				if fobj.Name == "" {
					fmt.Println("filter no Name")
					continue
				}
				f := createFilter(fobj)
				if f != nil {
					// fmt.Println("new filter =>", f)
					if fobj.Name == "root" {
						if f.level != NOUSE {
							logConfig.root.level = f.level
						}
						if f.writer != nil {
							logConfig.root.writer = f.writer
						}
					} else {
						logConfig.filters[fobj.Name] = *f
					}
				}
			}
		}
	}
	// default init
	if logConfig.root.writer == nil {
		initDefaultLogger()
	}
}

var (
	closed bool
)

func Close() {
	closed = true
	time.Sleep(1 * time.Millisecond)

	// Close all open loggers
	logConfig.filters = make(map[string]Filter)
	logConfig.root.writer = nil

	ws := logConfig.writers
	logConfig.writers = make(map[string]LogWriter)
	for _, w := range ws {
		w.Close()
	}
}

func initDefaultLogger() {
	dw := NewConsoleLogWriter(logBufferLength)
	logConfig.writers["console"] = dw

	fmt.Println("logger root default writer: console")
	logConfig.root.writer = dw
}

/******* Logging *******/
func getFilter(tag string, lvl level) (LogWriter, bool) {
	var l level = NOUSE
	var w LogWriter = nil
	fs := logConfig.filters
	if f, ok := fs[tag]; ok {
		l = f.level
		w = f.writer
	}

	if l == NOUSE {
		l = logConfig.root.level
	}
	if w == nil {
		w = logConfig.root.writer
	}

	// Determine if any logging will be done
	if lvl < l {
		return nil, false
	}
	if w == nil && !closed {
		initDefaultLogger()
		w = logConfig.root.writer
	}
	return w, true
}

// Send a formatted log message internally
func intLogf(tag string, lvl level, format string, args ...interface{}) {

	w, ok := getFilter(tag, lvl)

	if !ok {
		return
	}
	if w == nil {
		return
	}

	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}

	// Make the log record
	rec := &LogRecord{
		Level:   lvl,
		Created: time.Now(),
		Tag:     tag,
		Message: msg,
	}

	// Dispatch the logs
	w.LogWrite(rec)
}

// Send a closure log message internally
func intLogc(tag string, lvl level, closure func() string) {

	w, ok := getFilter(tag, lvl)

	if !ok {
		return
	}

	// Make the log record
	rec := &LogRecord{
		Level:   lvl,
		Created: time.Now(),
		Tag:     tag,
		Message: closure(),
	}

	// Dispatch the logs
	w.LogWrite(rec)
}

// Send a log message with manual level, source, and message.
func Log(tag string, lvl level, message string) {

	w, ok := getFilter(tag, lvl)

	if !ok {
		return
	}

	// Make the log record
	rec := &LogRecord{
		Level:   lvl,
		Created: time.Now(),
		Tag:     tag,
		Message: message,
	}

	// Dispatch the logs
	w.LogWrite(rec)
}

// Logf logs a formatted log message at the given log level, using the caller as
// its source.
func Logf(tag string, lvl level, format string, args ...interface{}) {
	intLogf(tag, lvl, format, args...)
}

// Logc logs a string returned by the closure at the given log level, using the caller as
// its source.  If no log message would be written, the closure is never called.
func Logc(tag string, lvl level, closure func() string) {
	intLogc(tag, lvl, closure)
}

func Enable(tag string, lvl level) bool {
	_, ok := getFilter(tag, lvl)
	return ok
}

func EnableDebug(tag string) bool {
	return Enable(tag, LEVEL_DEBUG)
}

func Debug(tag string, arg0 interface{}, args ...interface{}) {
	const (
		lvl = LEVEL_DEBUG
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		intLogf(tag, lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		intLogc(tag, lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		intLogf(tag, lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

func EnableInfo(tag string) bool {
	return Enable(tag, LEVEL_INFO)
}

// Info logs a message at the info log level.
// See Debug for an explanation of the arguments.
func Info(tag string, arg0 interface{}, args ...interface{}) {
	const (
		lvl = LEVEL_INFO
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		intLogf(tag, lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		intLogc(tag, lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		intLogf(tag, lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

func EnableWarn(tag string) bool {
	return Enable(tag, LEVEL_WARN)
}

// Warn logs a message at the warning log level and returns the formatted error.
// At the warning level and higher, there is no performance benefit if the
// message is not actually logged, because all formats are processed and all
// closures are executed to format the error message.
// See Debug for further explanation of the arguments.
func Warn(tag string, arg0 interface{}, args ...interface{}) error {
	const (
		lvl = LEVEL_WARN
	)
	var msg string
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		msg = fmt.Sprintf(first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		msg = first()
	default:
		// Build a format string so that it will be similar to Sprint
		msg = fmt.Sprintf(fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
	}
	intLogf(tag, lvl, msg)
	return errors.New(msg)
}

func EnableError(tag string) bool {
	return Enable(tag, LEVEL_ERROR)
}

// Error logs a message at the error log level and returns the formatted error,
// See Warn for an explanation of the performance and Debug for an explanation
// of the parameters.
func Error(tag string, arg0 interface{}, args ...interface{}) error {
	const (
		lvl = LEVEL_ERROR
	)
	var msg string
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		msg = fmt.Sprintf(first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		msg = first()
	default:
		// Build a format string so that it will be similar to Sprint
		msg = fmt.Sprintf(fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
	}
	intLogf(tag, lvl, msg)
	return errors.New(msg)
}

type formatCacheType struct {
	LastUpdateSeconds    int64
	shortTime, shortDate string
	longTime, longDate   string
}

var formatCache = &formatCacheType{}

// Known format codes:
// %T - Time (15:04:05 MST)
// %t - Time (15:04)
// %D - Date (2006/01/02)
// %d - Date (01/02/06)
// %L - Level (FNST, FINE, DEBG, TRAC, WARN, EROR, CRIT)
// %S - Source
// %M - Message
// Ignores unknown formats
// Recommended: "[%D %T] [%L] (%S) %M"
func FormatLogRecord(format string, rec *LogRecord) string {
	if rec == nil {
		return "<nil>"
	}
	if len(format) == 0 {
		return ""
	}

	out := bytes.NewBuffer(make([]byte, 0))
	secs := rec.Created.UnixNano() / 1e9

	cache := *formatCache
	if cache.LastUpdateSeconds != secs {
		month, day, year := rec.Created.Month(), rec.Created.Day(), rec.Created.Year()
		hour, minute, second := rec.Created.Hour(), rec.Created.Minute(), rec.Created.Second()
		updated := &formatCacheType{
			LastUpdateSeconds: secs,
			shortTime:         fmt.Sprintf("%02d:%02d", hour, minute),
			shortDate:         fmt.Sprintf("%02d-%02d-%02d", year%100, month, day),
			longTime:          fmt.Sprintf("%02d:%02d:%02d", hour, minute, second),
			longDate:          fmt.Sprintf("%04d-%02d-%02d", year, month, day),
		}
		cache = *updated
		formatCache = updated
	}

	// Split the string into pieces by % signs
	pieces := bytes.Split([]byte(format), []byte{'%'})

	// Iterate over the pieces, replacing known formats
	for i, piece := range pieces {
		if i > 0 && len(piece) > 0 {
			switch piece[0] {
			case 'T':
				out.WriteString(cache.longTime)
			case 't':
				out.WriteString(cache.shortTime)
			case 'D':
				out.WriteString(cache.longDate)
			case 'd':
				out.WriteString(cache.shortDate)
			case 'L':
				out.WriteString(levelStrings[rec.Level])
			case 'S':
				out.WriteString(rec.Tag)
			case 'M':
				out.WriteString(rec.Message)
			}
			if len(piece) > 1 {
				out.Write(piece[1:])
			}
		} else if len(piece) > 0 {
			out.Write(piece)
		}
	}
	out.WriteByte('\n')

	return out.String()
}

func Sprintf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
