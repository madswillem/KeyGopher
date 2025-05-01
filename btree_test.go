package keygopher

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

type DummieFile struct {
}

func (d *DummieFile) Read(p []byte) (n int, err error) {
	return 0, nil
}
func (d *DummieFile) Write(p []byte) (n int, err error) {
	fmt.Println(string(p))
	return len(p), nil
}
func (d *DummieFile) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}
func (d *DummieFile) Close() error {
	return nil
}

func TestBTreeSearch(t *testing.T) {
	// Create a sample B-tree (order 3)
	tree := InnitBTree(5)

	// Create child nodes
	child1 := InnitBTree(5)
	child2 := InnitBTree(5)
	child3 := InnitBTree(5)

	// Add keys to child nodes
	child1.Keys = append(child1.Keys, Key{Name: "apricot", Pointer: "101"})
	child1.Keys = append(child1.Keys, Key{Name: "avocado", Pointer: "201"})
	child2.Keys = append(child2.Keys, Key{Name: "blueberry", Pointer: "301"})
	child2.Keys = append(child2.Keys, Key{Name: "cranberry", Pointer: "401"})
	child3.Keys = append(child3.Keys, Key{Name: "fig", Pointer: "501"})
	child3.Keys = append(child3.Keys, Key{Name: "grape", Pointer: "601"})

	// Add keys to the root node
	tree.Keys = append(tree.Keys, Key{Name: "banana", Pointer: "200"})
	tree.Keys = append(tree.Keys, Key{Name: "date", Pointer: "400"})

	// Set children of the root node
	tree.Children = append(tree.Children, child1)
	tree.Children = append(tree.Children, child2)
	tree.Children = append(tree.Children, child3)

	// Perform a search for a key in the root node
	value, err := tree.Search("banana")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if value != "200" {
		t.Errorf("Expected value '200', got '%s'", value)
	}

	// Perform a search for a key in the first child node
	value, err = tree.Search("avocado")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if value != "201" {
		t.Errorf("Expected value '201', got '%s'", value)
	}

	// Perform a search for a key in the second child node
	value, err = tree.Search("cranberry")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if value != "401" {
		t.Errorf("Expected value '401', got '%s'", value)
	}

	// Perform a search for a key in the third child node
	value, err = tree.Search("fig")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if value != "501" {
		t.Errorf("Expected value '501', got '%s'", value)
	}

	// Search for a non-existent key
	_, err = tree.Search("abcde")
	if err == nil {
		t.Error("Expected 'key not found' error, got nil")
	}
}

func TestBTreeAdd(t *testing.T) {
	// Create an empty B-tree (order 5)
	tree := InnitBTree(5)

	// List of keys to add in order with value same as key.
	keysToAdd := []string{
		"A", "G", "F", "B", "K", "D", "H", "M", "J",
		"E", "S", "I", "R", "X", "C", "L", "N", "T", "U", "P",
	}

	for _, k := range keysToAdd {
		tree.Add(k, k) // Using key as its own value
		tree.PrintPyramid()		
		fmt.Println("____________________________________________________")
	}
	fmt.Println("Final B-tree structure:")
	tree.PrintPyramid()
	fmt.Println("First level keys:")
	fmt.Println(tree.Keys)
	fmt.Println("Second layer keys:")
	for _, child := range tree.Children {
		fmt.Println(child.Keys)
	}
	fmt.Println("____________________________________________________")


	// Verify that each key is retrievable with its corresponding value.
	for _, k := range keysToAdd {
		value, err := tree.Search(k)
		if err != nil {
			t.Errorf("Failed to find key %s: %v", k, err)
		}
		if value != k {
			t.Errorf("Expected value '%s' for key '%s', got '%s'", k, k, value)
		}
	}

	// Verify that the tree structure is correct thoroughly.
	// Helper: in-order traversal of the B-tree keys
	inOrderKeys := []string{}
	var traverse func(node *BTreeNode)
	traverse = func(node *BTreeNode) {
		if len(node.Children) == 0 {
			for _, k := range node.Keys {
				inOrderKeys = append(inOrderKeys, k.Name)
			}
		} else {
			for i, k := range node.Keys {
				traverse(node.Children[i])
				inOrderKeys = append(inOrderKeys, k.Name)
			}
			traverse(node.Children[len(node.Children)-1])
		}
	}
	traverse(tree)

	// Prepare expected sorted order of keys
	expected := make([]string, len(keysToAdd))
	copy(expected, keysToAdd)
	sort.Strings(expected)

	// Compare the in-order traversal with the expected sorted keys
	if !reflect.DeepEqual(inOrderKeys, expected) {
		t.Errorf("Expected in-order traversal to be %v, got %v", expected, inOrderKeys)
	}
}
