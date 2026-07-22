package nlsql

import (
	"fmt"
	"regexp"
	"strings"
)

var dangerousKeywords = []string{
	"INSERT", "UPDATE", "DELETE", "DROP", "ALTER", "CREATE",
	"TRUNCATE", "REPLACE", "GRANT", "REVOKE", "EXEC",
	"EXECUTE", "ATTACH", "DETACH", "PRAGMA", "VACUUM",
}

var allowedTables = []string{
	"productos", "categorias", "ventas", "venta_items",
	"clientes", "usuarios", "inventario_movimientos", "configuracion",
}

var jailbreakPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)ignor[ae]\s+(las\s+)?instrucciones`),
	regexp.MustCompile(`(?i)ignore\s+(all\s+)?(previous\s+)?instructions`),
	regexp.MustCompile(`(?i)forget\s+(all\s+)?(previous|your)`),
	regexp.MustCompile(`(?i)act\s+as\s+`),
	regexp.MustCompile(`(?i)system\s*prompt`),
	regexp.MustCompile(`(?i)bypass\s+(the\s+)?(filter|security)`),
	regexp.MustCompile(`(?i)\bDAN\b`),
	regexp.MustCompile(`(?i)jailbreak`),
}

// ValidateUserInput checks for jailbreak attempts and length.
func ValidateUserInput(query string) error {
	if query == "" {
		return fmt.Errorf("consulta vacia")
	}
	if len(query) > 500 {
		return fmt.Errorf("consulta demasiado larga (max 500 caracteres)")
	}
	for _, p := range jailbreakPatterns {
		if p.MatchString(query) {
			return fmt.Errorf("consulta no permitida")
		}
	}
	return nil
}

// ValidateSQL checks if a generated SQL query is safe to execute.
func ValidateSQL(sql string) error {
	if sql == "" {
		return fmt.Errorf("consulta SQL vacia")
	}
	normalized := strings.TrimSpace(sql)
	upper := strings.ToUpper(normalized)

	if !strings.HasPrefix(upper, "SELECT") && !strings.HasPrefix(upper, "WITH") {
		return fmt.Errorf("solo se permiten consultas SELECT")
	}
	for _, kw := range dangerousKeywords {
		pattern := fmt.Sprintf(`\b%s\b`, kw)
		if matched, _ := regexp.MatchString(pattern, upper); matched {
			return fmt.Errorf("palabra clave no permitida: %s", kw)
		}
	}
	if strings.Contains(normalized, "--") || strings.Contains(normalized, "/*") {
		return fmt.Errorf("no se permiten comentarios SQL")
	}
	if matched, _ := regexp.MatchString(`;[\s]*\S`, normalized); matched {
		return fmt.Errorf("no se permiten multiples sentencias")
	}
	// Table whitelist
	fromPattern := regexp.MustCompile(`\b(?:FROM|JOIN)\s+(\w+)`)
	matches := fromPattern.FindAllStringSubmatch(upper, -1)
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		tbl := strings.ToLower(m[1])
		if len(tbl) <= 3 {
			continue // alias
		}
		found := false
		for _, a := range allowedTables {
			if tbl == a {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("tabla no permitida: %s", tbl)
		}
	}
	return nil
}
