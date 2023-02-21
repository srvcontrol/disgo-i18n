package disgoi18n

import (
	"github.com/disgoorg/disgo/discord"
	"log"
)

func newMock() *translatorMock {
	return &translatorMock{}
}

func (mock *translatorMock) SetDefault(locale discord.Locale) {
	if mock.SetDefaultFunc != nil {
		mock.SetDefaultFunc(locale)
		return
	}

	log.Println("SetDefault not mocked")
}

func (mock *translatorMock) LoadBundle(locale discord.Locale, file string) error {
	if mock.LoadBundleFunc != nil {
		return mock.LoadBundleFunc(locale, file)
	}

	log.Println("LoadBundle not mocked")
	return nil
}

func (mock *translatorMock) Get(locale discord.Locale, key string, variables Vars) string {
	if mock.GetFunc != nil {
		return mock.GetFunc(locale, key, variables)
	}

	log.Println("Get not mocked")
	return ""
}

func (mock *translatorMock) GetLocalizations(key string, variables Vars) map[discord.Locale]string {
	if mock.GetFunc != nil {
		return mock.GetLocalizationsFunc(key, variables)
	}

	log.Println("GetLocalizations not mocked")
	return nil
}
