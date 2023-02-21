package disgoi18n

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFacade(t *testing.T) {
	var expectedFile, expectedKey = "File", "Key"
	var expectedValues Vars
	var called bool

	mock := newMock()
	mock.SetDefaultFunc = func(locale discord.Locale) {
		assert.Equal(t, discord.LocaleItalian, locale)
		called = true
	}
	mock.LoadBundleFunc = func(locale discord.Locale, file string) error {
		assert.Equal(t, discord.LocaleFrench, locale)
		assert.Equal(t, expectedFile, file)
		called = true
		return nil
	}
	mock.GetFunc = func(locale discord.Locale, key string, values Vars) string {
		assert.Equal(t, discord.LocaleChineseCN, locale)
		assert.Equal(t, expectedValues, values)
		assert.Equal(t, expectedKey, key)
		called = true
		return ""
	}
	mock.GetLocalizationsFunc = func(key string, values Vars) map[discord.Locale]string {
		assert.Equal(t, expectedValues, values)
		assert.Equal(t, expectedKey, key)
		called = true
		return nil
	}

	instance = mock

	called = false
	SetDefault(discord.LocaleItalian)
	assert.True(t, called)

	called = false
	assert.NoError(t, LoadBundle(discord.LocaleFrench, expectedFile))
	assert.True(t, called)

	called = false
	expectedValues = make(Vars)
	Get(discord.LocaleChineseCN, expectedKey)
	assert.True(t, called)

	called = false
	expectedValues = Vars{
		"Hi": "There",
	}
	Get(discord.LocaleChineseCN, expectedKey, expectedValues)
	assert.True(t, called)

	called = false
	expectedValues = Vars{
		"Hi":  "There",
		"Bye": "See u",
	}
	Get(discord.LocaleChineseCN, expectedKey, Vars{"Hi": "There"}, Vars{"Bye": "See u"})
	assert.True(t, called)

	called = false
	GetLocalizations(expectedKey, Vars{"Hi": "There"}, Vars{"Bye": "See u"})
	assert.True(t, called)
}
