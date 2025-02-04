package localization

import "fmt"

var languageNotFoundError = fmt.Errorf("language not found")

type Language struct {
	Language     string       `json:"language"`
	LanguageCode LanguageCode `json:"languageCode"`
	DisplayName  string       `json:"displayName"`
}

type LanguageCode string

const (
	// LanguageCodeAr Arabic
	LanguageCodeAr LanguageCode = "ar"
	// LanguageCodeZh Chinese
	LanguageCodeZh LanguageCode = "zh"
	// LanguageCodeEn English
	LanguageCodeEn LanguageCode = "en"
	// LanguageCodeFr French
	LanguageCodeFr LanguageCode = "fr"
	// LanguageCodeDe German
	LanguageCodeDe LanguageCode = "de"
	// LanguageCodeIt Italian
	LanguageCodeIt LanguageCode = "it"
	// LanguageCodeGa Irish
	LanguageCodeGa LanguageCode = "ga"
	// LanguageCodeJa Japanese
	LanguageCodeJa LanguageCode = "ja"
	// LanguageCodePt Portuguese
	LanguageCodePt LanguageCode = "pt"
	// LanguageCodeEs Spanish
	LanguageCodeEs LanguageCode = "es"
)

func (code LanguageCode) IsValid() bool {
	_, ok := languagesMap[code]
	return ok
}

const BaseLanguageCode = LanguageCodeEn

var languagesMap = map[LanguageCode]Language{
	LanguageCodeEn: {Language: "English", LanguageCode: LanguageCodeEn, DisplayName: "English"},
	LanguageCodeAr: {Language: "Arabic", LanguageCode: LanguageCodeAr, DisplayName: "العربية"},
	LanguageCodeZh: {Language: "Chinese", LanguageCode: LanguageCodeZh, DisplayName: "中文"},
	LanguageCodeFr: {Language: "French", LanguageCode: LanguageCodeFr, DisplayName: "Français"},
	LanguageCodeDe: {Language: "German", LanguageCode: LanguageCodeDe, DisplayName: "Deutsch"},
	LanguageCodeIt: {Language: "Italian", LanguageCode: LanguageCodeIt, DisplayName: "Italiano"},
	LanguageCodeGa: {Language: "Irish", LanguageCode: LanguageCodeGa, DisplayName: "Gaeilge"},
	LanguageCodeJa: {Language: "Japanese", LanguageCode: LanguageCodeJa, DisplayName: "日本語"},
	LanguageCodePt: {Language: "Portuguese", LanguageCode: LanguageCodePt, DisplayName: "Português"},
	LanguageCodeEs: {Language: "Spanish", LanguageCode: LanguageCodeEs, DisplayName: "Español"},
}

func (code LanguageCode) GetLanguage() (*Language, error) {
	if language, ok := languagesMap[code]; ok {
		return &language, nil
	}

	return nil, languageNotFoundError
}
