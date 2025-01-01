package sourcetool

import (
	"errors"
	"testing"

	"github.com/gofrs/uuid/v5"
)

func TestNew(t *testing.T) {
	apiKey := "test_api_key"
	st := New(apiKey)

	if st == nil {
		t.Fatal("New returned nil")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"APIKey", st.apiKey, apiKey},
		{"Subdomain", st.subdomain, "test"},
		{"Endpoint", st.endpoint, "ws://test.local.trysourcetool.com:8080/ws"},
		{"Navigations length", len(st.navigations), 0},
		{"Pages length", len(st.pages), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestPage(t *testing.T) {
	st := New("test_api_key")
	pageName := "TestPage"
	pageHandler := func(ui UIBuilder) error { return nil }

	// Add page
	st.Page(pageName, pageHandler)

	// Verify page was added
	pageID := st.generatePageID(pageName)
	page, exists := st.pages[pageID]
	if !exists {
		t.Fatal("Page was not added to pages map")
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Page ID", page.id, pageID},
		{"Page name", page.name, pageName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}

	// Test handler execution
	if err := page.run(nil); err != nil {
		t.Errorf("Page handler returned unexpected error: %v", err)
	}

	// Test handler with error
	errorHandler := func(ui UIBuilder) error {
		return errors.New("test error")
	}
	st.Page("ErrorPage", errorHandler)
	errorPageID := st.generatePageID("ErrorPage")
	errorPage := st.pages[errorPageID]
	if err := errorPage.run(nil); err == nil {
		t.Error("Expected error from handler, got nil")
	}
}

func TestGeneratePageID(t *testing.T) {
	st := New("test_api_key")
	pageName := "TestPage"

	// Test ID generation is deterministic
	id1 := st.generatePageID(pageName)
	id2 := st.generatePageID(pageName)

	if id1 == uuid.Nil {
		t.Error("Generated ID is nil")
	}

	if id1 != id2 {
		t.Error("Generated IDs are not deterministic")
	}

	// Test different page names generate different IDs
	otherID := st.generatePageID("OtherPage")
	if id1 == otherID {
		t.Error("Different page names generated same ID")
	}

	// Test different subdomains generate different IDs
	otherST := New("other_api_key")
	otherSTID := otherST.generatePageID(pageName)
	if id1 == otherSTID {
		t.Error("Different subdomains generated same ID")
	}
}

func TestPageManager(t *testing.T) {
	pages := make(map[uuid.UUID]*page)
	pageID := uuid.Must(uuid.NewV4())
	testPage := &page{
		id:   pageID,
		name: "TestPage",
	}
	pages[pageID] = testPage

	pm := newPageManager(pages)

	// Test getting existing page
	got := pm.getPage(pageID)
	if got != testPage {
		t.Error("getPage returned wrong page")
	}

	// Test getting non-existent page
	nonExistentID := uuid.Must(uuid.NewV4())
	got = pm.getPage(nonExistentID)
	if got != nil {
		t.Error("getPage returned non-nil for non-existent page")
	}
}
