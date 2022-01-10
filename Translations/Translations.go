package Translations

import (
	"encoding/json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var loc *i18n.Localizer

func init() {
	bundle := i18n.NewBundle(language.English)
	loc = i18n.NewLocalizer(bundle, language.English.String())

	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustLoadMessageFile("Translations/en.json")
	bundle.MustLoadMessageFile("Translations/nl.json")
}

func SLocalize(message *i18n.Message) string {
	return loc.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: message})
}

func Localize(message *i18n.Message, parameters map[string]interface{}) string {
	var pluralCount = 0
	if count, ok := parameters["count"]; ok {
		if count, ok := count.(int); ok {
			pluralCount = count
		}
	}

	return loc.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: message,
		TemplateData:   parameters,
		PluralCount:    pluralCount,
	})
}
