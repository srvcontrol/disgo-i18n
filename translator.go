package disgoi18n

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"math/rand"
	"os"
	"strings"
	"text/template"
	"time"
)

const (
	defaultLocale   = discord.LocaleEnglishUS
	leftDelim       = "{{"
	rightDelim      = "}}"
	keyDelim        = "."
	executionPolicy = "missingkey=error"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func newTranslatorImpl() *translatorImpl {
	return &translatorImpl{
		defaultLocale: defaultLocale,
		translations:  make(map[discord.Locale]bundle),
		loadedBundles: make(map[string]bundle),
	}
}

func (translator *translatorImpl) SetDefault(language discord.Locale) {
	translator.defaultLocale = language
}

func (translator *translatorImpl) LoadBundle(locale discord.Locale, path string) error {
	loadedBundle, found := translator.loadedBundles[path]
	if !found {

		buf, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var jsonContent map[string]interface{}
		err = json.Unmarshal(buf, &jsonContent)
		if err != nil {
			return err
		}

		newBundle := translator.mapBundleStructure(jsonContent)

		translator.loadedBundles[path] = newBundle
		translator.translations[locale] = newBundle

	} else {
		translator.translations[locale] = loadedBundle
	}

	return nil
}

func (translator *translatorImpl) Get(locale discord.Locale, key string, variables Vars) string {
	bundles, found := translator.translations[locale]
	if !found {
		if locale != translator.defaultLocale {
			return translator.Get(translator.defaultLocale, key, variables)
		} else {
			return key
		}
	}

	raws, found := bundles[key]
	if !found || len(raws) == 0 {
		if locale != translator.defaultLocale {
			return translator.Get(translator.defaultLocale, key, variables)
		} else {
			return key
		}
	}

	raw := raws[rand.Intn(len(raws))]

	if variables != nil && strings.Contains(raw, leftDelim) {
		t, err := template.New("").Delims(leftDelim, rightDelim).Option(executionPolicy).Parse(raw)
		if err != nil {
			return raw
		}

		var buf bytes.Buffer
		err = t.Execute(&buf, variables)
		if err != nil {
			return raw
		}
		return buf.String()
	}

	return raw
}

func (translator *translatorImpl) GetLocalizations(key string, variables Vars) map[discord.Locale]string {
	localizations := make(map[discord.Locale]string)

	for locale := range translator.translations {
		localizations[locale] = translator.Get(locale, key, variables)
	}

	return localizations
}

func (translator *translatorImpl) mapBundleStructure(jsonContent map[string]interface{}) bundle {
	bundle := make(map[string][]string)
	for key, content := range jsonContent {
		switch v := content.(type) {
		case string:
			bundle[key] = []string{v}
		case []interface{}:
			values := make([]string, 0)
			for _, value := range v {
				values = append(values, fmt.Sprintf("%v", value))
			}
			bundle[key] = values
		case map[string]interface{}:
			subValues := translator.mapBundleStructure(v)
			for subKey, subValue := range subValues {
				bundle[fmt.Sprintf("%s%s%s", key, keyDelim, subKey)] = subValue
			}
		default:
			bundle[key] = []string{fmt.Sprintf("%v", v)}
		}
	}

	return bundle
}
