package rules

import "testing"

func TestCheckLowercase(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    bool
	}{
		{
			name:    "lowercase ascii",
			message: "starting server",
			want:    true,
		},
		{
			name:    "uppercase ascii",
			message: "Starting server",
			want:    false,
		},
		{
			name:    "lowercase unicode",
			message: "запуск сервера",
			want:    true,
		},
		{
			name:    "starts with digit",
			message: "8080 started",
			want:    false,
		},
		{
			name:    "starts with punctuation",
			message: ".started",
			want:    false,
		},
		{
			name:    "empty message",
			message: "",
			want:    true,
		},
		{
			name:    "raw lowercase",
			message: "server\nstarted",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckLowercase(tt.message)
			if got != tt.want {
				t.Fatalf("CheckLowercase(%q) = %v, want %v", tt.message, got, tt.want)
			}
		})
	}
}
