package nextver

import (
	"regexp"
	"strings"
)

// CommitType represents the type of conventional commit.
type CommitType string

// feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert, BREAKING CHANGE
const (
	Feat     CommitType = "feat"
	Fix      CommitType = "fix"
	Docs     CommitType = "docs"
	Style    CommitType = "style"
	Refactor CommitType = "refactor"
	Perf     CommitType = "perf"
	Test     CommitType = "test"
	Build    CommitType = "build"
	Ci       CommitType = "ci"
	Chore    CommitType = "chore"
	Revert   CommitType = "revert"
	Other    CommitType = "other" // for commits that do not match any conventional type
)

// ConventionalCommit represents a conventional commit message.
type ConventionalCommit struct {
	Type         CommitType
	Scope        string            // optional scope
	Description  string            // short description of the commit
	Body         string            // optional body of the commit
	Footers      map[string]string // optional footers, e.g., "Reviewed-by", "Co-authored-by", "Refs", "BREAKING CHANGE", etc.
	Breaking     bool              // true if this commit is a breaking change
	Conventional bool              // true if this commit follows the conventional commit format
}

// ParseConventionalCommit parses a conventional commit message and returns a ConventionalCommit struct.
// Specification (as of 2025-06-05 at https://www.conventionalcommits.org/en/v1.0.0/#specification):
//
// The key words “MUST”, “MUST NOT”, “REQUIRED”, “SHALL”, “SHALL NOT”, “SHOULD”, “SHOULD NOT”, “RECOMMENDED”, “MAY”, and “OPTIONAL” in this document are to be interpreted as described in RFC 2119.
//
// 1. Commits MUST be prefixed with a type, which consists of a noun, feat, fix, etc., followed by the OPTIONAL scope, OPTIONAL !, and REQUIRED terminal colon and space.
// 2. The type feat MUST be used when a commit adds a new feature to your application or library.
// 3. The type fix MUST be used when a commit represents a bug fix for your application.
// 4. A scope MAY be provided after a type. A scope MUST consist of a noun describing a section of the codebase surrounded by parenthesis, e.g., fix(parser):
// 5. A description MUST immediately follow the colon and space after the type/scope prefix. The description is a short summary of the code changes, e.g., fix: array parsing issue when multiple spaces were contained in string.
// 6. A longer commit body MAY be provided after the short description, providing additional contextual information about the code changes. The body MUST begin one blank line after the description.
// 7. A commit body is free-form and MAY consist of any number of newline separated paragraphs.
// 8. One or more footers MAY be provided one blank line after the body. Each footer MUST consist of a word token, followed by either a :<space> or <space># separator, followed by a string value (this is inspired by the git trailer convention).
// 9. A footer’s token MUST use - in place of whitespace characters, e.g., Acked-by (this helps differentiate the footer section from a multi-paragraph body). An exception is made for BREAKING CHANGE, which MAY also be used as a token.
// 10. A footer’s value MAY contain spaces and newlines, and parsing MUST terminate when the next valid footer token/separator pair is observed.
// 11. Breaking changes MUST be indicated in the type/scope prefix of a commit, or as an entry in the footer.
// 12. If included as a footer, a breaking change MUST consist of the uppercase text BREAKING CHANGE, followed by a colon, space, and description, e.g., BREAKING CHANGE: environment variables now take precedence over config files.
// 13. If included in the type/scope prefix, breaking changes MUST be indicated by a ! immediately before the :. If ! is used, BREAKING CHANGE: MAY be omitted from the footer section, and the commit description SHALL be used to describe the breaking change.
// 14. Types other than feat and fix MAY be used in your commit messages, e.g., docs: update ref docs.
// 15. The units of information that make up Conventional Commits MUST NOT be treated as case sensitive by implementors, with the exception of BREAKING CHANGE which MUST be uppercase.
// 16. BREAKING-CHANGE MUST be synonymous with BREAKING CHANGE, when used as a token in a footer.
func ParseConventionalCommit(commitMessage string) ConventionalCommit {
	result := ConventionalCommit{
		Type:         Other,
		Conventional: false,
		Footers:      make(map[string]string),
	}

	// If empty message, return default
	if commitMessage == "" {
		return result
	}

	// Split the commit message into lines
	lines := strings.Split(commitMessage, "\n")
	if len(lines) == 0 {
		return result
	}

	// Parse the first line (header)
	header := lines[0]

	// Regular expression to match the conventional commit format
	// Format: type(scope)!: description
	// All parts except type and description are optional
	headerRegex := regexp.MustCompile(`^(?i)(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)(?:\(([^\)]+)\))?(!)?:\s+(.+)$`)
	matches := headerRegex.FindStringSubmatch(header)

	if len(matches) == 0 {
		// Not a conventional commit
		return result
	}

	// Mark as conventional since it matched the format
	result.Conventional = true

	// Extract type (case-insensitive as per rule 15)
	commitType := strings.ToLower(matches[1])
	switch commitType {
	case "feat":
		result.Type = Feat
	case "fix":
		result.Type = Fix
	case "docs":
		result.Type = Docs
	case "style":
		result.Type = Style
	case "refactor":
		result.Type = Refactor
	case "perf":
		result.Type = Perf
	case "test":
		result.Type = Test
	case "build":
		result.Type = Build
	case "ci":
		result.Type = Ci
	case "chore":
		result.Type = Chore
	case "revert":
		result.Type = Revert
	default:
		result.Type = Other
	}

	// Extract scope if present
	if matches[2] != "" {
		result.Scope = matches[2]
	}

	// Check for breaking change indicator in header
	if matches[3] == "!" {
		result.Breaking = true
	}

	// Extract description
	result.Description = matches[4]

	// Process the rest of the message (body and footers)
	if len(lines) > 1 {
		// First, identify all footer lines
		footerLines := []int{}
		bodyLines := []int{}

		// Skip the header and any blank lines immediately after it
		i := 1
		for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
			i++
		}

		// Process remaining lines
		for ; i < len(lines); i++ {
			line := strings.TrimSpace(lines[i])

			// Check if this line is a footer
			if isFooterLine(line) {
				footerLines = append(footerLines, i)
			} else if line != "" || len(bodyLines) > 0 {
				// Non-empty line or we've already started collecting body lines
				bodyLines = append(bodyLines, i)
			}
		}

		// Extract body if there are body lines
		if len(bodyLines) > 0 {
			// Check if any footer lines are mixed with body lines
			bodyEnd := bodyLines[len(bodyLines)-1]
			validBodyLines := true

			for _, footerLine := range footerLines {
				if footerLine < bodyEnd {
					// Footer line is before the end of the body, so it's part of the body
					validBodyLines = false
					break
				}
			}

			if validBodyLines {
				// Extract body with correct newline handling
				bodyText := make([]string, 0, len(bodyLines))

				// Group lines by paragraphs
				var currentParagraph []string
				for i, lineIdx := range bodyLines {
					line := lines[lineIdx]

					if strings.TrimSpace(line) == "" {
						// Empty line marks the end of a paragraph
						if len(currentParagraph) > 0 {
							// Join the paragraph lines with newlines
							bodyText = append(bodyText, strings.Join(currentParagraph, "\n"))
							currentParagraph = nil
						}
					} else {
						// Add line to current paragraph
						currentParagraph = append(currentParagraph, line)
					}

					// Handle the last paragraph
					if i == len(bodyLines)-1 && len(currentParagraph) > 0 {
						bodyText = append(bodyText, strings.Join(currentParagraph, "\n"))
					}
				}

				// Join paragraphs with double newlines
				result.Body = strings.Join(bodyText, "\n\n")
			}
		}

		// Process footers
		for _, lineIdx := range footerLines {
			line := strings.TrimSpace(lines[lineIdx])

			// Handle BREAKING CHANGE footers
			if strings.HasPrefix(line, "BREAKING CHANGE:") {
				value := strings.TrimSpace(line[len("BREAKING CHANGE:"):])
				result.Footers["BREAKING CHANGE"] = value
				result.Breaking = true
			} else if strings.HasPrefix(line, "BREAKING-CHANGE:") {
				value := strings.TrimSpace(line[len("BREAKING-CHANGE:"):])
				result.Footers["BREAKING-CHANGE"] = value
				result.Breaking = true
			} else {
				// Regular footer
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					token := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					result.Footers[token] = value
				} else {
					// Try hash format
					parts = strings.SplitN(line, " # ", 2)
					if len(parts) == 2 {
						token := strings.TrimSpace(parts[0])
						value := strings.TrimSpace(parts[1])
						result.Footers[token] = value
					}
				}
			}
		}
	}

	return result
}

// isFooterLine checks if a line is a footer line
func isFooterLine(line string) bool {
	if line == "" {
		return false
	}

	// Check for BREAKING CHANGE footers
	if strings.HasPrefix(line, "BREAKING CHANGE:") || strings.HasPrefix(line, "BREAKING-CHANGE:") {
		return true
	}

	// Check for regular footers (token: value or token # value)
	colonRegex := regexp.MustCompile(`^([A-Za-z0-9-]+):\s+(.+)$`)
	hashRegex := regexp.MustCompile(`^([A-Za-z0-9-]+)\s+#\s+(.+)$`)

	return colonRegex.MatchString(line) || hashRegex.MatchString(line)
}
