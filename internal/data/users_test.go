package data

import (
	"strings"
	"testing"
)

func TestPasswordSet(t *testing.T) {
	var p password

	// Test valid password.
	err := p.Set("valid_password123")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if p.plaintext == nil || *p.plaintext != "valid_password123" {
		t.Errorf("expected plaintext to be stored correctly")
	}
	if len(p.hash) == 0 {
		t.Errorf("expected hash to be generated")
	}

	// Test Matches.
	matches, err := p.Matches("valid_password123")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !matches {
		t.Errorf("expected password to match")
	}

	matches, err = p.Matches("wrong_password")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if matches {
		t.Errorf("expected password not to match")
	}

	// Test password exceeding bcrypt limit (> 72 bytes).
	longPassword := strings.Repeat("a", 73)
	err = p.Set(longPassword)
	if err == nil {
		t.Errorf("expected error for password exceeding 72 bytes, got nil")
	}
}
