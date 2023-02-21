package disgoi18n

import (
	"github.com/disgoorg/disgo/discord"
)

var instance translator

func init() {
	instance = newTranslatorImpl()
}

// SetDefault sets the locale used as a fallback.
// Not thread-safe; designed to be called during initialization.
func SetDefault(language discord.Locale) {
	instance.SetDefault(language)
}

// LoadBundle loads a translation file corresponding to a specified locale.
// Not thread-safe; designed to be called during initialization.
func LoadBundle(language discord.Locale, file string) error {
	return instance.LoadBundle(language, file)
}

// Get gets a translation corresponding to a locale and a key.
// Optional Vars parameter is used to inject variables in the translation.
// When a key does not match any translations in the desired locale,
// the default locale is used instead. If the situation persists with the fallback,
// key is returned. If more than one translation is available for dedicated key,
// it is picked randomly. Thread-safe.
func Get(language discord.Locale, key string, values ...Vars) string {
	args := make(Vars)

	for _, variables := range values {
		for variable, value := range variables {
			args[variable] = value
		}
	}

	return instance.Get(language, key, args)
}

// GetLocalizations retrieves translations from every loaded bundles.
// Aims to simplify discord.SlashCommandCreate instantiations by providing
// localizations structures that can be used for any localizable field (example:
// command name, description, etc). Thread-safe.
func GetLocalizations(key string, values ...Vars) map[discord.Locale]string {
	args := make(Vars)

	for _, variables := range values {
		for variable, value := range variables {
			args[variable] = value
		}
	}

	return instance.GetLocalizations(key, args)
}
