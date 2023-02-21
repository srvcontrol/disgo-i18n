package disgoi18n

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMock(t *testing.T) {

	// Must not panic
	mock := newMock()
	mock.SetDefault(discord.LocaleChineseCN)
	assert.NoError(t, mock.LoadBundle(discord.LocaleSpanishES, ""))
	assert.Empty(t, mock.Get(discord.LocaleCroatian, "", nil))
	assert.Nil(t, mock.GetLocalizations("", nil))

	var called bool
	mock.SetDefaultFunc = func(locale discord.Locale) {
		called = true
	}
	mock.LoadBundleFunc = func(locale discord.Locale, file string) error {
		called = true
		return nil
	}
	mock.GetFunc = func(locale discord.Locale, key string, values Vars) string {
		called = true
		return ""
	}
	mock.GetLocalizationsFunc = func(key string, values Vars) map[discord.Locale]string {
		called = true
		return nil
	}

	called = false
	mock.SetDefault(discord.LocaleChineseCN)
	assert.True(t, called)

	called = false
	assert.NoError(t, mock.LoadBundle(discord.LocaleSpanishES, ""))
	assert.True(t, called)

	called = false
	assert.Empty(t, mock.Get(discord.LocaleCroatian, "", nil))
	assert.True(t, called)

	called = false
	assert.Empty(t, mock.GetLocalizations("", nil))
	assert.True(t, called)
}
