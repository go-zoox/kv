package test

import (
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/go-zoox/kv/typing"
)

// RunTestCases runs all test cases.
func RunTestCases(t *testing.T, client typing.KV, casesDisabled ...[]string) {
	client.Clear()
	defer client.Clear()

	casesDisabledX := map[string]bool{
		"main":        true,
		"keys":        true,
		"forEach":     true,
		"maxAge":      true,
		"plainValues": true,
	}
	if len(casesDisabled) > 0 {
		for _, c := range casesDisabled[0] {
			casesDisabledX[c] = false
		}
	}

	if casesDisabledX["main"] {
		RunMainTestCase(t, client)
	}

	if casesDisabledX["keys"] {
		RunKeysTestCase(t, client)
	}

	if casesDisabledX["forEach"] {
		RunForEachTestCase(t, client)
	}

	if casesDisabledX["maxAge"] {
		RunMaxAgeTestCase(t, client)
	}

	if casesDisabledX["plainValues"] {
		RunPlainValuesTestCase(t, client)
	}
}

// RunMainTestCase tests the main functionality.
func RunMainTestCase(t *testing.T, client typing.KV) {
	t.Log("Testing main test case")

	if err := client.Clear(); err != nil {
		t.Fatal(err)
	}

	if client.Size() != 0 {
		t.Errorf("Expected size 0, got %d", client.Size())
	}

	if err := client.Set("key", "value"); err != nil {
		t.Fatal(err)
	}
	if !client.Has("key") {
		t.Error("Expected key to be set")
	}
	if client.Size() != 1 {
		t.Errorf("Expected size 1, got %d", client.Size())
	}

	if err := client.Delete("key"); err != nil {
		t.Fatal(err)
	}
	if client.Has("key") {
		t.Error("Expected key to be deleted")
	}
	if client.Size() != 0 {
		t.Errorf("Expected size 0, got %d", client.Size())
	}

	if err := client.Set("key", "value"); err != nil {
		t.Fatal(err)
	}
	if client.Size() != 1 {
		t.Errorf("Expected size 1, got %d", client.Size())
	}
	if err := client.Set("key2", "value2"); err != nil {
		t.Fatal(err)
	}
	if client.Size() != 2 {
		t.Errorf("Expected size 1, got %d", client.Size())
	}

	if err := client.Clear(); err != nil {
		t.Fatal(err)
	}
	if client.Size() != 0 {
		t.Errorf("Expected size 0, got %d", client.Size())
	}

	type User struct {
		Name string
	}
	user1 := &User{
		Name: "Zero",
	}
	if err := client.Set("user", user1); err != nil {
		t.Fatal(err)
	}

	user2 := &User{}
	if err := client.Get("user", user2); err != nil {
		t.Fatal(err)
	}
	if user2.Name != user1.Name {
		t.Fatalf("expect user.name %s, but got %s", user1.Name, user2.Name)
	}

	var user3 *User = nil
	if err := client.Set("user", user3); err != nil {
		t.Fatal(err)
	}

	user4 := &User{}
	if err := client.Get("user", user4); err != nil {
		t.Fatal(err)
	}
	if user4.Name != "" {
		t.Fatalf("expect user.name %s, but got %s", "", user2.Name)
	}
}

// RunKeysTestCase tests the keys functionality.
func RunKeysTestCase(t *testing.T, client typing.KV) {
	t.Log("Testing keys test case")

	if err := client.Clear(); err != nil {
		t.Fatal(err)
	}

	if err := client.Set("key1", "value1"); err != nil {
		t.Fatal(err)
	}
	if err := client.Set("key2", "value2"); err != nil {
		t.Fatal(err)
	}
	if err := client.Set("key3", "value3"); err != nil {
		t.Fatal(err)
	}

	keys := client.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected len 3, got %d", len(keys))
	}
	sort.Slice(keys, func(i, j int) bool {
		return strings.Compare(keys[i], keys[j]) < 0
	})
	if strings.Join(keys, ",") != "key1,key2,key3" {
		t.Errorf("Expected keys to be key1,key2,key3, got %v", strings.Join(keys, ","))
	}
}

// RunForEachTestCase tests the ForEach functionality.
func RunForEachTestCase(t *testing.T, client typing.KV) {
	t.Log("Testing forEach test case")

	if err := client.Clear(); err != nil {
		t.Fatal(err)
	}

	if err := client.Set("key1", "value1"); err != nil {
		t.Fatal(err)
	}
	if err := client.Set("key2", "value2"); err != nil {
		t.Fatal(err)
	}
	if err := client.Set("key3", "value3"); err != nil {
		t.Fatal(err)
	}

	// client.ForEach(func(key string, value interface{}) {
	// 	if key == "key1" && value != "value1" {
	// 		t.Errorf("Expected value to be 'value1', but got %s", value)
	// 	}
	// 	if key == "key2" && value != "value2" {
	// 		t.Errorf("Expected value to be 'value2', but got %s", value)
	// 	}
	// 	if key == "key3" && value != "value3" {
	// 		t.Errorf("Expected value to be 'value3', but got %s", value)
	// 	}
	// })
}

// RunMaxAgeTestCase tests the maxAge functionality.
func RunMaxAgeTestCase(t *testing.T, client typing.KV) {
	t.Log("Testing max age test case")

	client.Clear()
	defer client.Clear()

	key1 := "key1"
	value1 := "value1"
	if err := client.Set(key1, &value1, 500); err != nil {
		t.Fatal(err)
	}

	var value string
	if err := client.Get("key1", &value); err != nil || value != value1 {
		t.Error("Expected value to be 'value1'")
	}

	// fmt.Println("xxx1:", client.Get("key1"))
	time.Sleep(2 * time.Second)
	// done := make(chan bool)
	// go func() {
	// 	fmt.Println("xxx2:", client.Get("key1"))
	// 	time.Sleep(2 * time.Second)
	// 	fmt.Println("xxx3:", client.Get("key1"))
	// 	done <- true
	// }()
	// <-done
	// fmt.Println("xxx4:", client.Get("key1"))

	client.Get("key1", &value)
	if client.Has("key1") {
		t.Errorf("Expected value to be '', but got %s", value)
	}
}

// RunPlainValuesTestCase tests Get/Set with plain (non-pointer) values like string, int, bool.
func RunPlainValuesTestCase(t *testing.T, client typing.KV) {
	t.Log("Testing plain values test case")

	if err := client.Clear(); err != nil {
		t.Fatal(err)
	}

	// plain string
	if err := client.Set("s", "hello"); err != nil {
		t.Fatal(err)
	}
	var s string
	if err := client.Get("s", &s); err != nil {
		t.Fatalf("Get plain string: %v", err)
	}
	if s != "hello" {
		t.Fatalf("expected 'hello', got %q", s)
	}

	// plain int
	if err := client.Set("i", 42); err != nil {
		t.Fatal(err)
	}
	var i int
	if err := client.Get("i", &i); err != nil {
		t.Fatalf("Get plain int: %v", err)
	}
	if i != 42 {
		t.Fatalf("expected 42, got %d", i)
	}

	// plain bool
	if err := client.Set("b", true); err != nil {
		t.Fatal(err)
	}
	var b bool
	if err := client.Get("b", &b); err != nil {
		t.Fatalf("Get plain bool: %v", err)
	}
	if !b {
		t.Fatal("expected true")
	}

	// plain float64
	if err := client.Set("f", 3.14); err != nil {
		t.Fatal(err)
	}
	var f float64
	if err := client.Get("f", &f); err != nil {
		t.Fatalf("Get plain float64: %v", err)
	}
	if f != 3.14 {
		t.Fatalf("expected 3.14, got %f", f)
	}

	// plain int64
	if err := client.Set("i64", int64(999)); err != nil {
		t.Fatal(err)
	}
	var i64 int64
	if err := client.Get("i64", &i64); err != nil {
		t.Fatalf("Get plain int64: %v", err)
	}
	if i64 != 999 {
		t.Fatalf("expected 999, got %d", i64)
	}

	// string slice
	if err := client.Set("ss", []string{"a", "b", "c"}); err != nil {
		t.Fatal(err)
	}
	var ss []string
	if err := client.Get("ss", &ss); err != nil {
		t.Fatalf("Get string slice: %v", err)
	}
	if len(ss) != 3 || ss[0] != "a" || ss[2] != "c" {
		t.Fatalf("expected [a b c], got %v", ss)
	}

	// pointer to pointer (double pointer) should work
	type User struct{ Name string }
	user := &User{Name: "Zero"}
	if err := client.Set("user_ptr", user); err != nil {
		t.Fatal(err)
	}
	var result User
	if err := client.Get("user_ptr", &result); err != nil {
		t.Fatalf("Get struct pointer: %v", err)
	}
	if result.Name != "Zero" {
		t.Fatalf("expected Zero, got %q", result.Name)
	}

	// nil pointer should be retrievable
	var nilPtr *User = nil
	if err := client.Set("nil_ptr", nilPtr); err != nil {
		t.Fatal(err)
	}
	var nilResult *User
	if err := client.Get("nil_ptr", &nilResult); err != nil {
		t.Fatalf("Get nil pointer: %v", err)
	}
	if nilResult != nil {
		t.Fatal("expected nil result")
	}

	// overwrite plain with pointer and back
	if err := client.Set("x", "plain"); err != nil {
		t.Fatal(err)
	}
	var x1 string
	client.Get("x", &x1)
	if x1 != "plain" {
		t.Fatal("expected 'plain'")
	}

	if err := client.Set("x", &x1); err != nil {
		t.Fatal(err)
	}
	var x2 string
	client.Get("x", &x2)
	if x2 != "plain" {
		t.Fatal("expected 'plain' after pointer set")
	}
}
