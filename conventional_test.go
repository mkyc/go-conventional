package nextver

import (
	"reflect"
	"testing"
)

func TestParseConventionalCommit(t *testing.T) {
	tests := []struct {
		name          string
		commitMessage string
		expected      ConventionalCommit
	}{
		{
			name: "Commit message with description and breaking change footer",
			commitMessage: `feat: allow provided config object to extend other configs

BREAKING CHANGE: ` + "`extends`" + ` key in config file is now used for extending other config files`,
			expected: ConventionalCommit{
				Type:         Feat,
				Description:  "allow provided config object to extend other configs",
				Breaking:     true,
				Conventional: true,
				Footers: map[string]string{
					"BREAKING CHANGE": "`extends` key in config file is now used for extending other config files",
				},
			},
		},
		{
			name:          "Commit message with ! to draw attention to breaking change",
			commitMessage: "feat!: send an email to the customer when a product is shipped",
			expected: ConventionalCommit{
				Type:         Feat,
				Description:  "send an email to the customer when a product is shipped",
				Breaking:     true,
				Conventional: true,
				Footers:      map[string]string{},
			},
		},
		{
			name:          "Commit message with scope and ! to draw attention to breaking change",
			commitMessage: "feat(api)!: send an email to the customer when a product is shipped",
			expected: ConventionalCommit{
				Type:         Feat,
				Scope:        "api",
				Description:  "send an email to the customer when a product is shipped",
				Breaking:     true,
				Conventional: true,
				Footers:      map[string]string{},
			},
		},
		{
			name: "Commit message with both ! and BREAKING CHANGE footer",
			commitMessage: `chore!: drop support for Node 6

BREAKING CHANGE: use JavaScript features not available in Node 6.`,
			expected: ConventionalCommit{
				Type:         Chore,
				Description:  "drop support for Node 6",
				Breaking:     true,
				Conventional: true,
				Footers: map[string]string{
					"BREAKING CHANGE": "use JavaScript features not available in Node 6.",
				},
			},
		},
		{
			name:          "Commit message with no body",
			commitMessage: "docs: correct spelling of CHANGELOG",
			expected: ConventionalCommit{
				Type:         Docs,
				Description:  "correct spelling of CHANGELOG",
				Conventional: true,
				Footers:      map[string]string{},
			},
		},
		{
			name:          "Commit message with scope",
			commitMessage: "feat(lang): add Polish language",
			expected: ConventionalCommit{
				Type:         Feat,
				Scope:        "lang",
				Description:  "add Polish language",
				Conventional: true,
				Footers:      map[string]string{},
			},
		},
		{
			name: "Commit message with multi-paragraph body and multiple footers",
			commitMessage: `fix: prevent racing of requests

Introduce a request id and a reference to latest request. Dismiss
incoming responses other than from latest request.

Remove timeouts which were used to mitigate the racing issue but are
obsolete now.

Reviewed-by: Z
Refs: #123`,
			expected: ConventionalCommit{
				Type:         Fix,
				Description:  "prevent racing of requests",
				Body:         "Introduce a request id and a reference to latest request. Dismiss\nincoming responses other than from latest request.\n\nRemove timeouts which were used to mitigate the racing issue but are\nobsolete now.",
				Conventional: true,
				Footers: map[string]string{
					"Reviewed-by": "Z",
					"Refs":        "#123",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseConventionalCommit(tt.commitMessage)

			if result.Type != tt.expected.Type {
				t.Errorf("Type = %v, want %v", result.Type, tt.expected.Type)
			}

			if result.Scope != tt.expected.Scope {
				t.Errorf("Scope = %v, want %v", result.Scope, tt.expected.Scope)
			}

			if result.Description != tt.expected.Description {
				t.Errorf("Description = %v, want %v", result.Description, tt.expected.Description)
			}

			if result.Body != tt.expected.Body {
				t.Errorf("Body = %v, want %v", result.Body, tt.expected.Body)
			}

			if result.Breaking != tt.expected.Breaking {
				t.Errorf("Breaking = %v, want %v", result.Breaking, tt.expected.Breaking)
			}

			if result.Conventional != tt.expected.Conventional {
				t.Errorf("Conventional = %v, want %v", result.Conventional, tt.expected.Conventional)
			}

			if !reflect.DeepEqual(result.Footers, tt.expected.Footers) {
				t.Errorf("Footers = %v, want %v", result.Footers, tt.expected.Footers)
			}
		})
	}
}
