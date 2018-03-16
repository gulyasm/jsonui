package main

import "testing"

func TestListTree(t *testing.T) {
	raw := []byte(`[1,234,35]`)
	tree, err := FromBytes(raw)
	if err != nil {
		t.Fatalf("failed to convert JSON to tree")
	}
	v, ok := tree.(*ListNode)
	if !ok {
		t.Fatalf("root element should be a ListNode")
	}
	if len(v.data) != 3 {
		t.Fatalf("root element should have 3 children")
	}
}

func TestEmpty(t *testing.T) {
	var tp TreePosition = []string{"alma", "barack", "banan"}
	if tp.Empty() {
		t.Fatalf("non-empty position returned true for Empty")
	}
}

func TestEmpty2(t *testing.T) {
	var tp TreePosition = []string{}
	if !tp.Empty() {
		t.Fatalf("empty position returned false for Empty")
	}
}

// TODO(gulyasm): nicer fail messages
func TestShift(t *testing.T) {
	var tp TreePosition = []string{"alma", "barack", "banan"}
	var expected TreePosition = TreePosition([]string{"barack", "banan"})
	new := tp.Shift()
	if len(new) != 2 {
		t.Fatalf("length is not 2 after Shift")
	}
	for i := 0; i < len(expected); i++ {
		if new[i] != expected[i] {
			t.Fatalf("shift did not return what was expected")
		}
	}
}
