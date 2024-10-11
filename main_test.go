package main

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestName(t *testing.T) {
	input := "\\path\\to\\file"
	output := filepath.FromSlash(input)
	fmt.Println(output)
}
