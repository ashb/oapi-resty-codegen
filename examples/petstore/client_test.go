package petstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient("a")
	assert.Nil(t, err)
	assert.NotNil(t, client)
}
