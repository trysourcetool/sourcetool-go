package pagerunner

import "github.com/trysourcetool/sourcetool-go/ui"

type PageRunner struct {
	pages []*ui.Page
}

func New(pages []*ui.Page) *PageRunner {
	return &PageRunner{
		pages: pages,
	}
}
