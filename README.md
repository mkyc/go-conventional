# go-conventional

A Go implementation of the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification.

## Overview

This package provides functionality to parse and validate commit messages according to the Conventional Commits specification v1.0.0. It allows you to extract structured information from commit messages, including type, scope, description, body, footers, and breaking change indicators.

## Installation

```bash
go get github.com/mkyc/go-conventional
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/mkyc/go-conventional"
)

func main() {
	// Parse a conventional commit message
	commitMsg := `feat(api): add new endpoint for user authentication

This commit adds a new REST API endpoint for user authentication
with JWT tokens.

Reviewed-by: Jane Doe
BREAKING CHANGE: The old authentication endpoint is removed`

	commit := conventional.ParseConventionalCommit(commitMsg)

	// Access the parsed information
	fmt.Printf("Type: %s\n", commit.Type)
	fmt.Printf("Scope: %s\n", commit.Scope)
	fmt.Printf("Description: %s\n", commit.Description)
	fmt.Printf("Body: %s\n", commit.Body)
	fmt.Printf("Is Breaking Change: %t\n", commit.Breaking)
	fmt.Printf("Is Conventional: %t\n", commit.Conventional)
	
	// Access footers
	for token, value := range commit.Footers {
		fmt.Printf("Footer - %s: %s\n", token, value)
	}
}
```

## Supported Commit Types

The package recognizes the following conventional commit types:

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `build`: Changes that affect the build system or external dependencies
- `ci`: Changes to CI configuration files and scripts
- `chore`: Other changes that don't modify src or test files
- `revert`: Reverts a previous commit

## Implementation Details

The package implements the full Conventional Commits specification v1.0.0, including:

- Parsing commit types and optional scopes
- Detecting breaking changes (via `!` or `BREAKING CHANGE` footer)
- Extracting commit descriptions, bodies, and footers
- Handling multi-line descriptions and bodies
- Processing various footer formats

The implementation is case-insensitive for commit types as per the specification, but preserves the case of other elements like scope, description, and footers.

## Specification Compliance

This implementation follows the [Conventional Commits specification v1.0.0](https://www.conventionalcommits.org/en/v1.0.0/#specification), adhering to all the rules defined in the specification (or at least that is what I hope). Key points include:

1. Commits are prefixed with a type, optional scope, optional breaking change indicator, and required colon and space
2. The types `feat` and `fix` are used for new features and bug fixes respectively
3. Other types are supported for different kinds of changes
4. Breaking changes are indicated either in the type/scope prefix with `!` or as a footer
5. Footers follow the git trailer convention

## License

See the [LICENSE](LICENSE) file for details.
