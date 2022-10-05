package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCustomResponseError(t *testing.T) {
	err := NewCustomResponseError(500, "message")

	assert.Equal(t, 500, err.StatusCode())
	assert.Equal(t, "message", err.Error())
	assert.Equal(t, "message", err.Message())
}

func TestSendBusinessErrorJson(t *testing.T) {

}
