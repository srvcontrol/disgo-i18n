package disgoi18n

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	translatorNominalCase1         = "translatorNominalCase1.json"
	translatorNominalCase2         = "translatorNominalCase2.json"
	translatorFailedUnmarshallCase = "translatorFailedUnmarshallCase.json"
	translatorFileDoesNotExistCase = "translatorFileDoesNotExistCase.json"

	content1 = `
	{
		"hi": ["this is a {{ .Test }}"],
		"with": ["all"],
		"the": ["elements", "we"],
		"can": ["find"],
		"in": ["a","json"],
		"config": ["file", "! {{ .Author }}"],
		"parse": ["{{if $foo}}{{end}}"]
	}
	`

	content2 = `
	{
		"this": ["is a {{ .Test }}"],
		"with.a.file": ["containing", "less", "variables"],
		"bye": ["see you"]
	}
	`

	badContent = `
	 [
		"content",
		"not",
		"ok",
		"test"
	 ]
	`
)

var (
	translatorTest *translatorImpl
)

func setUp() {
	translatorTest = newTranslatorImpl()
	if err := os.WriteFile(translatorNominalCase1, []byte(content1), os.ModePerm); err != nil {
		log.Fatalf("'%s' could not be created, test stopped: %v", translatorNominalCase1, err)
	}
	if err := os.WriteFile(translatorNominalCase2, []byte(content2), os.ModePerm); err != nil {
		log.Fatalf("'%s' could not be created, test stopped: %v", translatorNominalCase2, err)
	}
	if err := os.WriteFile(translatorFailedUnmarshallCase, []byte(badContent), os.ModePerm); err != nil {
		log.Fatalf("'%s' could not be created, test stopped: %v", translatorFailedUnmarshallCase, err)
	}
}

func tearDown() {
	translatorTest = nil
	if err := os.Remove(translatorNominalCase1); err != nil {
		log.Fatalf("'%s' could not be deleted: %v", translatorNominalCase1, err)
	}
	if err := os.Remove(translatorNominalCase2); err != nil {
		log.Fatalf("'%s' could not be deleted: %v", translatorNominalCase2, err)
	}
	if err := os.Remove(translatorFailedUnmarshallCase); err != nil {
		log.Fatalf("'%s' could not be deleted: %v", translatorFailedUnmarshallCase, err)
	}
}

func TestNew(t *testing.T) {
	setUp()
	defer tearDown()

	assert.Empty(t, translatorTest.translations)
	assert.Empty(t, translatorTest.loadedBundles)
}

func TestSetDefault(t *testing.T) {
	setUp()
	defer tearDown()

	assert.Equal(t, defaultLocale, translatorTest.defaultLocale)
	translatorTest.SetDefault(discord.LocaleItalian)
	assert.Equal(t, discord.LocaleItalian, translatorTest.defaultLocale)
}

func TestLoadBundle(t *testing.T) {
	setUp()
	defer tearDown()

	// Bad case, file does not exist
	_, err := os.Stat(translatorFileDoesNotExistCase)
	assert.Error(t, os.ErrNotExist, err)
	assert.Error(t, translatorTest.LoadBundle(discord.LocaleFrench, translatorFileDoesNotExistCase))
	assert.Empty(t, translatorTest.translations)
	assert.Empty(t, translatorTest.loadedBundles)

	// Bad case, file is not well structured
	assert.Error(t, translatorTest.LoadBundle(discord.LocaleFrench, translatorFailedUnmarshallCase))
	assert.Empty(t, translatorTest.translations)
	assert.Empty(t, translatorTest.loadedBundles)

	// Nominal case, load an existing and well structured bundle
	assert.NoError(t, translatorTest.LoadBundle(discord.LocaleFrench, translatorNominalCase1))
	assert.Equal(t, 1, len(translatorTest.loadedBundles))
	assert.Equal(t, 1, len(translatorTest.translations))
	assert.Equal(t, 7, len(translatorTest.translations[discord.LocaleFrench]))

	// Nominal case, reload a bundle
	assert.NoError(t, translatorTest.LoadBundle(discord.LocaleFrench, translatorNominalCase2))
	assert.Equal(t, 2, len(translatorTest.loadedBundles))
	assert.Equal(t, 1, len(translatorTest.translations))
	assert.Equal(t, 3, len(translatorTest.translations[discord.LocaleFrench]))

	// Nominal case, load a bundle already loaded but for another locale
	assert.NoError(t, translatorTest.LoadBundle(discord.LocaleEnglishGB, translatorNominalCase2))
	assert.Equal(t, 2, len(translatorTest.loadedBundles))
	assert.Equal(t, 2, len(translatorTest.translations))
	assert.Equal(t, 3, len(translatorTest.translations[discord.LocaleEnglishGB]))

	// Nominal case, reload a bundle linked to two locales
	assert.NoError(t, translatorTest.LoadBundle(discord.LocaleEnglishGB, translatorNominalCase1))
	assert.Equal(t, 2, len(translatorTest.loadedBundles))
	assert.Equal(t, 2, len(translatorTest.translations))
	assert.Equal(t, 7, len(translatorTest.translations[discord.LocaleEnglishGB]))
}

func TestGet(t *testing.T) {
	setUp()
	defer tearDown()

	// Nominal case, get without bundle
	assert.Equal(t, "hi", translatorTest.Get(discord.LocaleDutch, "hi", nil))

	// Nominal case, unexisting key with bundle loaded
	assert.NoError(t, translatorTest.LoadBundle(discord.LocaleDutch, translatorNominalCase1))
	assert.NoError(t, translatorTest.LoadBundle(defaultLocale, translatorNominalCase2))
	assert.Equal(t, "does_not_exist", translatorTest.Get(discord.LocaleDutch, "does_not_exist", nil))

	// Nominal case, Get existing key from loaded bundle
	assert.Equal(t, "this is a {{ .Test }}", translatorTest.Get(discord.LocaleDutch, "hi", nil))
	assert.Equal(t, "this is a test :)", translatorTest.Get(discord.LocaleDutch, "hi", Vars{"Test": "test :)"}))

	// Nominal case, Get key not present in loaded bundle but available in default
	assert.Equal(t, "see you", translatorTest.Get(discord.LocaleDutch, "bye", nil))

	// Bad case, value is not well structured to be parsed
	assert.Equal(t, "{{if $foo}}{{end}}", translatorTest.Get(discord.LocaleDutch, "parse", Vars{}))

	// Bad case, value is well structured but cannot inject value
	assert.Equal(t, "this is a {{ .Test }}", translatorTest.Get(discord.LocaleDutch, "hi", Vars{}))
}

func TestMapBundleStructure(t *testing.T) {
	setUp()
	defer tearDown()

	tests := []struct {
		Description    string
		Input          map[string]interface{}
		ExpectedBundle bundle
	}{
		{
			Description:    "Nil Input",
			Input:          nil,
			ExpectedBundle: make(bundle),
		},
		{
			Description:    "Empty Input",
			Input:          make(map[string]interface{}),
			ExpectedBundle: make(bundle),
		},
		{
			Description: "Simple string Input",
			Input: map[string]interface{}{
				"simple":       "translation",
				"variabilized": "translation {{ .translation }}",
			},
			ExpectedBundle: bundle{
				"simple":       []string{"translation"},
				"variabilized": []string{"translation {{ .translation }}"},
			},
		},
		{
			Description: "Different types handled",
			Input: map[string]interface{}{
				"pi":                                  3.14,
				"answer_to_ultimate_question_of_life": 42,
				"some_prime_numbers":                  []interface{}{2, "3", 5.0, 7},
			},
			ExpectedBundle: bundle{
				"pi":                                  []string{"3.14"},
				"answer_to_ultimate_question_of_life": []string{"42"},
				"some_prime_numbers":                  []string{"2", "3", "5", "7"},
			},
		},
		{
			Description: "Deep structure",
			Input: map[string]interface{}{
				"command": map[string]interface{}{
					"salutation": map[string]interface{}{
						"hi":  "Hello there!",
						"bye": []interface{}{"Bye {{ .anyone }}!", "See u {{ .anyone }}"},
					},
					"speak": map[string]interface{}{
						"random": []interface{}{"love to talk", "how are u?", "u're so interesting"},
					},
				},
				"panic": "I've panicked!",
			},
			ExpectedBundle: bundle{
				"command.salutation.hi":  []string{"Hello there!"},
				"command.salutation.bye": []string{"Bye {{ .anyone }}!", "See u {{ .anyone }}"},
				"command.speak.random":   []string{"love to talk", "how are u?", "u're so interesting"},
				"panic":                  []string{"I've panicked!"},
			},
		},
	}

	for _, test := range tests {
		bundle := translatorTest.mapBundleStructure(test.Input)
		assert.True(t, reflect.DeepEqual(test.ExpectedBundle, bundle),
			fmt.Sprintf("%s:\n\nExpecting: %v\n\nGot      : %v", test.Description, test.ExpectedBundle, bundle))
	}
}

func TestGetLocalizations(t *testing.T) {
	setUp()
	defer tearDown()

	// Nominal case: empty map when no bundle loaded
	assert.Empty(t, translatorTest.GetLocalizations("hi", Vars{}))

	// Nominal case: two bundles loaded so two translations expected
	assert.NoError(t, translatorTest.LoadBundle(discord.LocaleDutch, translatorNominalCase1))
	assert.NoError(t, translatorTest.LoadBundle(defaultLocale, translatorNominalCase2))
	assert.Equal(t, 2, len(translatorTest.GetLocalizations("hi", Vars{})))

	// Nominal case: three bundles loaded so three translations expected
	assert.NoError(t, translatorTest.LoadBundle(discord.LocaleChineseCN, translatorNominalCase1))
	assert.Equal(t, 3, len(translatorTest.GetLocalizations("hi", Vars{})))
}
