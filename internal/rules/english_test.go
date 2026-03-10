package rules

import "testing"

func TestCheckEnglish(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    bool
	}{
		{
			name:    "english words",
			message: "starting server",
			want:    true,
		},
		{
			name:    "english with digits",
			message: "server started on port 8080",
			want:    true,
		},
		{
			name:    "english with punctuation",
			message: "server started!",
			want:    true,
		},
		{
			name:    "russian text",
			message: "запуск сервера",
			want:    false,
		},
		{
			name:    "mixed english and russian",
			message: "starting сервер",
			want:    false,
		},
		{
			name:    "accented latin",
			message: "cafe déjà vu",
			want:    false,
		},
		{
			name:    "numbers only",
			message: "8080",
			want:    true,
		},
		{
			name:    "empty message",
			message: "",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckEnglish(tt.message)
			if got != tt.want {
				t.Fatalf("CheckEnglish(%q) = %v, want %v", tt.message, got, tt.want)
			}
		})
	}
}
