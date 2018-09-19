package main

import (
	"log"
	"os"
	"testing"
)

func TestSublevelSearch(t *testing.T) {
	raw := []byte(`{
    "name": "gulyasm",
    "age": 12,
    "price": 876.2341234,
    "fake": {
        "zip": "H-1056"
	},
    "address": {
        "zip": "H-1056",
        "city": "Budapest",
        "gateways": ["Sopron", "Vienna", "Budapest"]
		}
	}`)
	tree, err := fromBytes(raw)
	if err != nil {
		t.Fatalf("failed to convert JSON to tree")
	}
	v, ok := tree.(*complexNode)
	if !ok {
		t.Fatalf("failed to convert tree to complexNode")
	}
	subtree, err := v.search("ga")
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
	log.Println(len(v2.data))
	v2.draw(os.Stdout, 2, 0)
	if len(v2.data) != 1 {
		t.Fatalf("root complex node should have a single child")
	}
	_, ok = v2.data["address"]
	if !ok {
		t.Fatalf("root complex node should have an address child")
	}

}

func TestSingleResult(t *testing.T) {
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

	subtree, err := tree.search("a")
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

	if len(v2.data) != 4 {
		t.Fatalf("searched subtree element should have 4 children")
	}
}
