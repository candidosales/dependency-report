package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-github/v90/github"
)

func TestGetTypeDependency(t *testing.T) {
	tests := []struct {
		name     string
		quantity int
		want     string
	}{
		{name: "zero", quantity: 0, want: "good"},
		{name: "lower boundary good", quantity: 1, want: "good"},
		{name: "upper boundary good", quantity: 2, want: "good"},
		{name: "lower boundary warning", quantity: 3, want: "warning"},
		{name: "upper boundary warning", quantity: 5, want: "warning"},
		{name: "lower boundary bad", quantity: 6, want: "bad"},
		{name: "upper boundary bad", quantity: 10, want: "bad"},
		{name: "lower boundary terrible", quantity: 11, want: "terrible"},
		{name: "large quantity terrible", quantity: 100, want: "terrible"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTypeDependency(tt.quantity); got != tt.want {
				t.Errorf("getTypeDependency(%d) = %q, want %q", tt.quantity, got, tt.want)
			}
		})
	}
}

func TestSplitRepositoryURL(t *testing.T) {
	tests := []struct {
		name    string
		repo    Repository
		want    map[string]string
		wantErr bool
	}{
		{
			name: "project url",
			repo: Repository{URL: "https://github.com/vendasta/galaxy", Type: TypeProject},
			want: map[string]string{
				"owner":       "vendasta",
				"repo":        "galaxy",
				"packageJSON": "package.json",
			},
		},
		{
			name: "component url",
			repo: Repository{
				URL:  "https://github.com/vendasta/frontend/tree/master/angular/projects/billing-ui",
				Type: TypeComponent,
			},
			want: map[string]string{
				"owner":       "vendasta",
				"repo":        "frontend",
				"packageJSON": "/angular/projects/billing-ui/package.json",
			},
		},
		{
			name:    "malformed project url",
			repo:    Repository{URL: "https://github.com", Type: TypeProject},
			wantErr: true,
		},
		{
			name:    "component url too short",
			repo:    Repository{URL: "https://github.com/vendasta/repo", Type: TypeComponent},
			wantErr: true,
		},
		{
			name:    "empty url",
			repo:    Repository{URL: "", Type: TypeProject},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := splitRepositoryURL(tt.repo)
			if (err != nil) != tt.wantErr {
				t.Fatalf("splitRepositoryURL() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitRepositoryURL() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestReadConfig(t *testing.T) {
	t.Run("valid file", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "config.json")
		content := `{"filters":["@angular/core_9"],"repositories":[{"url":"https://github.com/owner/repo","type":"project"}]}`
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test config: %v", err)
		}

		config, err := readConfig(path)
		if err != nil {
			t.Fatalf("readConfig() error = %v", err)
		}
		if len(config.Filters) != 1 || config.Filters[0] != "@angular/core_9" {
			t.Errorf("Filters = %v, want [@angular/core_9]", config.Filters)
		}
		if len(config.Repositories) != 1 || config.Repositories[0].URL != "https://github.com/owner/repo" {
			t.Errorf("Repositories = %+v, want one repo with the expected URL", config.Repositories)
		}
	})

	t.Run("missing file", func(t *testing.T) {
		_, err := readConfig(filepath.Join(t.TempDir(), "does-not-exist.json"))
		if err == nil {
			t.Error("readConfig() error = nil, want error for missing file")
		}
	})

	t.Run("malformed json", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "config.json")
		if err := os.WriteFile(path, []byte("{not valid json"), 0644); err != nil {
			t.Fatalf("failed to write test config: %v", err)
		}

		_, err := readConfig(path)
		if err == nil {
			t.Error("readConfig() error = nil, want error for malformed json")
		}
	})
}

func TestIsInconsistent(t *testing.T) {
	tests := []struct {
		name string
		repo RepositoryClientData
		want bool
	}{
		{name: "nil updates", repo: RepositoryClientData{Updates: nil}, want: false},
		{name: "empty updates", repo: RepositoryClientData{Updates: []*UpdateComponent{}}, want: false},
		{name: "with updates", repo: RepositoryClientData{Updates: []*UpdateComponent{{Name: "foo"}}}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isInconsistent(&tt.repo); got != tt.want {
				t.Errorf("isInconsistent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsVulnerable(t *testing.T) {
	tests := []struct {
		name string
		repo RepositoryClientData
		want bool
	}{
		{name: "nil notifications", repo: RepositoryClientData{Notifications: nil}, want: false},
		{name: "empty notifications", repo: RepositoryClientData{Notifications: []NotificationData{}}, want: false},
		{name: "with notifications", repo: RepositoryClientData{Notifications: []NotificationData{{Reason: "security_alert"}}}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isVulnerable(&tt.repo); got != tt.want {
				t.Errorf("isVulnerable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateSummary(t *testing.T) {
	app := &AppConfig{}

	projectsClientData := []RepositoryClientData{
		{Name: "clean"},
		{Name: "inconsistent", Updates: []*UpdateComponent{{Name: "foo"}}},
		{Name: "vulnerable", Notifications: []NotificationData{{Reason: "security_alert"}}},
		{
			Name:          "both",
			Updates:       []*UpdateComponent{{Name: "bar"}},
			Notifications: []NotificationData{{Reason: "security_alert"}},
		},
	}

	summary := app.generateSummary(projectsClientData)

	if !reflect.DeepEqual(summary.Updated, []string{"clean"}) {
		t.Errorf("Updated = %v, want [clean]", summary.Updated)
	}
	if !reflect.DeepEqual(summary.Inconsistent, []string{"inconsistent", "both"}) {
		t.Errorf("Inconsistent = %v, want [inconsistent both]", summary.Inconsistent)
	}
	if !reflect.DeepEqual(summary.Vulnerable, []string{"vulnerable", "both"}) {
		t.Errorf("Vulnerable = %v, want [vulnerable both]", summary.Vulnerable)
	}
}

func TestToNotificationData(t *testing.T) {
	updatedAt := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)

	tests := []struct {
		name string
		in   *github.Notification
		want NotificationData
	}{
		{
			name: "all fields populated",
			in: &github.Notification{
				Reason:    github.Ptr("security_alert"),
				URL:       github.Ptr("https://api.github.com/notifications/threads/1"),
				UpdatedAt: &github.Timestamp{Time: updatedAt},
				Subject:   &github.NotificationSubject{Title: github.Ptr("Bump lodash")},
			},
			want: NotificationData{
				Reason:    "security_alert",
				Title:     "Bump lodash",
				URL:       "https://api.github.com/notifications/threads/1",
				UpdatedAt: updatedAt,
			},
		},
		{
			name: "all fields nil",
			in:   &github.Notification{},
			want: NotificationData{},
		},
		{
			name: "subject present but title nil",
			in:   &github.Notification{Subject: &github.NotificationSubject{}},
			want: NotificationData{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toNotificationData(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toNotificationData() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestStatsCountDependenciesByVersions(t *testing.T) {
	app := &AppConfig{}

	projects := []Repository{
		{PackageJSON: &PackageJSON{Name: "a", Dependencies: map[string]string{"lodash": "4.17.21"}}},
		{PackageJSON: &PackageJSON{Name: "b", Dependencies: map[string]string{"lodash": "4.17.21"}}},
		{PackageJSON: &PackageJSON{Name: "c", Dependencies: map[string]string{"lodash": "3.10.1"}}},
	}

	got := app.statsCountDependenciesByVersions(projects)

	lodash, ok := got["lodash"]
	if !ok {
		t.Fatalf("expected lodash entry, got %+v", got)
	}
	if lodash.Type != "good" {
		t.Errorf("Type = %q, want good", lodash.Type)
	}
	if lodash.Versions["4.17.21"].Quantity != 2 {
		t.Errorf("4.17.21 Quantity = %d, want 2", lodash.Versions["4.17.21"].Quantity)
	}
	if !reflect.DeepEqual(lodash.Versions["4.17.21"].Projects, []string{"a", "b"}) {
		t.Errorf("4.17.21 Projects = %v, want [a b]", lodash.Versions["4.17.21"].Projects)
	}
	if lodash.Versions["3.10.1"].Quantity != 1 {
		t.Errorf("3.10.1 Quantity = %d, want 1", lodash.Versions["3.10.1"].Quantity)
	}
}

func TestStatsCountProjectsByFilters(t *testing.T) {
	app := &AppConfig{
		config: Config{Filters: []string{"@angular/core_9", "@angular/core_10"}},
	}

	projects := []Repository{
		{PackageJSON: &PackageJSON{Name: "a", Dependencies: map[string]string{"@angular/core": "9.0.0"}}},
		{PackageJSON: &PackageJSON{Name: "b", Dependencies: map[string]string{"@angular/core": "9.0.0"}}},
		{PackageJSON: &PackageJSON{Name: "c", Dependencies: map[string]string{"@angular/core": "8.0.0"}}},
	}
	projectsClientData := []RepositoryClientData{{Name: "a"}, {Name: "b"}, {Name: "c"}}

	got := app.statsCountProjectsByFilters(projects, projectsClientData)

	want := []FilterCount{
		{Filter: "@angular/core_9", Count: 2},
		{Filter: "@angular/core_10", Count: 0},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("statsCountProjectsByFilters() = %+v, want %+v", got, want)
	}
	if projectsClientData[0].Filter != "@angular/core_9" || projectsClientData[1].Filter != "@angular/core_9" {
		t.Errorf("expected matched projects to have Filter set, got %+v", projectsClientData)
	}
	if projectsClientData[2].Filter != "" {
		t.Errorf("expected unmatched project to have empty Filter, got %q", projectsClientData[2].Filter)
	}
}

func TestStatsCountComponentsByFilters(t *testing.T) {
	app := &AppConfig{
		config: Config{Filters: []string{"@angular/core_9"}},
	}

	components := []Repository{
		{PackageJSON: &PackageJSON{Name: "ui", PeerDependencies: map[string]string{"@angular/core": "9.0.0"}}},
	}
	componentsClientData := []RepositoryClientData{{Name: "ui"}}

	got := app.statsCountComponentsByFilters(components, componentsClientData)

	want := []FilterCount{{Filter: "@angular/core_9", Count: 1}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("statsCountComponentsByFilters() = %+v, want %+v", got, want)
	}
	if componentsClientData[0].Filter != "@angular/core_9" {
		t.Errorf("expected component Filter set, got %q", componentsClientData[0].Filter)
	}
}

func TestStatusProjectsByComponents(t *testing.T) {
	app := &AppConfig{}

	projects := []Repository{
		{PackageJSON: &PackageJSON{Name: "app", Dependencies: map[string]string{"ui-lib": "1.0.0"}}},
	}
	projectsClientData := []RepositoryClientData{{Name: "app"}}
	components := []RepositoryClientData{{Name: "ui-lib", Version: "2.0.0"}}

	app.statusProjectsByComponents(projects, projectsClientData, components)

	if len(projectsClientData[0].Updates) != 1 {
		t.Fatalf("Updates = %+v, want 1 entry", projectsClientData[0].Updates)
	}
	update := projectsClientData[0].Updates[0]
	if update.Name != "ui-lib" || update.Current != "1.0.0" || update.Update != "2.0.0" {
		t.Errorf("Updates[0] = %+v, want {ui-lib 1.0.0 2.0.0}", update)
	}
}

func TestStatusProjectsByComponents_NoMismatch(t *testing.T) {
	app := &AppConfig{}

	projects := []Repository{
		{PackageJSON: &PackageJSON{Name: "app", Dependencies: map[string]string{"ui-lib": "2.0.0"}}},
	}
	projectsClientData := []RepositoryClientData{{Name: "app"}}
	components := []RepositoryClientData{{Name: "ui-lib", Version: "2.0.0"}}

	app.statusProjectsByComponents(projects, projectsClientData, components)

	if len(projectsClientData[0].Updates) != 0 {
		t.Errorf("Updates = %+v, want no entries when versions match", projectsClientData[0].Updates)
	}
}

func TestSplitProjectsComponents(t *testing.T) {
	app := &AppConfig{
		config: Config{
			Repositories: []Repository{
				{Type: TypeProject, PackageJSON: &PackageJSON{Name: "p1"}},
				{Type: TypeComponent, PackageJSON: &PackageJSON{Name: "c1"}},
				{Type: TypeProject, PackageJSON: &PackageJSON{Name: "p2"}},
			},
		},
	}

	projects, components, projectsClientData, componentsClientData := app.splitProjectsComponents()

	if len(*projects) != 2 {
		t.Errorf("projects = %+v, want 2 entries", *projects)
	}
	if len(*components) != 1 {
		t.Errorf("components = %+v, want 1 entry", *components)
	}
	if len(*projectsClientData) != 2 {
		t.Errorf("projectsClientData = %+v, want 2 entries", *projectsClientData)
	}
	if len(*componentsClientData) != 1 {
		t.Errorf("componentsClientData = %+v, want 1 entry", *componentsClientData)
	}
}
