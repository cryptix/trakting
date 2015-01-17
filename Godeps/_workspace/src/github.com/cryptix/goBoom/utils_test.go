package goBoom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverse(t *testing.T) {
	assert.Equal(t, reverse("test 123"), "321 tset")
}
