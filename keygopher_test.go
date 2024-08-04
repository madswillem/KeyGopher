package keygopher

import (
	"syscall"
	"testing"
)

func TestNew(t *testing.T) {
	err, _ := New()
	if err != nil {
		t.Errorf("Error: %e", err)
	}
}
func TestWrite(t *testing.T) {
	err, db := Load("testFile.txt")
	if err != nil {
		t.Errorf("Error while loading: %e, %v", err, syscall.Errno(9))
	}

	err = db.Write("test", "1234")
	if err != nil {
		t.Errorf("Error while writing: %e", err)
	}
}
func TestRead(t *testing.T) {
	err, db := Load("testFile.txt")
	if err != nil {
		t.Errorf("Error while loading: %e, %v", err, syscall.Errno(9))
	}

	err, r := db.Read("test")
	if err != nil {
		t.Errorf("Error while reading: %e", err)
	}

	if r != "1234" {
		t.Errorf("r not %s but %s", "1234", r)
	}
}
