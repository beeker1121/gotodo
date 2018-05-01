package creek

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Logger defines our custom Logger type.
type Logger struct {
	Filename string // the file to log to
	MaxSize  int64  // max file size in MB
	file     *os.File
	size     int64
	mu       sync.Mutex
}

// New creates a new creek logger.
func New(filename string, maxSize int64) *Logger {
	return &Logger{
		Filename: filename,
		MaxSize:  maxSize,
	}
}

// Write satisfies the io.Writer interface.
func (l *Logger) Write(p []byte) (n int, err error) {
	// Lock the mutex.
	l.mu.Lock()
	defer l.mu.Unlock()

	writeLen := int64(len(p))

	// If the data to write exceeds our max file size, error out.
	if writeLen > l.maxSize() {
		return 0, fmt.Errorf("Write length %d exceeds maximum file size %d", writeLen, l.maxSize())
	}

	// Get current log file.
	if l.file == nil {
		if err = l.openExistingOrNew(len(p)); err != nil {
			return 0, err
		}
	}

	// If writing the new data will go over our max file size, rotate the log file.
	if l.size+writeLen > l.maxSize() {
		if err := l.rotate(); err != nil {
			return 0, err
		}
	}

	// Write to the log file.
	n, err = l.file.Write(p)
	l.size += int64(n)

	return n, err
}

// close closes the log file if it's open.
func (l *Logger) close() error {
	if l.file == nil {
		return nil
	}

	err := l.file.Close()
	l.file = nil

	return err
}

// rotate rotates the log file.
func (l *Logger) rotate() error {
	// Close the current log file.
	if err := l.close(); err != nil {
		return err
	}

	// Open a new log file.
	if err := l.openNew(); err != nil {
		return err
	}

	return nil
}

// compressLogFile compresses a log file.
func compressLogFile(name string) {
	// Open the given log file for reading.
	file, err := os.Open(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open log file for compression: %s\n", err)
		return
	}
	defer file.Close()

	// Create a file to save the compressed data to.
	filegz, err := os.OpenFile(name+".gz", os.O_CREATE|os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create file for compression: %s\n", err)
		return
	}
	defer filegz.Close()

	// Create new gzip Writer.
	gz := gzip.NewWriter(filegz)
	defer gz.Close()

	// Compress the log file.
	if _, err = io.Copy(gz, file); err != nil {
		fmt.Fprintf(os.Stderr, "Error compressing log file: %s\n", err)
		return
	}

	// Remove the old file.
	if err = os.Remove(name); err != nil {
		fmt.Fprintf(os.Stderr, "Error removing old log file: %s\n", err)
		return
	}

	return
}

// openExistingOrNew tries to open the existing log file.
func (l *Logger) openExistingOrNew(writeLen int) error {
	// Get or create the log file.
	info, err := os.Stat(l.Filename)
	if os.IsNotExist(err) {
		return l.openNew()
	}
	if err != nil {
		return fmt.Errorf("Error getting log file info: %s", err)
	}

	// See if we should rotate the log file.
	if info.Size()+int64(writeLen) >= l.maxSize() {
		return l.rotate()
	}

	// Try to open the current log file.
	file, err := os.OpenFile(l.Filename, os.O_APPEND|os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		// If we fail to open, just ignore and open a new one.
		return l.openNew()
	}

	l.file = file
	l.size = info.Size()

	return nil
}

// openNew tries to open a new log file, creating a backup if one
// already exists.
func (l *Logger) openNew() error {
	// Create the log file directories.
	err := os.MkdirAll(filepath.Dir(l.Filename), 0744)
	if err != nil {
		return fmt.Errorf("Could not create directories for new log file: %s", err)
	}

	mode := os.FileMode(0644)
	info, err := os.Stat(l.Filename)
	if err == nil {
		// Copy mode from the old log file.
		mode = info.Mode()

		// Rename existing log file as backup.
		backup := backupName(l.Filename)
		if err := os.Rename(l.Filename, backup); err != nil {
			return fmt.Errorf("Could not rename log file: %s", err)
		}

		// Compress the backup log file.
		go compressLogFile(backup)
	}

	file, err := os.OpenFile(l.Filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("Could not open new log file: %s", err)
	}

	// Update the instance file info.
	l.file = file
	l.size = 0

	return nil
}

// backupName returns a new backup name for a log file.
func backupName(name string) string {
	// Get the parts of the filepath.
	dir := filepath.Dir(name)
	filename := filepath.Base(name)
	ext := filepath.Ext(filename)
	prefix := filename[:len(filename)-len(ext)]

	// Get a timestamp in RFC3339 format (2006-01-02T15:04:05Z07:00).
	timestamp := time.Now().UTC().Format(time.RFC3339)

	// Return the full path and filename with timestamp.
	return filepath.Join(dir, fmt.Sprintf("%s-%s%s", prefix, timestamp, ext))
}

// maxSize returns the maximum size in bytes of the log file before
// rolling over.
func (l *Logger) maxSize() int64 {
	megabyte := int64(1024 * 1024)
	return l.MaxSize * megabyte
}
