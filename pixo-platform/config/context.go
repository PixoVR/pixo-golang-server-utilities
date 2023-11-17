package config

type ContextRequest string

func (c ContextRequest) String() string {
	return string(c)
}
