package response

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDefaultHeaders(t *testing.T) {
	h := GetDefaultHeaders(12)

	contentLength, ok := h.Get("content-length")
	require.True(t, ok)
	assert.Equal(t, "12", contentLength)

	connection, ok := h.Get("connection")
	require.True(t, ok)
	assert.Equal(t, "close", connection)

	contentType, ok := h.Get("content-type")
	require.True(t, ok)
	assert.Equal(t, "text/plain", contentType)
}

func TestWriteStatusLine(t *testing.T) {
	var out bytes.Buffer
	w := NewWriter(&out)

	err := w.WriteStatusLine(StatusInternalServerError)
	require.NoError(t, err)
	assert.Equal(t, "HTTP/1.1 500 Internal Server Error\r\n", out.String())
}

func TestWriteStatusLineUnknownCode(t *testing.T) {
	var out bytes.Buffer
	w := NewWriter(&out)

	err := w.WriteStatusLine(StatusCode(999))
	require.Error(t, err)
	assert.Equal(t, "", out.String())
}