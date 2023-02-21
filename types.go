package disgoi18n

import "github.com/disgoorg/disgo/discord"

// Vars is the collection used to inject variables during translation.
// This type only exists to be less verbose.
type Vars map[string]interface{}

type translator interface {
	SetDefault(locale discord.Locale)
	LoadBundle(locale discord.Locale, file string) error
	Get(locale discord.Locale, key string, values Vars) string
	GetLocalizations(key string, variables Vars) map[discord.Locale]string
}

type translatorImpl struct {
	defaultLocale discord.Locale
	translations  map[discord.Locale]bundle
	loadedBundles map[string]bundle
}

type translatorMock struct {
	SetDefaultFunc       func(locale discord.Locale)
	LoadBundleFunc       func(locale discord.Locale, file string) error
	GetFunc              func(locale discord.Locale, key string, values Vars) string
	GetLocalizationsFunc func(key string, variables Vars) map[discord.Locale]string
}

type bundle map[string][]string
