package form

type Options struct {
	ClearOnSubmit bool
}

type Option func(*Options)
