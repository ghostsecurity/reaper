package server

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Server_callMethod(t *testing.T) {
	s := New(nil)
	outputs, err := s.callMethod("Test", []string{`"Hello world!"`})
	require.NoError(t, err)
	require.Len(t, outputs, 1)
	require.Equal(t, `"Hello world!"`, outputs[0])
}
