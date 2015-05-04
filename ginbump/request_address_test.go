package ginbump

import (
	"bytes"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPublicIP(t *testing.T) {
	parsed := net.ParseIP("10.0.0.15")
	assert.False(t, IsPublicIP(parsed))

	parsed = net.ParseIP("8.8.8.8")
	assert.True(t, IsPublicIP(parsed))

	parsed = net.ParseIP("172.16.0.4")
	assert.False(t, IsPublicIP(parsed))

	parsed = net.ParseIP("8.8.4.4")
	assert.True(t, IsPublicIP(parsed))

	parsed = net.ParseIP("::0")
	assert.False(t, IsPublicIP(parsed))
}

func TestParseForwarded(t *testing.T) {
	parsed := ParseForwarded("10.0.0.1, 8.8.8.8")
	assert.Equal(t, "8.8.8.8", parsed)

	parsed = ParseForwarded("8.8.4.4, 8.8.8.8, 10.0.0.1")
	assert.Equal(t, "8.8.4.4", parsed)

	parsed = ParseForwarded("10.0.0.1")
	assert.Equal(t, "", parsed)

	parsed = ParseForwarded("")
	assert.Equal(t, "", parsed)
}

func TestGetRequesterAddress(t *testing.T) {
	b := bytes.NewBufferString("some body")
	request, _ := http.NewRequest("GET", "test", b)

	address := GetRequesterAddress(request)
	assert.Equal(t, "", address)

	request.RemoteAddr = "127.0.0.1:30475"

	address = GetRequesterAddress(request)
	assert.Equal(t, "127.0.0.1", address)

	request.Header.Set("X-Forwarded-For", "10.0.0.3,::0, 8.8.8.8")

	address = GetRequesterAddress(request)
	assert.Equal(t, "8.8.8.8", address)

	request.Header.Set("X-Forwarded-For", "208.0.0.1, 9.9.4.4")

	address = GetRequesterAddress(request)
	assert.Equal(t, "208.0.0.1", address)
}
