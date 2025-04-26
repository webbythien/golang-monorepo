package l

import (
	"strings"
	"unicode"
)

// baseKeywords contains the base words that need to be protected
var baseKeywords = []string{
	// Authentication & Authorization
	"password", "pass", "token", "secret", "signature", "credential", "authorization", "auth",
	"jwt", "session", "identity", "oauth", "key",

	// Personal Information
	"email", "phone", "address", "passport", "license", "account", "ssn", "social", "security",
	"tax", "dob", "birth", "name",

	// Financial
	"credit", "card", "cvv", "ccv", "pin", "bank",

	// Infrastructure & Configuration
	"config", "endpoint", "connection", "url", "host", "ip", "server", "database", "db",
	"user", "admin", "root",

	// Cloud & Services
	"aws", "gcp", "azure", "encryption", "certificate", "webhook", "smtp", "mailgun", "sendgrid",
	"slack", "discord", "twilio",

	// Other Sensitive
	"salt", "pepper", "private", "access", "refresh", "client", "api",
}

// wordVariations generates different forms of a word
type wordVariations struct {
	Original    string   // Original form
	Plural      string   // Plural form
	CamelCase   string   // camelCase
	PascalCase  string   // PascalCase
	SnakeCase   string   // snake_case
	KebabCase   string   // kebab-case
	AllVariants []string // All variations combined
}

// pluralize returns the plural form of a word
func pluralize(word string) string {
	// Handle compound words by finding the last word
	parts := splitWords(word)
	if len(parts) > 1 {
		lastWord := parts[len(parts)-1]
		pluralLastWord := pluralize(lastWord)
		return strings.Join(parts[:len(parts)-1], "") + pluralLastWord
	}

	switch {
	case strings.HasSuffix(word, "y"):
		// Handle words ending in consonant + y
		if len(word) > 1 && !isVowel(rune(word[len(word)-2])) {
			return strings.TrimSuffix(word, "y") + "ies"
		}
		return word + "s"
	case strings.HasSuffix(word, "z"):
		return word + word[len(word)-1:] + "es" // "quiz" -> "quizzes"
	case strings.HasSuffix(word, "s"), strings.HasSuffix(word, "x"),
		strings.HasSuffix(word, "ch"), strings.HasSuffix(word, "sh"):
		return word + "es"
	default:
		return word + "s"
	}
}

// isVowel checks if a character is a vowel
func isVowel(c rune) bool {
	return strings.ContainsRune("aeiouAEIOU", c)
}

// containsOnlyNumbers checks if a string contains only numbers
func containsOnlyNumbers(s string) bool {
	for _, r := range s {
		if !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}

// isAcronym checks if a word is likely an acronym
func isAcronym(s string) bool {
	if len(s) <= 1 {
		return false
	}

	// Special cases for common acronyms
	commonAcronyms := map[string]bool{
		"AWS":    true,
		"API":    true,
		"OAuth":  true,
		"OAuth2": true,
	}
	if commonAcronyms[s] {
		return true
	}

	// Check if all characters are uppercase or numbers
	for _, r := range s {
		if !unicode.IsUpper(r) && !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}

// splitWords splits a string into words based on case, spaces, underscores, and hyphens
func splitWords(s string) []string {
	var words []string
	var currentWord strings.Builder

	// Helper function to determine character type
	getCharType := func(r rune) int {
		switch {
		case unicode.IsUpper(r):
			return 2
		case unicode.IsLower(r):
			return 1
		case unicode.IsNumber(r):
			return 3
		default:
			return 4
		}
	}

	// Helper function to add word and reset builder
	addWord := func() {
		if currentWord.Len() > 0 {
			word := currentWord.String()
			// Check if the current word is part of a known acronym
			if len(words) > 0 {
				prevWord := words[len(words)-1]
				combined := prevWord + word
				if isAcronym(combined) {
					words[len(words)-1] = combined
					currentWord.Reset()
					return
				}
			}
			words = append(words, word)
			currentWord.Reset()
		}
	}

	for i, r := range s {
		currentType := getCharType(r)

		// Handle separators
		if r == '_' || r == '-' || r == ' ' {
			addWord()
			continue
		}

		// Handle transitions
		if i > 0 {
			prevType := getCharType(rune(s[i-1]))

			// Start new word on:
			// 1. After a separator
			// 2. Lower to upper transition (except in acronyms)
			// 3. Number to letter transition
			// 4. End of acronym (multiple uppercase followed by lowercase)
			if prevType == 4 || // after separator
				(prevType == 1 && currentType == 2) || // lower to upper
				(prevType == 3 && (currentType == 1 || currentType == 2)) || // number to letter
				(prevType == 2 && currentType == 1 && // upper to lower
					i > 1 && getCharType(rune(s[i-2])) == 2 && // but not first letter of word
					!isAcronym(currentWord.String())) { // and not part of acronym
				addWord()
			}
		}

		currentWord.WriteRune(r)
	}

	addWord()

	// Post-process words to handle acronyms and numbers
	for i := 0; i < len(words); i++ {
		// Join numbers with their preceding word if it exists
		if i > 0 && containsOnlyNumbers(words[i]) {
			words[i-1] += words[i]
			words = append(words[:i], words[i+1:]...)
			i--
			continue
		}
	}

	return words
}

// toCamelCase converts a string to camelCase
func toCamelCase(s string) string {
	words := splitWords(s)
	if len(words) == 0 {
		return ""
	}

	// First word is always lowercase unless it's an acronym
	result := words[0]
	if isAcronym(result) {
		result = strings.ToLower(result)
	}

	// Subsequent words are capitalized
	for i := 1; i < len(words); i++ {
		word := words[i]
		if isAcronym(word) {
			result += word
		} else {
			result += strings.Title(strings.ToLower(word))
		}
	}
	return result
}

// toPascalCase converts a string to PascalCase
func toPascalCase(s string) string {
	words := splitWords(s)
	var result string
	for _, word := range words {
		result += strings.Title(strings.ToLower(word))
	}
	return result
}

// toSnakeCase converts a string to snake_case
func toSnakeCase(s string) string {
	words := splitWords(s)
	for i := range words {
		words[i] = strings.ToLower(words[i])
	}
	return strings.Join(words, "_")
}

// toKebabCase converts a string to kebab-case
func toKebabCase(s string) string {
	words := splitWords(s)
	for i := range words {
		words[i] = strings.ToLower(words[i])
	}
	return strings.Join(words, "-")
}

// GenerateKeywordVariations generates all variations of the base keywords
func GenerateKeywordVariations() []string {
	var allVariations []string
	seen := make(map[string]struct{})

	// Helper to add unique variations
	addUnique := func(s string) {
		if s != "" {
			if _, exists := seen[s]; !exists {
				allVariations = append(allVariations, s)
				seen[s] = struct{}{}
			}
		}
	}

	// Process each base keyword
	for _, base := range baseKeywords {
		v := wordVariations{
			Original:   base,
			Plural:     pluralize(base),
			CamelCase:  toCamelCase(base),
			PascalCase: toPascalCase(base),
			SnakeCase:  toSnakeCase(base),
			KebabCase:  toKebabCase(base),
		}

		// Add all forms
		addUnique(v.Original)
		addUnique(v.Plural)
		addUnique(v.CamelCase)
		addUnique(v.PascalCase)
		addUnique(v.SnakeCase)
		addUnique(v.KebabCase)

		// Generate plural forms of different cases
		addUnique(toCamelCase(v.Plural))
		addUnique(toPascalCase(v.Plural))
		addUnique(toSnakeCase(v.Plural))
		addUnique(toKebabCase(v.Plural))

		// Handle compound words (e.g., "apiKey", "creditCard")
		for _, other := range baseKeywords {
			if base != other {
				// Create compound in different forms
				compound := base + other
				addUnique(compound)
				addUnique(pluralize(compound))
				addUnique(toCamelCase(compound))
				addUnique(toPascalCase(compound))
				addUnique(toSnakeCase(compound))
				addUnique(toKebabCase(compound))

				// Add plural forms of compounds
				pluralCompound := pluralize(compound)
				addUnique(toCamelCase(pluralCompound))
				addUnique(toPascalCase(pluralCompound))
				addUnique(toSnakeCase(pluralCompound))
				addUnique(toKebabCase(pluralCompound))

				// Add special compound patterns
				if (base == "password" && other == "hash") ||
					(base == "secret" && other == "key") {
					snakeCase := base + "_" + other
					addUnique(snakeCase)
					addUnique(pluralize(snakeCase))
				}
			}
		}
	}

	// Add special compound patterns that might be missed
	specialPatterns := []string{
		"api_key", "api_keys",
		"api-key", "api-keys",
		"credit_card", "credit_cards",
		"credit-card", "credit-cards",
		"password_hash", "password_hashes",
		"secret_key", "secret_keys",
	}
	for _, pattern := range specialPatterns {
		addUnique(pattern)
	}

	return allVariations
}
