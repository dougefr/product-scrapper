package businesserr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBusinessError(t *testing.T) {
	be := NewBusinessError("code", "error")
	assert.Equal(t, "code", be.Code())
	assert.Equal(t, "error", be.Error())
}
