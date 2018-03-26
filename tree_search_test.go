package main

import "testing"

func TestSearchSimpleKeys(t *testing.T) {
	raw := []byte(`{
		"alma": "mu",
		"name": "barack",
		"age": 12,
		"tags": ["good", "excellent", "ceu"]
	}`)
	tree, err := fromBytes(raw)
	if err != nil {
		t.Fatalf("failed to convert JSON to tree")
	}
	v, ok := tree.(*complexNode)
	if !ok {
		t.Fatalf("failed to convert tree to complexNode")
	}

	if len(v.data) != 4 {
		t.Fatalf("root element should have 4 children")
	}

	subtree, err := tree.search("alma")
	if err != nil {
		t.Fatalf("failed to search tree: %q", err.Error())
	}
	if subtree == nil {
		t.Fatalf("subtree returned nil")
	}
	v2, ok := subtree.(*complexNode)
	if !ok {
		t.Fatalf("failed to convert tree to complexNode")
	}

	if len(v2.data) != 1 {
		t.Fatalf("searched subtree element should have 1 children")
	}
	anode, ok := v2.data["alma"]
	if !ok {
		t.Fatalf("first node should be alma node")
	}
	if anode.String(0, 0) != "\"mu\"" {
		t.Fatalf("searched subtree element should be mu. Instead it was %q", anode.String(0, 0))
	}
}
