package keygopher

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)
type BTreeEngine struct {
	Filepath string
	File     ReadWriteSeekerCloser
	BTree	 *BTreeNode
}

type Key struct {
	Name  	string
	Pointer string
}

type BTreeNode struct {
	Children []*BTreeNode
	Keys 	 []Key
	Order   int
}

func InnitBTree(order int) *BTreeNode {
	b := &BTreeNode{
		Children: make([]*BTreeNode, 0, order),
		Keys:     make([]Key, 0, order-1),
		Order:    order,
	}
	return b
}


func (b *BTreeNode) Search(key string) (string, error) {
	for i, k := range b.Keys {
		if strings.Compare(key, k.Name) == 0 {
			return k.Pointer, nil
		}
		if strings.Compare(key, k.Name) < 0 {
			if len(b.Children) > i {
				return b.Children[i].Search(key)
			} else {
				return "", fmt.Errorf("key not found 1: %s", key)
			}
		}
		if i == len(b.Keys)-1 {
			if len(b.Children) > i {
				return b.Children[i+1].Search(key)
			} else {
				return "", fmt.Errorf("key not found 2: %s", key)
			}
		}
	}
	return "", fmt.Errorf("key not found 4: %s", key)
}
func (b *BTreeNode) Add(key string, value string) {
	
}


func InnitBTreeEngine(name string) (BTreeEngine, error) {
	e := BTreeEngine{}
	err := e.Load(name + ".db")
	return e, err
}

func (e BTreeEngine) Close() error {
	if e.File != nil {
		return e.File.Close()
	}
	return nil
}
func (e *BTreeEngine) Load(filepath string) error {
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
func (e BTreeEngine) Write(key, value string) error {
	if _, err := e.File.Seek(0, io.SeekEnd); err != nil {
        return err
    }
    line := key + "=" + value + "\n"
    _, err := e.File.Write([]byte(line))
    return err
}
func (e BTreeEngine) Get(key string) (string, error) {
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
