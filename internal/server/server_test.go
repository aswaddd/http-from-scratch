package server

import (
	"bytes"
	"io"
	"testing"

	"aswad.http.module/internal/request"
	"aswad.http.module/internal/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testConn struct {
	reader *bytes.Reader
	out    bytes.Buffer
	closed bool
}

func newTestConn(in string) *testConn {
	return &testConn{reader: bytes.NewReader([]byte(in))}
}

func (c *testConn) Read(p []byte) (int, error) {
	return c.reader.Read(p)
}

func (c *testConn) Write(p []byte) (int, error) {
	return c.out.Write(p)
}

func (c *testConn) Close() error {
	c.closed = true
	return nil
}

func TestRunConnectionBadRequest(t *testing.T) {
	s := &Server{}
	conn := newTestConn("NOT A VALID REQUEST")

	runConnection(s, conn)

	assert.True(t, conn.closed)
	assert.Contains(t, conn.out.String(), "HTTP/1.1 400 Bad Request\r\n")
	assert.Contains(t, conn.out.String(), "connection: close\r\n")
}

func TestRunConnectionCallsHandler(t *testing.T) {
	called := false
	handler := func(w *response.Writer, req *request.Request) {
		called = true
		_ = w.WriteStatusLine(response.StatusOK)
		h := response.GetDefaultHeaders(5)
		_ = w.WriteHeaders(*h)
		_, _ = w.WriteBody([]byte("hello"))
	}

	s := &Server{handler: handler}
	conn := newTestConn("GET / HTTP/1.1\r\nHost: localhost\r\n\r\n")

	runConnection(s, conn)

	require.True(t, called)
	assert.True(t, conn.closed)
	assert.Contains(t, conn.out.String(), "HTTP/1.1 200 OK\r\n")
	assert.Contains(t, conn.out.String(), "hello")
}

func TestRunConnectionReaderError(t *testing.T) {
	conn := &errorConn{}
	s := &Server{}

	runConnection(s, conn)

	assert.True(t, conn.closed)
	assert.Contains(t, conn.out.String(), "HTTP/1.1 400 Bad Request\r\n")
}

type errorConn struct {
	out    bytes.Buffer
	closed bool
}

func (c *errorConn) Read(_ []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func (c *errorConn) Write(p []byte) (int, error) {
	return c.out.Write(p)
}

func (c *errorConn) Close() error {
	c.closed = true
	return nil
}
