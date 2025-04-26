package l

import (
	"sort"
	"strings"
	"testing"
)

func TestPluralize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Regular word", "password", "passwords"},
		{"Word ending in y", "identity", "identities"},
		{"Word ending in s", "address", "addresses"},
		{"Word ending in x", "index", "indexes"},
		{"Word ending in z", "quiz", "quizzes"},
		{"Word ending in ch", "match", "matches"},
		{"Word ending in sh", "hash", "hashes"},
		{"Compound word", "apiKey", "apiKeys"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pluralize(tt.input)
			if result != tt.expected {
				t.Errorf("pluralize(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCaseConversions(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantCamel  string
		wantPascal string
		wantSnake  string
		wantKebab  string
	}{
		{
			name:       "Simple word",
			input:      "password",
			wantCamel:  "password",
			wantPascal: "Password",
			wantSnake:  "password",
			wantKebab:  "password",
		},
		{
			name:       "Compound camelCase",
			input:      "apiKey",
			wantCamel:  "apiKey",
			wantPascal: "ApiKey",
			wantSnake:  "api_key",
			wantKebab:  "api-key",
		},
		{
			name:       "Compound PascalCase",
			input:      "ApiKey",
			wantCamel:  "apiKey",
			wantPascal: "ApiKey",
			wantSnake:  "api_key",
			wantKebab:  "api-key",
		},
		{
			name:       "Snake case",
			input:      "api_key",
			wantCamel:  "apiKey",
			wantPascal: "ApiKey",
			wantSnake:  "api_key",
			wantKebab:  "api-key",
		},
		{
			name:       "Kebab case",
			input:      "api-key",
			wantCamel:  "apiKey",
			wantPascal: "ApiKey",
			wantSnake:  "api_key",
			wantKebab:  "api-key",
		},
		{
			name:       "Multiple words",
			input:      "database_connection_string",
			wantCamel:  "databaseConnectionString",
			wantPascal: "DatabaseConnectionString",
			wantSnake:  "database_connection_string",
			wantKebab:  "database-connection-string",
		},
		{
			name:       "With numbers",
			input:      "oauth2_client",
			wantCamel:  "oauth2Client",
			wantPascal: "Oauth2Client",
			wantSnake:  "oauth2_client",
			wantKebab:  "oauth2-client",
		},
		{
			name:       "Mixed case with acronym",
			input:      "AWS_secret_key",
			wantCamel:  "awsSecretKey",
			wantPascal: "AwsSecretKey",
			wantSnake:  "aws_secret_key",
			wantKebab:  "aws-secret-key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toCamelCase(tt.input); got != tt.wantCamel {
				t.Errorf("toCamelCase(%q) = %q, want %q", tt.input, got, tt.wantCamel)
			}
			if got := toPascalCase(tt.input); got != tt.wantPascal {
				t.Errorf("toPascalCase(%q) = %q, want %q", tt.input, got, tt.wantPascal)
			}
			if got := toSnakeCase(tt.input); got != tt.wantSnake {
				t.Errorf("toSnakeCase(%q) = %q, want %q", tt.input, got, tt.wantSnake)
			}
			if got := toKebabCase(tt.input); got != tt.wantKebab {
				t.Errorf("toKebabCase(%q) = %q, want %q", tt.input, got, tt.wantKebab)
			}
		})
	}
}

func TestSplitWords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"Simple word", "password", []string{"password"}},
		{"Camel case", "apiKey", []string{"api", "Key"}},
		{"Pascal case", "ApiKey", []string{"Api", "Key"}},
		{"Snake case", "api_key", []string{"api", "key"}},
		{"Kebab case", "api-key", []string{"api", "key"}},
		{"Multiple words", "database_connection_string", []string{"database", "connection", "string"}},
		{"Mixed case", "OAuth2ClientSecret", []string{"OAuth2", "Client", "Secret"}},
		{"With spaces", "api key", []string{"api", "key"}},
		{"Complex mixed", "AWS_secretKey-Value", []string{"AWS", "secret", "Key", "Value"}},
		{"With numbers", "oauth2Client", []string{"oauth2", "Client"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitWords(tt.input)
			if !stringSliceEqual(result, tt.expected) {
				t.Errorf("splitWords(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGenerateKeywordVariations(t *testing.T) {
	variations := GenerateKeywordVariations()

	// Test for uniqueness
	t.Run("Uniqueness", func(t *testing.T) {
		seen := make(map[string]bool)
		duplicates := []string{}
		for _, v := range variations {
			if seen[v] {
				duplicates = append(duplicates, v)
			}
			seen[v] = true
		}
		if len(duplicates) > 0 {
			t.Errorf("Found %d duplicate variations: %v", len(duplicates), duplicates)
		}
	})

	// Test for expected variations of specific keywords
	t.Run("Expected variations", func(t *testing.T) {
		expectedVariations := []string{
			// password variations
			"password", "passwords", "Password", "Passwords",
			// api key variations
			"apikey", "apikeys", "api_key", "api_keys", "api-key", "api-keys",
			"ApiKey", "ApiKeys",
			// credit card variations
			"creditcard", "creditcards", "credit_card", "credit_cards",
			"credit-card", "credit-cards", "CreditCard", "CreditCards",
		}

		missing := []string{}
		for _, expected := range expectedVariations {
			if !containsString(variations, expected) {
				missing = append(missing, expected)
			}
		}
		if len(missing) > 0 {
			t.Errorf("Missing expected variations: %v", missing)
		}
	})

	// Test that all base keywords are included
	t.Run("Base keywords included", func(t *testing.T) {
		missing := []string{}
		for _, base := range baseKeywords {
			if !containsString(variations, base) {
				missing = append(missing, base)
			}
		}
		if len(missing) > 0 {
			t.Errorf("Missing base keywords: %v", missing)
		}
	})

	// Test compound word generation
	t.Run("Compound words", func(t *testing.T) {
		compounds := []string{
			"apikey", "apitoken", "apicredential",
			"userpassword", "usertoken", "useremail",
		}

		missing := []string{}
		for _, compound := range compounds {
			found := false
			for _, v := range variations {
				if strings.ToLower(v) == compound {
					found = true
					break
				}
			}
			if !found {
				missing = append(missing, compound)
			}
		}
		if len(missing) > 0 {
			t.Errorf("Missing compound words: %v", missing)
		}
	})

	// Test for common patterns
	t.Run("Common patterns", func(t *testing.T) {
		patterns := []struct {
			base     string
			expected []string
		}{
			{
				base: "password",
				expected: []string{
					"password", "passwords",
					"userPassword", "userPasswords",
					"password_hash", "password_hashes",
				},
			},
			{
				base: "key",
				expected: []string{
					"key", "keys",
					"apiKey", "apiKeys",
					"secret_key", "secret_keys",
				},
			},
		}

		for _, p := range patterns {
			missing := []string{}
			for _, expected := range p.expected {
				if !containsString(variations, expected) {
					missing = append(missing, expected)
				}
			}
			if len(missing) > 0 {
				t.Errorf("Missing patterns for %q: %v", p.base, missing)
			}
		}
	})
}

// Helper function to compare string slices
func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Helper function to check if a string is in a slice
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if strings.EqualFold(item, s) { // Case-insensitive comparison
			return true
		}
	}
	return false
}

// Benchmark the generation of keyword variations
func BenchmarkGenerateKeywordVariations(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		GenerateKeywordVariations()
	}
}

// Example test to demonstrate usage
func ExampleGenerateKeywordVariations() {
	// Get a small subset of variations for demonstration
	variations := GenerateKeywordVariations()

	// Sort them for consistent output
	subset := make([]string, 0)
	for _, v := range variations {
		if strings.Contains(strings.ToLower(v), "password") {
			subset = append(subset, v)
		}
	}
	sort.Strings(subset)

	// Print the first few password-related variations
	for i := 0; i < len(subset) && i < 5; i++ {
		println(subset[i])
	}
}
