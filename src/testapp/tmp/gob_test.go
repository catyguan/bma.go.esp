package tmp

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"
)

type myval struct {
	V interface{}
}

type P struct {
	X, Y, Z int
	Name    string
}

func TestGOB(t *testing.T) {
	if true {
		gob.Register(map[string]interface{}{})
		gob.Register([]interface{}{})
		gob.Register(P{})
	}

	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.
	dec := gob.NewDecoder(&network) // Will read from network.

	// Encode (send) some values.
	if true {
		o := make(map[string]interface{})
		o["X"] = 3
		o["Y"] = 4
		o["Name"] = "Test"
		a := make([]interface{}, 0)
		a = append(a, 1)
		a = append(a, 2)
		a = append(a, true)
		p := new(P)
		p.X = 1
		p.Y = 2
		p.Z = 3
		p.Name = "abc"
		err := enc.Encode(&myval{o})
		if err != nil {
			t.Errorf("encode error: %s", err)
			return
		}
		fmt.Printf("encoding -> %d\n", network.Len())
	}

	// Decode (receive) and print the values.
	if true {
		// var o myval
		// var o interface{}
		var o map[string]interface{}
		err := dec.Decode(&o)
		if err != nil {
			t.Errorf("decode error: %s", err)
			return
		}
		fmt.Printf("%+v\n", o)
	}
}
