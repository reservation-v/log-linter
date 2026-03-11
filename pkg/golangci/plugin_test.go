package golangci

import "testing"

func TestSettingsConfig(t *testing.T) {
	settings := Settings{
		Disable: []string{"english", "sensitive"},
	}

	cfg := settings.config()
	if !cfg.Lowercase {
		t.Fatal("Lowercase = false, want true")
	}
	if cfg.English {
		t.Fatal("English = true, want false")
	}
	if !cfg.Symbols {
		t.Fatal("Symbols = false, want true")
	}
	if cfg.Sensitive {
		t.Fatal("Sensitive = true, want false")
	}
}

func TestNewRejectsUnknownDisabledCheck(t *testing.T) {
	_, err := New(map[string]any{
		"disable": []string{"unknown"},
	})
	if err == nil {
		t.Fatal("New() error = nil, want non-nil")
	}
}

func TestNewAcceptsNestedSettings(t *testing.T) {
	plugin, err := New(map[string]any{
		"type":        "module",
		"description": "test",
		"settings": map[string]any{
			"disable": []string{"english"},
		},
	})
	if err != nil {
		t.Fatalf("New() error = %v, want nil", err)
	}

	instance, ok := plugin.(*Plugin)
	if !ok {
		t.Fatalf("New() returned %T, want *Plugin", plugin)
	}
	if len(instance.settings.Disable) != 1 || instance.settings.Disable[0] != "english" {
		t.Fatalf("settings.Disable = %v, want [english]", instance.settings.Disable)
	}
}
