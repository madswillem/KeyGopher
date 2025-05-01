package keygopher

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type BTreeEngine struct {
	Filepath string
	File     ReadWriteSeekerCloser
	BTree    *BTreeNode
}

type Key struct {
	Name    string
	Pointer string
}

type BTreeNode struct {
	Parent   *BTreeNode
	Children []*BTreeNode
	Keys     []Key
	Order    int
}

func InnitBTree(order int) *BTreeNode {
	b := &BTreeNode{
		Children: make([]*BTreeNode, 0, order),
		Keys:     make([]Key, 0, order-1),
		Order:    order,
	}
	return b
}

func (b *BTreeNode) GetNodeByKey(key string) *BTreeNode {
	for i, k := range b.Keys {
		if strings.Compare(key, k.Name) == 0 {
			return nil
		}
		if strings.Compare(key, k.Name) < 0 {
			if len(b.Children) > i {
				return b.Children[i].GetNodeByKey(key)
			} else {
				return b
			}
		}
		if i == len(b.Keys)-1 {
			if len(b.Children) > i {
				return b.Children[i+1].GetNodeByKey(key)
			} else {
				return b
			}
		}
	}
	return b
}

func (b *BTreeNode) InsertChild(node *BTreeNode) error {
	first_key := node.Keys[0].Name
	for i, child := range b.Children {
		if strings.Compare(child.Keys[0].Name, first_key) > 0 {
			if i == 0 {
				b.Children = append([]*BTreeNode{node}, b.Children...)
				return nil
			}
			b.Children = append(b.Children[:i], append([]*BTreeNode{node}, b.Children[i:]...)...)
			return nil
		}
		if i == len(b.Children)-1 {
			b.Children = append(b.Children, node)
			return nil
		}
	}
	return nil;
}

func splitAtMedian(node *BTreeNode) (*BTreeNode, *BTreeNode, Key) {
	// Ensure keys are in order.
	sort.Slice(node.Keys, func(i, j int) bool {
		return strings.Compare(node.Keys[i].Name, node.Keys[j].Name) < 0
	})

	mid := len(node.Keys) / 2
	median := node.Keys[mid]

	// Split keys.
	leftKeys := make([]Key, mid)
	copy(leftKeys, node.Keys[:mid])
	rightKeys := make([]Key, len(node.Keys)-mid-1)
	copy(rightKeys, node.Keys[mid+1:])

	// Split children if present.
	var leftChildren, rightChildren []*BTreeNode
	if len(node.Children) > 0 {
		// In a B-tree, number of children is keys+1.
		leftChildCount := len(leftKeys) + 1
		leftChildren = make([]*BTreeNode, leftChildCount)
		copy(leftChildren, node.Children[:leftChildCount])
		rightChildren = make([]*BTreeNode, len(node.Children)-leftChildCount)
		copy(rightChildren, node.Children[leftChildCount:])
	}

	leftNode := &BTreeNode{
		Parent:   node.Parent,
		Keys:     leftKeys,
		Children: leftChildren,
		Order:    node.Order,
	}
	for _, child := range leftNode.Children {
		child.Parent = leftNode
	}
	rightNode := &BTreeNode{
		Parent:   node.Parent,
		Keys:     rightKeys,
		Children: rightChildren,
		Order:    node.Order,
	}
	for _, child := range leftNode.Children {
		child.Parent = rightNode
	}
	return leftNode, rightNode, median
}
func (b *BTreeNode) PrintPyramid() {
    // Gather tree levels via BFS.
    levels := [][]*BTreeNode{}
    queue := []*BTreeNode{b}
    for len(queue) > 0 {
        levelNodes := []*BTreeNode{}
        nextQueue := []*BTreeNode{}
        for _, node := range queue {
            levelNodes = append(levelNodes, node)
            nextQueue = append(nextQueue, node.Children...)
        }
        levels = append(levels, levelNodes)
        queue = nextQueue
    }

    // Use the width of the last level for spacing.
    if len(levels) == 0 {
        return
    }
    maxWidth := len(levels[len(levels)-1])

    // Print each level with padding to simulate a pyramid.
    for _, level := range levels {
        // Calculate indent based on the difference between max width and current level count.
        indentCountRaw := (maxWidth - len(level)) * 2
        // Ensure indentCount is not negative.
        indentCount := indentCountRaw
        if indentCount < 0 {
            indentCount = 0
        }
        indent := strings.Repeat(" ", indentCount)
        line := indent
        spaceBetween := strings.Repeat(" ", indentCount*2+2)
        for j, node := range level {
            // Represent each node as [key1,key2,...]
            s := "["
            for k, key := range node.Keys {
                if k > 0 {
                    s += ","
                }
                s += key.Name
            }
            s += "]"
            line += s
            if j < len(level)-1 {
                line += spaceBetween
            }
        }
        fmt.Println(line)
    }
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

func (b *BTreeNode) AddKey(key Key){
	if len(b.Keys) < b.Order-1 {
		b.Keys = append(b.Keys, key)
		sort.Slice(b.Keys, func(i, j int) bool {
			return strings.Compare(b.Keys[i].Name, b.Keys[j].Name) < 0
		},
		)
		return
	}
	b.Keys = append(b.Keys, key)
	left, right, median := splitAtMedian(b)
	if b.Parent != nil {
		b.Children = left.Children;
		b.Keys = left.Keys
		b.Parent.InsertChild(right)
		b.Parent.AddKey(median)
		return
	}
	b.Keys = make([]Key, 0, b.Order-1)
	b.Keys = append(b.Keys, median)
	b.Children = make([]*BTreeNode, 0, b.Order)
	left.Parent = b
	right.Parent = b
	b.Children = append(b.Children, left)
	b.Children = append(b.Children, right)
}

func (b *BTreeNode) Add(key string, value string) {
	node := b.GetNodeByKey(key)
	fmt.Println("Node to add to:", node)
	if node == nil {
		return
	}
	node.AddKey(Key{
		Name:    key,
		Pointer: value,
	})
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
