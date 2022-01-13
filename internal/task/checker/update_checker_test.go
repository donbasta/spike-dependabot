package checker

import (
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
)

func TestGetSourceTypeFromUrl(t *testing.T) {
	url := faker.URL()
	_, err := getSourceTypeFromURL(url)
	assert.NotNil(t, err)

	url = "https://source.golabs.io/cloud-platform/automation/stateful-component-portal"
	sourceType, err := getSourceTypeFromURL(url)
	assert.Nil(t, err)
	assert.Equal(t, "golabs", sourceType)

	url = "https://github.com/UnderGreen/ansible-role-mongodb"
	sourceType, err = getSourceTypeFromURL(url)
	assert.Nil(t, err)
	assert.Equal(t, "github", sourceType)
}
