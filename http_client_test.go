package main

import (
	"fmt"
	"testing"
)

func TestHttpClient_Get(t *testing.T) {
	client, err := NewHttpClient("https://www.google.com")
	if err != nil {
		t.Errorf("Failed to create a new client: %v", err)
		return
	}

	result, err := client.DoPost()
	if err != nil {
		t.Errorf("Failed to do http get: %v", err)
	}

	fmt.Printf("result: %v\n", result)
}

func TestHttpClient_Get2(t *testing.T) {
	addr := "https://en.wikipedia.org/wiki/Wikipedia:Community_portal"
	client, err := NewHttpClient(addr)
	if err != nil {
		t.Errorf("Failed to create a new client: %v", err)
		return
	}

	result, err := client.DoPost()
	if err != nil {
		t.Errorf("Failed to do http get: %v", err)
	}

	fmt.Printf("result: %v\n", result)
}

func TestHttpClient_Get3(t *testing.T) {
	addr := "http://localhost:28080/memwork.php/?value=110&memory=10"
	client, err := NewHttpClient(addr)
	if err != nil {
		t.Errorf("Failed to create a new client: %v", err)
		return
	}

	result, err := client.DoPost()
	if err != nil {
		t.Errorf("Failed to do http get: %v", err)
	}

	fmt.Printf("result: %v\n", result)
}
