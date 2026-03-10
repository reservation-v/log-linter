package rules

import "testing"

func TestCheckSymbols(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    bool
	}{
		{
			name:    "letters and spaces only",
			message: "server started",
			want:    true,
		},
		{
			name:    "letters digits and spaces",
			message: "server started on port 8080",
			want:    true,
		},
		{
			name:    "exclamation marks",
			message: "connection failed!!!",
			want:    false,
		},
		{
			name:    "colon",
			message: "warning: something went wrong",
			want:    false,
		},
		{
			name:    "ellipsis",
			message: "something went wrong...",
			want:    false,
		},
		{
			name:    "emoji",
			message: "server started 🚀",
			want:    false,
		},
		{
			name:    "underscore",
			message: "request_id received",
			want:    false,
		},
		{
			name:    "newline is allowed as whitespace",
			message: "server\nstarted",
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
			got := CheckSymbols(tt.message)
			if got != tt.want {
				t.Fatalf("CheckSymbols(%q) = %v, want %v", tt.message, got, tt.want)
			}
		})
	}
}
