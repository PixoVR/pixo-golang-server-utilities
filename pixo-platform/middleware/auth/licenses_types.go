package auth

const LicenseTypePixoBuilder = "pixo-builder"

type TokenLicense struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}
