package keygopher

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type SimpleEngine struct {
	Filepath string
	File     *os.File
}

func InnitSimpleEngine(name string) (SimpleEngine, error) {
	e := SimpleEngine{}
	err := e.Load(name + ".db")
	return e, err
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
	err := e.Load(e.Filepath)
	if err != nil {
		return err
	}
	defer e.File.Close()

	string := key + "=" + value + "\n"
	_, err = e.File.WriteString(string)
	return err
}
func (e SimpleEngine) Get(key string) (string, error) {
	err := e.Load(e.Filepath)
	if err != nil {
		return "", err
	}
	defer e.File.Close()

	scanner := bufio.NewScanner(e.File)

	var v string
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		l := strings.Split(scanner.Text(), "=")
		if l[0] == key {
			v = l[1]
		}
	}

	return v, scanner.Err()
}
