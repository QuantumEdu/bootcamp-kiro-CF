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

// allowedTables is the whitelist of tables that can be queried.
var allowedTables = []string{
	"productos", "categorias", "ventas", "venta_items",
	"clientes", "usuarios", "inventario_movimientos", "configuracion",
}

// jailbreakPatterns detects prompt injection attempts.
var jailbreakPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)ignor[ae]\s+(las\s+)?instrucciones`),
	regexp.MustCompile(`(?i)ignore\s+(all\s+)?(previous\s+)?instructions`),
	regexp.MustCompile(`(?i)forget\s+(all\s+)?(previous|your)\s+`),
	regexp.MustCompile(`(?i)olvida\s+(las\s+)?instrucciones`),
	regexp.MustCompile(`(?i)act\s+as\s+(a|an)\s+`),
	regexp.MustCompile(`(?i)actua\s+como\s+`),
	regexp.MustCompile(`(?i)you\s+are\s+now\s+`),
	regexp.MustCompile(`(?i)ahora\s+eres\s+`),
	regexp.MustCompile(`(?i)system\s*prompt`),
	regexp.MustCompile(`(?i)new\s+instructions?`),
	regexp.MustCompile(`(?i)nuevas?\s+instrucciones`),
	regexp.MustCompile(`(?i)override\s+(the\s+)?(system|rules)`),
	regexp.MustCompile(`(?i)bypass\s+(the\s+)?(filter|security|validation)`),
	regexp.MustCompile(`(?i)saltar\s+(el\s+)?(filtro|seguridad|validacion)`),
	regexp.MustCompile(`(?i)disable\s+(the\s+)?`),
	regexp.MustCompile(`(?i)desactivar\s+`),
	regexp.MustCompile(`(?i)\bDAN\b`),
	regexp.MustCompile(`(?i)do\s+anything\s+now`),
	regexp.MustCompile(`(?i)jailbreak`),
	regexp.MustCompile(`(?i)pretend\s+`),
	regexp.MustCompile(`(?i)simula\s+`),
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

	// Validate table whitelist
	if err := validateTableWhitelist(upper); err != nil {
		return err
	}

	return nil
}

// ValidateUserInput checks if the user's natural language query contains jailbreak attempts.
func ValidateUserInput(query string) error {
	if query == "" {
		return fmt.Errorf("consulta vacia")
	}

	// Max length check
	if len(query) > 500 {
		return fmt.Errorf("consulta demasiado larga (max 500 caracteres)")
	}

	// Jailbreak detection
	for _, pattern := range jailbreakPatterns {
		if pattern.MatchString(query) {
			return fmt.Errorf("consulta no permitida: posible intento de manipulacion")
		}
	}

	return nil
}

// validateTableWhitelist ensures only allowed tables are referenced.
func validateTableWhitelist(upperSQL string) error {
	// Extract table names from FROM and JOIN clauses
	fromPattern := regexp.MustCompile(`\b(?:FROM|JOIN)\s+(\w+)`)
	matches := fromPattern.FindAllStringSubmatch(upperSQL, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		tableName := strings.ToLower(match[1])

		// Skip common SQL keywords that might appear after FROM/JOIN
		if tableName == "select" || tableName == "where" || tableName == "on" {
			continue
		}

		allowed := false
		for _, t := range allowedTables {
			if tableName == t {
				allowed = true
				break
			}
		}
		if !allowed {
			// Check aliases (single letters like v, p, vi, c are ok)
			if len(tableName) <= 3 {
				continue // likely an alias
			}
			return fmt.Errorf("tabla no permitida: %s", tableName)
		}
	}

	return nil
}

