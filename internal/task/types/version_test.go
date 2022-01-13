package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeVersion(t *testing.T) {
	t.Run("from invalid version should error", func(t *testing.T) {
		_, err := MakeVersion("v.1.2.3")
		assert.NotNil(t, err)
	})

	t.Run("from invalid version should error", func(t *testing.T) {
		_, err := MakeVersion("va.b.c")
		assert.NotNil(t, err)
	})

	t.Run("from latest version", func(t *testing.T) {
		version, err := MakeVersion("latest")
		assert.Nil(t, err)
		assert.True(t, LatestVersion(version))
	})

	t.Run("should return correct parsing", func(t *testing.T) {
		version, err := MakeVersion("v61.12.134")
		assert.Nil(t, err)
		assert.Equal(t, 61, version.major)
		assert.Equal(t, 12, version.minor)
		assert.Equal(t, 134, version.patch)
	})
}

func TestVersionToString(t *testing.T) {
	t.Run("should return correct version string representation", func(t *testing.T) {
		version, err := MakeVersion("v1.2.3")
		assert.Nil(t, err)
		assert.Equal(t, version.String(), "v1.2.3")
	})

	t.Run("should return latest", func(t *testing.T) {
		version, err := MakeVersion("latest")
		assert.Nil(t, err)
		assert.Equal(t, version.String(), "latest")
	})
}
