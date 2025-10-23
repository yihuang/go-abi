package abi

import (
	"testing"
)

func TestMinimalDirectTuple(t *testing.T) {
	// Test the most basic case first
	humanABI := []string{
		"function test((uint256, uint256) pair)",
	}

	_, err := ParseHumanReadableABI(humanABI)
	if err != nil {
		t.Fatalf("Failed to parse minimal direct tuple: %v", err)
	}

	t.Log("✅ Minimal direct tuple parsing works!")
}