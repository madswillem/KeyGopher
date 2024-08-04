package keygopher

import (
	"bufio"
	"os"
	"strings"
)

type DB struct {
	FilePath string
	File     *os.File
}

func New() (error, *DB) {
	f, err := os.Create("testFile.txt")
	if err != nil {
		return err, nil
	}
	db := DB{
		FilePath: f.Name(),
		File:     f,
	}
	return nil, &db
}
func Load(path string) (error, *DB) {
	f, err := os.OpenFile("./testFile.txt", os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return err, nil
	}
	if err != nil {
		return err, nil
	}
	db := DB{
		FilePath: path,
		File:     f,
	}
	return nil, &db
}

func (db *DB) Write(key, value string) error {
	string := key + "=" + value + "\n"
	_, err := db.File.WriteString(string)
	return err
}
func (db *DB) Read(key string) (error, string) {
	scanner := bufio.NewScanner(db.File)

	var v string
	for scanner.Scan() {
		l := strings.Split(scanner.Text(), "=")
		if l[0] == key {
			v = l[1]
		}
	}

	return scanner.Err(), v
}
