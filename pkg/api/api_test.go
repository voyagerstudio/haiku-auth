package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO: Make this a real test
// Just made this to get CI working
func TestNewServer(t *testing.T) {
	s := NewServer("", 0, nil)
	require.NotNil(t, s)
}
