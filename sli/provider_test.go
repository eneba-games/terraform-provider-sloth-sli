package sli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	err := Provider(GeneratorConfigureFunc).InternalValidate()

	assert.NoError(t, err)
}
