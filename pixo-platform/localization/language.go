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
	// LanguageCodeHi Hindi
	LanguageCodeHi LanguageCode = "hi"
	// LanguageCodeKn Kannada
	LanguageCodeKn LanguageCode = "kn"
	// LanguageCodeBn Bengali
	LanguageCodeBn LanguageCode = "bn"
	// LanguageCodePa Punjabi
	LanguageCodePa LanguageCode = "pa"
	// LanguageCodeMr Marathi
	LanguageCodeMr LanguageCode = "mr"
	// LanguageCodeMs Malay
	LanguageCodeMs LanguageCode = "ms"
	// LanguageCodeTr Turkish
	LanguageCodeTr LanguageCode = "tr"
	// LanguageCodeKo Korean
	LanguageCodeKo LanguageCode = "ko"
	// LanguageCodeRu Russian
	LanguageCodeRu LanguageCode = "ru"
	// LanguageCodeVi Vietnamese
	LanguageCodeVi LanguageCode = "vi"
	// LanguageCodeJv Javanese
	LanguageCodeJv LanguageCode = "jv"
	// LanguageCodeGu Gujarati
	LanguageCodeGu LanguageCode = "gu"
	// LanguageCodeNl Dutch
	LanguageCodeNl LanguageCode = "nl"
	// LanguageCodeSv Swedish
	LanguageCodeSv LanguageCode = "sv"
	// LanguageCodeCs Czech
	LanguageCodeCs LanguageCode = "cs"
	// LanguageCodePl Polish
	LanguageCodePl LanguageCode = "pl"
	// LanguageCodeLv Latvian
	LanguageCodeLv LanguageCode = "lv"
	// LanguageCodeSk Slovak
	LanguageCodeSk LanguageCode = "sk"
	// LanguageCodeEl Greek
	LanguageCodeEl LanguageCode = "el"
	// LanguageCodeLt Lithuanian
	LanguageCodeLt LanguageCode = "lt"
	// LanguageCodeUk Ukrainian
	LanguageCodeUk LanguageCode = "uk"
	// LanguageCodeFa Farsi
	LanguageCodeFa LanguageCode = "fa"
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
	LanguageCodeHi: {Language: "Hindi", LanguageCode: LanguageCodeHi, DisplayName: "हिन्दी"},
	LanguageCodeKn: {Language: "Kannada", LanguageCode: LanguageCodeKn, DisplayName: "ಕನ್ನಡ"},
	LanguageCodeBn: {Language: "Bengali", LanguageCode: LanguageCodeBn, DisplayName: "বাংলা"},
	LanguageCodePa: {Language: "Punjabi", LanguageCode: LanguageCodePa, DisplayName: "ਪੰਜਾਬੀ; پنجابی"},
	LanguageCodeMr: {Language: "Marathi", LanguageCode: LanguageCodeMr, DisplayName: "मराठी"},
	LanguageCodeMs: {Language: "Malay", LanguageCode: LanguageCodeMs, DisplayName: "മലയാളം"},
	LanguageCodeTr: {Language: "Turkish", LanguageCode: LanguageCodeTr, DisplayName: "Türkçe"},
	LanguageCodeKo: {Language: "Korean", LanguageCode: LanguageCodeKo, DisplayName: "한국어"},
	LanguageCodeRu: {Language: "Russian", LanguageCode: LanguageCodeRu, DisplayName: "Русский язык"},
	LanguageCodeVi: {Language: "Vietnamese", LanguageCode: LanguageCodeVi, DisplayName: "tiếng Việt"},
	LanguageCodeJv: {Language: "Javanese", LanguageCode: LanguageCodeJv, DisplayName: "ꦧꦱꦗꦮ"},
	LanguageCodeGu: {Language: "Gujarati", LanguageCode: LanguageCodeGu, DisplayName: "ગુજરાતી"},
	LanguageCodeNl: {Language: "Dutch", LanguageCode: LanguageCodeNl, DisplayName: "Nederlands"},
	LanguageCodeSv: {Language: "Swedish", LanguageCode: LanguageCodeSv, DisplayName: "Svenska"},
	LanguageCodeCs: {Language: "Czech", LanguageCode: LanguageCodeCs, DisplayName: "Čeština"},
	LanguageCodePl: {Language: "Polish", LanguageCode: LanguageCodePl, DisplayName: "Polski"},
	LanguageCodeLv: {Language: "Latvian", LanguageCode: LanguageCodeLv, DisplayName: "Latviski"},
	LanguageCodeSk: {Language: "Slovak", LanguageCode: LanguageCodeSk, DisplayName: "Slovenčina"},
	LanguageCodeEl: {Language: "Greek", LanguageCode: LanguageCodeEl, DisplayName: "Νέα Ελληνικά"},
	LanguageCodeLt: {Language: "Lithuanian", LanguageCode: LanguageCodeLt, DisplayName: "Lietuvių"},
	LanguageCodeUk: {Language: "Ukrainian", LanguageCode: LanguageCodeUk, DisplayName: "Українська"},
	LanguageCodeFa: {Language: "Farsi", LanguageCode: LanguageCodeFa, DisplayName: "فارسی"},
}

func (code LanguageCode) GetLanguage() (*Language, error) {
	if language, ok := languagesMap[code]; ok {
		return &language, nil
	}

	return nil, languageNotFoundError
}
