package form

type Options struct {
	ButtonLabel    string
	ButtonDisabled bool
	ClearOnSubmit  bool
}

type Option func(*Options)
