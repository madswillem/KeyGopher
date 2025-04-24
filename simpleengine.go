package keygopher

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type ReadWriteSeekerCloser interface {
	io.ReadWriteSeeker
	io.Closer
}

type SimpleEngine struct {
	Filepath string
	File     ReadWriteSeekerCloser
}

func InnitSimpleEngine(name string) (SimpleEngine, error) {
	e := SimpleEngine{}
	err := e.Load(name + ".db")
	return e, err
}

func (e SimpleEngine) Close() error {
	if e.File != nil {
		return e.File.Close()
	}
	return nil
}
func (e *SimpleEngine) Load(filepath string) error {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_RDWR, 0644)
	if os.IsNotExist(err) {
		f, err = os.Create(filepath)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	e.Filepath = filepath
	e.File = f

	return nil
}
func (e SimpleEngine) Write(key, value string) error {
	if _, err := e.File.Seek(0, io.SeekEnd); err != nil {
        return err
    }
    line := key + "=" + value + "\n"
    _, err := e.File.Write([]byte(line))
    return err
}
func (e SimpleEngine) Get(key string) (string, error) {
	if _, err := e.File.Seek(0, io.SeekStart); err != nil {
        return "", err
    }
    scanner := bufio.NewScanner(e.File)
    for scanner.Scan() {
        line := scanner.Text()
        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 {
            continue // skip malformed line
        }
        if parts[0] == key {
            return parts[1], nil
        }
    }
    if err := scanner.Err(); err != nil {
        return "", err
    }
    return "", fmt.Errorf("key not found: %s", key)
}
