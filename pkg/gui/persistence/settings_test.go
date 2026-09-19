package persistence

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadSettings(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		assert := assert.New(t)

		oldFolder := ConfigFolder()
		t.Cleanup(func() {
			SetConfigFolder(oldFolder)
		})
		configFolder = "testdata/valid"

		expectedSettings := &Settings{
			Remotes: []RemoteServer{
				{
					Name: "TestServer",
					URL:  "https://example.com",
				},
			},
		}

		s, err := LoadSettings()
		assert.NoError(err, "Should load settings")
		assert.Equal(expectedSettings, s, "Should load the same settings")
	})
	t.Run("NoSettings", func(t *testing.T) {
		assert := assert.New(t)

		oldFolder := ConfigFolder()
		t.Cleanup(func() {
			SetConfigFolder(oldFolder)
		})
		configFolder = "testdata/nothing"

		s, err := LoadSettings()
		assert.NoError(err, "Should load settings")
		assert.Equal(&Settings{}, s, "Should return empty settings")
	})
	t.Run("InvalidSettings", func(t *testing.T) {
		assert := assert.New(t)

		oldFolder := ConfigFolder()
		t.Cleanup(func() {
			SetConfigFolder(oldFolder)
		})
		configFolder = "testdata/invalid"

		s, err := LoadSettings()
		assert.Error(err, "Should load settings")
		assert.Equal(&Settings{}, s, "Should return empty settings")
	})
}

func TestSaveSettings(t *testing.T) {
	require := require.New(t)

	oldFolder := ConfigFolder()
	t.Cleanup(func() {
		SetConfigFolder(oldFolder)
	})
	configFolder = t.TempDir()

	s := &Settings{
		Remotes: []RemoteServer{
			{
				Name: "TestServer",
				URL:  "https://example.com",
			},
		},
	}

	require.NoError(s.Save(), "Should save settings")

	loadedSettings, err := LoadSettings()
	require.NoError(err, "Should load settings")

	require.Equal(s, loadedSettings, "Should load the same settings")
}
