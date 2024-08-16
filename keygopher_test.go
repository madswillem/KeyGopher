package keygopher

import (
	"os"
	"syscall"
	"testing"
)

func TestNew(t *testing.T) {
	defer os.Remove("test.db")

	err, _ := New(&Config{Name: "test"})
	if err != nil {
		t.Errorf("Error: %+v", err)
	}
}
func TestWrite(t *testing.T) {
	defer os.Remove("test.db")

	err, db := New(&Config{Name: "test"})
	if err != nil {
		t.Errorf("Error while loading: %e, %v", err, syscall.Errno(9))
	}

	err = db.Write("test", "1234")
	if err != nil {
		t.Errorf("Error while writing: %e", err)
	}
}
func TestRead(t *testing.T) {
	defer os.Remove("test.db")

	err, db := New(&Config{Name: "test"})
	if err != nil {
		t.Errorf("Error while loading: %e, %v", err, syscall.Errno(9))
	}

	err = db.Write("test", "1234")
	if err != nil {
		t.Errorf("Error while writing: %e", err)
	}

	r, err := db.Get("test")
	if err != nil {
		t.Errorf("Error while reading: %e", err)
	}

	if r != "1234" {
		t.Errorf("r not %s but %s", "1234", r)
	}
}
