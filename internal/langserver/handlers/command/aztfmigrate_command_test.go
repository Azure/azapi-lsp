package command

import (
	"testing"
)

func Test_GetWorkingDirectory(t *testing.T) {
	tests := []struct {
		name  string
		os    string
		input string
		want  string
	}{
		{
			name:  "basic on darwin",
			os:    "darwin",
			input: "file:///Users/username/go/src/github.com/Azure/azapi-lsp/main.tf",
			want:  "/Users/username/go/src/github.com/Azure/azapi-lsp",
		},

		{
			name:  "basic on windows",
			os:    "windows",
			input: "file:///c%3A/Users/username/go/src/github.com/Azure/azapi-lsp/main.tf",
			want:  "c:\\Users\\username\\go\\src\\github.com\\Azure\\azapi-lsp",
		},

		{
			name:  "path with spaces on windows",
			os:    "windows",
			input: "file:///c%3A/Users/username/go/src/github.com/Azure/azapi-lsp/terraform%20files/main.tf",
			want:  "c:\\Users\\username\\go\\src\\github.com\\Azure\\azapi-lsp\\terraform files",
		},

		{
			name:  "path with Chinese characters on windows",
			os:    "windows",
			input: "file:///c%3A/Users/username/go/src/github.com/Azure/azapi-lsp/%E4%B8%AD%E6%96%87/main.tf",
			want:  "c:\\Users\\username\\go\\src\\github.com\\Azure\\azapi-lsp\\中文",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getWorkingDirectory(tt.input, tt.os); got != tt.want {
				t.Errorf("getWorkingDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}
