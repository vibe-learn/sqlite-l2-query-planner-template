package main

import (
	"os"
	"testing"
)

// TestDatabasePathDefault — pure unit test. No file needed.
func TestDatabasePathDefault(t *testing.T) {
	if os.Getenv("DATABASE_PATH") == "" && DatabasePath() != ":memory:" {
		t.Errorf("DatabasePath() default = %q, want :memory:", DatabasePath())
	}
}

// TestOpenDBInMemory — открыть :memory:-БД и убедиться, что схема создалась.
// SQLite встроена — серверу подниматься не нужно, тест бежит в CI как есть.
func TestOpenDBInMemory(t *testing.T) {
	db, err := OpenDB(":memory:")
	if err != nil {
		t.Fatalf("OpenDB(:memory:) failed: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		t.Fatalf("ping failed: %v", err)
	}
	// TODO: вызови свои реализованные функции на db и проверь поведение урока «Планировщик запросов».
}
