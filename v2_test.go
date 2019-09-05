package main

import (
	"testing"
)

func TestCompileMatcher(t *testing.T) {
	t.Errorf("%v\n", compileMatcher([]string{"he", "she", "hers", "his"}))
}
