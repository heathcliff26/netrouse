package main

import (
	"strings"
	"testing"

	"github.com/heathcliff26/netrouse/pkg/version"
	"github.com/stretchr/testify/assert"
)

func TestNewRootCommand(t *testing.T) {
	cmd := NewRootCommand()

	name := strings.ToLower(version.Name)
	assert.Equal(t, name, cmd.Use)
}
