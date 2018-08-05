package timeutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSecond(t *testing.T) {

	tests := map[int]int{
		1: 86400,
		7: 604800,
		0: 0,
	}

	for days, expected := range tests {

		assert := assert.New(t)

		assert.Equal(expected, GetSeconds(days))
	}
}
