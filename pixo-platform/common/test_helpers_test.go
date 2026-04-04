package common_test

type Dummy interface {
	GetValue() int
}

type DummyStruct struct {
	Value int
}

func (d DummyStruct) GetValue() int {
	return d.Value
}
