package sourcetool

import (
	"testing"
)

func TestJoinPath(t *testing.T) {
	tests := []struct {
		name     string
		basePath string
		path     string
		want     string
	}{
		{
			name:     "Empty base path",
			basePath: "",
			path:     "/users",
			want:     "/users",
		},
		{
			name:     "Base path with trailing slash",
			basePath: "/admin/",
			path:     "users",
			want:     "/admin/users",
		},
		{
			name:     "Path without leading slash",
			basePath: "/admin",
			path:     "users",
			want:     "/admin/users",
		},
		{
			name:     "Both with slashes",
			basePath: "/admin/",
			path:     "/users/",
			want:     "/admin/users/",
		},
		{
			name:     "Nested paths",
			basePath: "/api/v1",
			path:     "users/list",
			want:     "/api/v1/users/list",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &router{basePath: tt.basePath}
			got := r.joinPath(tt.path)
			if got != tt.want {
				t.Errorf("joinPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name   string
		groups []string
		want   []string
	}{
		{
			name:   "No duplicates",
			groups: []string{"admin", "user", "guest"},
			want:   []string{"admin", "user", "guest"},
		},
		{
			name:   "With duplicates",
			groups: []string{"admin", "user", "admin", "guest", "user"},
			want:   []string{"admin", "user", "guest"},
		},
		{
			name:   "Empty slice",
			groups: []string{},
			want:   []string{},
		},
		{
			name:   "Single element",
			groups: []string{"admin"},
			want:   []string{"admin"},
		},
		{
			name:   "All duplicates",
			groups: []string{"admin", "admin", "admin"},
			want:   []string{"admin"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeDuplicates(tt.groups)
			if len(got) != len(tt.want) {
				t.Errorf("removeDuplicates() length = %v, want %v", len(got), len(tt.want))
			}
			for _, w := range tt.want {
				found := false
				for _, g := range got {
					if g == w {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("removeDuplicates() missing element %v", w)
				}
			}
		})
	}
}

func TestGeneratePageID(t *testing.T) {
	r := &router{
		namespaceDNS: "test.trysourcetool.com",
	}

	tests := []struct {
		name     string
		path     string
		wantSame bool
	}{
		{
			name:     "Simple path",
			path:     "/users",
			wantSame: true,
		},
		{
			name:     "Nested path",
			path:     "/admin/users/list",
			wantSame: true,
		},
		{
			name:     "Root path",
			path:     "/",
			wantSame: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id1 := r.generatePageID(tt.path)
			id2 := r.generatePageID(tt.path)

			if tt.wantSame && id1 != id2 {
				t.Errorf("generatePageID() generated different IDs for same path")
			}

			differentPath := tt.path + "/different"
			id3 := r.generatePageID(differentPath)
			if id1 == id3 {
				t.Errorf("generatePageID() generated same ID for different paths")
			}
		})
	}
}

func TestRouterAccessGroups(t *testing.T) {
	pageHandler := func(ui UIBuilder) error { return nil }

	t.Run("Router level access groups", func(t *testing.T) {
		st := New("test_api_key")
		users := st.Group("/users")
		users.AccessGroups("admin")
		users.Page("/list", "List users", pageHandler)

		page := findPageByPath(st.pages, "/users/list")
		if page == nil {
			t.Fatal("Page not found")
		}

		if len(page.accessGroups) != 1 || page.accessGroups[0] != "admin" {
			t.Errorf("Expected [admin] access group, got %v", page.accessGroups)
		}
	})

	t.Run("Page specific access groups", func(t *testing.T) {
		st := New("test_api_key")
		users := st.Group("/users")
		users.AccessGroups("admin")
		users.AccessGroups("customer_support").Page("/create", "Create user", pageHandler)

		page := findPageByPath(st.pages, "/users/create")
		if page == nil {
			t.Fatal("Page not found")
		}

		expectedGroups := []string{"admin", "customer_support"}
		if len(page.accessGroups) != len(expectedGroups) {
			t.Errorf("Expected %v access groups, got %v", expectedGroups, page.accessGroups)
		}

		for _, group := range expectedGroups {
			found := false
			for _, actualGroup := range page.accessGroups {
				if actualGroup == group {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected group %s not found in %v", group, page.accessGroups)
			}
		}
	})

	t.Run("Nested groups inheritance", func(t *testing.T) {
		st := New("test_api_key")
		users := st.Group("/users")
		users.AccessGroups("admin")
		products := users.Group("/products")
		products.AccessGroups("product_manager")
		products.Page("/list", "List products", pageHandler)

		page := findPageByPath(st.pages, "/users/products/list")
		if page == nil {
			t.Fatal("Page not found")
		}

		expectedGroups := []string{"admin", "product_manager"}
		if len(page.accessGroups) != len(expectedGroups) {
			t.Errorf("Expected %v access groups, got %v", expectedGroups, page.accessGroups)
		}

		for _, group := range expectedGroups {
			found := false
			for _, actualGroup := range page.accessGroups {
				if actualGroup == group {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected group %s not found in %v", group, page.accessGroups)
			}
		}
	})

	t.Run("Page specific groups do not affect siblings", func(t *testing.T) {
		st := New("test_api_key")
		users := st.Group("/users")
		users.AccessGroups("admin")
		users.Page("/list", "List users", pageHandler)
		users.AccessGroups("customer_support").Page("/create", "Create user", pageHandler)

		listPage := findPageByPath(st.pages, "/users/list")
		if listPage == nil {
			t.Fatal("List page not found")
		}

		if len(listPage.accessGroups) != 1 || listPage.accessGroups[0] != "admin" {
			t.Errorf("Expected [admin] access group for list page, got %v", listPage.accessGroups)
		}

		createPage := findPageByPath(st.pages, "/users/create")
		if createPage == nil {
			t.Fatal("Create page not found")
		}

		expectedGroups := []string{"admin", "customer_support"}
		if len(createPage.accessGroups) != len(expectedGroups) {
			t.Errorf("Expected %v access groups for create page, got %v", expectedGroups, createPage.accessGroups)
		}
	})
}

func TestRouterGroup(t *testing.T) {
	pageHandler := func(ui UIBuilder) error { return nil }

	t.Run("Base path construction", func(t *testing.T) {
		st := New("test_api_key")
		admin := st.Group("/admin")
		settings := admin.Group("/settings")
		settings.Page("/users", "User Settings", pageHandler)

		page := findPageByPath(st.pages, "/admin/settings/users")
		if page == nil {
			t.Fatal("Page not found")
		}

		if page.path != "/admin/settings/users" {
			t.Errorf("Expected path /admin/settings/users, got %s", page.path)
		}
	})

	t.Run("Multiple nested groups", func(t *testing.T) {
		st := New("test_api_key")
		api := st.Group("/api")
		v1 := api.Group("/v1")
		users := v1.Group("/users")
		users.Page("/list", "List Users", pageHandler)

		page := findPageByPath(st.pages, "/api/v1/users/list")
		if page == nil {
			t.Fatal("Page not found")
		}

		if page.path != "/api/v1/users/list" {
			t.Errorf("Expected path /api/v1/users/list, got %s", page.path)
		}
	})

	t.Run("Group access group inheritance", func(t *testing.T) {
		st := New("test_api_key")
		api := st.Group("/api")
		api.AccessGroups("api_user")
		v1 := api.Group("/v1")
		v1.AccessGroups("v1_user")
		v1.Page("/users", "Users API", pageHandler)

		page := findPageByPath(st.pages, "/api/v1/users")
		if page == nil {
			t.Fatal("Page not found")
		}

		expectedGroups := []string{"api_user", "v1_user"}
		if len(page.accessGroups) != len(expectedGroups) {
			t.Errorf("Expected %v access groups, got %v", expectedGroups, page.accessGroups)
		}

		for _, group := range expectedGroups {
			found := false
			for _, actualGroup := range page.accessGroups {
				if actualGroup == group {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected group %s not found in %v", group, page.accessGroups)
			}
		}
	})
}
