package ports

import "testing"

func TestMaskAPIKey_VariousInputs(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want string
	}{
		{"empty string", "", ""},
		{"1 char", "a", "*"},
		{"2 chars", "ab", "**"},
		{"3 chars", "abc", "***"},
		{"exactly 4 chars", "abcd", "abcd"},
		{"5 chars", "12345", "*2345"},
		{"8 chars shows last 4", "sk-12345", "****2345"},
		{"long key masks properly", "sk-or-abcdefghijklmnop", "******************mnop"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaskAPIKey(tt.key)
			if got != tt.want {
				t.Errorf("MaskAPIKey(%q) = %q, want %q", tt.key, got, tt.want)
			}
			// Verify length is preserved for non-empty strings
			if tt.key != "" && len(got) != len(tt.key) {
				t.Errorf("MaskAPIKey(%q) length = %d, want %d", tt.key, len(got), len(tt.key))
			}
		})
	}
}
