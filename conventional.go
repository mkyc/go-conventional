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
		Conventional: true,
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

	// If the header matches the conventional commit format
	if len(matches) > 0 {
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

		// Extract description from header
		description := matches[4]
		result.Description = description
	} else {
		// If the header does not match the conventional commit format, treat it as a non-conventional commit
		result.Conventional = false
		result.Description = header
	}

	if len(lines) < 2 {
		// If there is neither body nor footers, return the result with just the header
		return result
	}

	// Process the rest of the message (description continuation, body, and footers)
	// Find the first empty line after the header
	firstEmptyLine := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "" {
			firstEmptyLine = i
			break
		}
	}

	// If there's no empty line, all lines are part of the description
	if firstEmptyLine == -1 {
		// Append all lines to the description
		for i := 1; i < len(lines); i++ {
			result.Description += "\n" + lines[i]
		}
		return result
	}

	// Lines between header and first empty line are part of the description
	if firstEmptyLine > 1 {
		for i := 1; i < firstEmptyLine; i++ {
			result.Description += "\n" + lines[i]
		}
	}

	// Skip the header, description continuation, and the empty line
	bodyStart := firstEmptyLine + 1
	for bodyStart < len(lines) && strings.TrimSpace(lines[bodyStart]) == "" {
		bodyStart++
	}

	// Find the first footer line
	firstFooterLine := -1
	for i := bodyStart; i < len(lines); i++ {
		if isFooterLine(strings.TrimSpace(lines[i])) {
			firstFooterLine = i
			break
		}
	}

	// Extract body (everything between header and first footer)
	if bodyStart < len(lines) && (firstFooterLine == -1 || bodyStart < firstFooterLine) {
		bodyEnd := len(lines)
		if firstFooterLine != -1 {
			bodyEnd = firstFooterLine
		}

		// Join all body lines, preserving empty lines
		result.Body = strings.TrimSpace(strings.Join(lines[bodyStart:bodyEnd], "\n"))
	}

	if firstFooterLine == -1 {
		// If no footer lines found, return the result with just the header, description, and body
		return result
	}

	// Process footers (everything after the first footer line)
	// Find all footer line numbers
	var footerLineNumbers []int
	for i := firstFooterLine; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line != "" && isFooterLine(line) {
			footerLineNumbers = append(footerLineNumbers, i)
		}
	}

	// Divide text into blocks starting with footer lines
	var footerBlocks [][]string
	for i := range footerLineNumbers {
		startIdx := footerLineNumbers[i]
		endIdx := len(lines)
		if i < len(footerLineNumbers)-1 {
			endIdx = footerLineNumbers[i+1]
		}

		// Create a block of text for this footer
		block := make([]string, 0, endIdx-startIdx)
		for j := startIdx; j < endIdx; j++ {
			block = append(block, lines[j])
		}
		footerBlocks = append(footerBlocks, block)
	}

	// Process each footer block
	for _, block := range footerBlocks {
		token, value := extractFooter(block)
		if token != "" && value != "" {
			result.Footers[token] = value
			// Check for breaking change
			if token == "BREAKING CHANGE" || token == "BREAKING-CHANGE" {
				result.Breaking = true
			}
		}
	}

	return result
}

// patternMatcher is a utility function to return a matches
func patternMatcher(line string) (string, string) {
	// Define regex patterns for different footer formats
	patterns := []struct {
		regex      *regexp.Regexp
		valueIndex int
	}{
		// Breaking change format with optional colon
		{regexp.MustCompile(`^(BREAKING CHANGE|BREAKING-CHANGE)(:)?(.*)$`), 3},
		// Colon format: token: value
		{regexp.MustCompile(`^([A-Za-z0-9-]+):\s+(.+)$`), 2},
		// Hash format: token #value
		{regexp.MustCompile(`^([A-Za-z0-9-]+)\s+#(.+)$`), 2},
	}

	for _, pattern := range patterns {
		matches := pattern.regex.FindStringSubmatch(line)
		if len(matches) > 0 {
			token := matches[1]
			value := strings.TrimSpace(matches[pattern.valueIndex])
			return token, value
		}
	}

	return "", ""
}

// extractFooter takes a block of text starting with a footer line and returns the token and value
func extractFooter(block []string) (string, string) {
	token, value := patternMatcher(block[0])
	if len(token) == 0 || len(value) == 0 {
		// No valid footer line found
		return "", ""
	}
	// If the footer line is a valid footer, we can extract the value
	for i := 1; i < len(block); i++ {
		line := strings.TrimSpace(block[i])
		if line != "" {
			value += "\n" + line
		}
	}

	return token, value
}

// isFooterLine checks if a line is a footer line
func isFooterLine(line string) bool {
	token, value := patternMatcher(line)
	return len(token) > 0 && len(value) > 0
}
