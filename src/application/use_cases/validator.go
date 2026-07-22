package usecases

import (
	"fmt"
	"regexp"
	"strings"
)

// dangerousKeywords contains SQL keywords that are NOT allowed.
var dangerousKeywords = []string{
	"INSERT", "UPDATE", "DELETE", "DROP", "ALTER", "CREATE",
	"TRUNCATE", "REPLACE", "GRANT", "REVOKE", "EXEC",
	"EXECUTE", "ATTACH", "DETACH", "PRAGMA", "VACUUM",
	"REINDEX", "SAVEPOINT", "RELEASE", "ROLLBACK", "COMMIT",
	"BEGIN", "END TRANSACTION",
}

// multiStatementRegex detects multiple SQL statements (;) which could be injection.
var multiStatementRegex = regexp.MustCompile(`;[\s]*\S`)

// ValidateSQL checks if a SQL query is safe to execute.
// It only allows SELECT statements and blocks dangerous keywords.
func ValidateSQL(sql string) error {
	if sql == "" {
		return fmt.Errorf("consulta SQL vacia")
	}

	normalized := strings.TrimSpace(sql)
	upper := strings.ToUpper(normalized)

	// Must start with SELECT or WITH (for CTEs)
	if !strings.HasPrefix(upper, "SELECT") && !strings.HasPrefix(upper, "WITH") {
		return fmt.Errorf("solo se permiten consultas SELECT (encontrado: %.20s...)", normalized)
	}

	// Check for dangerous keywords (as whole words)
	for _, kw := range dangerousKeywords {
		pattern := fmt.Sprintf(`\b%s\b`, kw)
		matched, _ := regexp.MatchString(pattern, upper)
		if matched {
			return fmt.Errorf("palabra clave no permitida: %s", kw)
		}
	}

	// Check for multiple statements (SQL injection attempt)
	if multiStatementRegex.MatchString(normalized) {
		return fmt.Errorf("no se permiten multiples sentencias SQL")
	}

	// Check for comment injection
	if strings.Contains(normalized, "--") || strings.Contains(normalized, "/*") {
		return fmt.Errorf("no se permiten comentarios SQL")
	}

	return nil
}
