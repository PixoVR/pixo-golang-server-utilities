package auth

type LicenceTypeEnum int

const (
	LicenceTypePixoBuilder LicenceTypeEnum = iota
)

func (e LicenceTypeEnum) String() string {
	return []string{
		"pixo-builder",
	}[e]
}

type TokenLicense struct {
	ID LicenceTypeEnum `json:"id"`
}
