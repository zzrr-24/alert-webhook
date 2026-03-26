package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	err := Init("info", "json")
	assert.NoError(t, err)
	assert.NotNil(t, Logger)

	err = Init("debug", "console")
	assert.NoError(t, err)
	assert.NotNil(t, Logger)
}
