package main

import (
	"testing"
)

func TestGetAlias(t *testing.T) {
	tests := []struct {
		name    string
		pkgName string
		version string
		want    string
	}{
		{name: "simple", pkgName: "@angular/core", version: "9.0.0", want: "@angular/core_9.0.0"},
		{name: "empty version", pkgName: "@angular/core", version: "", want: "@angular/core_"},
		{name: "empty name", pkgName: "", version: "1.0.0", want: "_1.0.0"},
		{name: "both empty", pkgName: "", version: "", want: "_"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetAlias(tt.pkgName, tt.version)
			if got != tt.want {
				t.Errorf("GetAlias(%q, %q) = %q, want %q", tt.pkgName, tt.version, got, tt.want)
			}
		})
	}
}

func TestPackageJSON_getAlias(t *testing.T) {
	p := &PackageJSON{Name: "@angular/core", Version: "9.0.0"}
	want := "@angular/core_9.0.0"
	if got := p.getAlias(); got != want {
		t.Errorf("getAlias() = %q, want %q", got, want)
	}
}

func TestClearDependenciesVersions(t *testing.T) {
	tests := []struct {
		name string
		in   map[string]string
		want map[string]string
	}{
		{
			name: "caret prefix stripped",
			in:   map[string]string{"@angular/core": "^9.0.0"},
			want: map[string]string{"@angular/core": "9.0.0"},
		},
		{
			name: "tilde prefix stripped",
			in:   map[string]string{"lodash": "~4.17.21"},
			want: map[string]string{"lodash": "4.17.21"},
		},
		{
			name: "range wildcard segment kept",
			in:   map[string]string{"foo": "1.x"},
			want: map[string]string{"foo": "1.x"},
		},
		{
			name: "bare wildcard has no match",
			in:   map[string]string{"foo": "*"},
			want: map[string]string{"foo": ""},
		},
		{
			name: "already clean version",
			in:   map[string]string{"foo": "1.2.3"},
			want: map[string]string{"foo": "1.2.3"},
		},
		{
			name: "empty map",
			in:   map[string]string{},
			want: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PackageJSON{
				Dependencies:     tt.in,
				DevDependencies:  tt.in,
				PeerDependencies: tt.in,
			}
			if err := p.clearDependenciesVersions(); err != nil {
				t.Fatalf("clearDependenciesVersions() error = %v", err)
			}
			for key, want := range tt.want {
				if got := p.Dependencies[key]; got != want {
					t.Errorf("Dependencies[%q] = %q, want %q", key, got, want)
				}
			}
		})
	}
}

func TestPackageJSON_Validate(t *testing.T) {
	tests := []struct {
		name    string
		pkg     PackageJSON
		wantErr bool
	}{
		{name: "name present", pkg: PackageJSON{Name: "my-package"}, wantErr: false},
		{name: "name empty", pkg: PackageJSON{Name: ""}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.pkg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_getRepositoryClientData(t *testing.T) {
	r := &Repository{
		URL: "https://github.com/owner/repo",
		PackageJSON: &PackageJSON{
			Name:    "repo",
			Version: "1.0.0",
		},
		Documentation: Documentation{Frontend: "https://example.com"},
		Notifications: []NotificationData{{Reason: "security_alert"}},
	}

	got := r.getRepositoryClientData()

	if got.Name != "repo" {
		t.Errorf("Name = %q, want %q", got.Name, "repo")
	}
	if got.Version != "1.0.0" {
		t.Errorf("Version = %q, want %q", got.Version, "1.0.0")
	}
	if got.URL != r.URL {
		t.Errorf("URL = %q, want %q", got.URL, r.URL)
	}
	if got.Documentation != r.Documentation {
		t.Errorf("Documentation = %+v, want %+v", got.Documentation, r.Documentation)
	}
	if len(got.Notifications) != 1 {
		t.Errorf("Notifications = %+v, want 1 entry", got.Notifications)
	}
	if got.Updates == nil || len(got.Updates) != 0 {
		t.Errorf("Updates = %+v, want empty non-nil slice", got.Updates)
	}
	if got.Tags == nil || len(got.Tags) != 0 {
		t.Errorf("Tags = %+v, want empty non-nil slice", got.Tags)
	}
}
