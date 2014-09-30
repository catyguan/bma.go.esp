package logger

import (
	"fmt"
	"os"
	"time"
)

type RotateFilenameCreator func(tm time.Time, num int) string

type RotateFile struct {
	// The opened file
	filename        string
	file            *os.File
	filenameCreator RotateFilenameCreator

	// Rotate at linecount
	maxlines          int
	maxlines_curlines int

	// Rotate at size
	maxsize         int
	maxsize_cursize int

	// Rotate daily
	daily          bool
	daily_opendate int

	// Keep old logfiles (.001, .002, etc)
	rotate bool
}

func NewRotateFile(fnc RotateFilenameCreator, rotate bool) *RotateFile {
	w := &RotateFile{
		filenameCreator: fnc,
		rotate:          rotate,
	}

	// open the file for the first time
	if err := w.initRotate(); err != nil {
		return nil
	}
	return w
}

func (w *RotateFile) Close() {
	if w.file != nil {
		w.file.Close()
		w.file = nil
	}
}

func (w *RotateFile) Write(day int, msg string) bool {
	now := time.Now()
	if day == 0 {
		day = now.Day()
	}
	if (w.maxlines > 0 && w.maxlines_curlines >= w.maxlines) ||
		(w.maxsize > 0 && w.maxsize_cursize >= w.maxsize) ||
		(w.daily && day != w.daily_opendate) {
		if err := w.initRotate(); err != nil {
			fmt.Printf("ERROR: RotateFile(%q): %s\n", w.filename, err)
			return false
		}
	}

	// Perform the write
	n, err := fmt.Fprint(w.file, msg)
	if err != nil {
		fmt.Printf("ERROR: RotateFile(%q): %s\n", w.filename, err)
		return false
	}

	// Update the counts
	w.maxlines_curlines++
	w.maxsize_cursize += n
	return true
}

// Request that the logs rotate
func (w *RotateFile) Rotate() bool {
	if err := w.initRotate(); err != nil {
		fmt.Printf("ERROR: FileLogWriter(%q): %s\n", w.filename, err)
		return false
	}
	return true
}

// If this is called in a threaded context, it MUST be synchronized
func (w *RotateFile) initRotate() error {
	// Close any log file that may be open
	w.Close()

	now := time.Now()
	// If we are keeping log files, move it to the next available number
	fn := w.filenameCreator(now, 0)
	if w.rotate && fn == w.filename {
		_, err := os.Lstat(w.filename)
		if err == nil { // file exists
			// Find the next available number
			num := 1
			fname := ""
			for ; err == nil && num <= 999; num++ {
				fname = w.filenameCreator(now, num)
				_, err = os.Lstat(fname)
			}
			// return error if the last file checked still existed
			if err == nil {
				return fmt.Errorf("ERROR: RotateFile Cannot find free log number to rename %s\n", w.filename)
			}

			// Rename the file to its newfound home
			err = os.Rename(w.filename, fname)
			if err != nil {
				return fmt.Errorf("ERROR: RotateFile rename fail %s\n", err)
			}
		}
	}

	// Set the daily open date to the current date
	w.daily_opendate = now.Day()
	w.filename = fn

	// Open the log file
	fd, err := os.OpenFile(w.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		return err
	}
	w.file = fd

	// initialize rotation values
	w.maxlines_curlines = 0
	w.maxsize_cursize = 0

	return nil
}

// Set rotate at linecount (chainable). Must be called before the first log
// message is written.
func (w *RotateFile) SetRotateLines(maxlines int) *RotateFile {
	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotateLines: %v\n", maxlines)
	w.maxlines = maxlines
	return w
}

// Set rotate at size (chainable). Must be called before the first log message
// is written.
func (w *RotateFile) SetRotateSize(maxsize int) *RotateFile {
	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotateSize: %v\n", maxsize)
	w.maxsize = maxsize
	return w
}

// Set rotate daily (chainable). Must be called before the first log message is
// written.
func (w *RotateFile) SetRotateDaily(daily bool) *RotateFile {
	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotateDaily: %v\n", daily)
	w.daily = daily
	return w
}

// SetRotate changes whether or not the old logs are kept. (chainable) Must be
// called before the first log message is written.  If rotate is false, the
// files are overwritten; otherwise, they are rotated to another file before the
// new log is opened.
func (w *RotateFile) SetRotate(rotate bool) *RotateFile {
	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotate: %v\n", rotate)
	w.rotate = rotate
	return w
}
