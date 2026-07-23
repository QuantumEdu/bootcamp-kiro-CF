package nlsql

import (
	"fmt"
	"strconv"
	"strings"
)

// moneyColumns contains column name patterns that indicate monetary values.
var moneyColumns = []string{
	"precio", "total", "subtotal", "monto", "costo",
	"precio_venta", "precio_compra", "ingreso",
}

// FormatResults converts ChatResult columns/rows into a human-readable Spanish summary.
// Rules:
//   - 0 rows → "No se encontraron resultados."
//   - 1 row, 1 column → return the value directly (formatted if money)
//   - ≤5 rows → numbered list
//   - >5 rows → first 5 with "...y X más"
func FormatResults(columns []string, rows [][]string) string {
	if len(rows) == 0 {
		return "No se encontraron resultados."
	}

	// Single value case
	if len(rows) == 1 && len(columns) == 1 {
		val := rows[0][0]
		if isMoneyColumn(columns[0]) {
			val = formatMoney(val)
		}
		return val
	}

	// Single row, multiple columns → describe as key: value pairs
	if len(rows) == 1 {
		parts := make([]string, 0, len(columns))
		for i, col := range columns {
			val := rows[0][i]
			if isMoneyColumn(col) {
				val = formatMoney(val)
			}
			parts = append(parts, fmt.Sprintf("%s: %s", humanizeColumn(col), val))
		}
		return strings.Join(parts, " | ")
	}

	// Multiple rows
	var sb strings.Builder
	limit := 5
	if len(rows) <= limit {
		limit = len(rows)
	}

	for i := 0; i < limit; i++ {
		sb.WriteString(fmt.Sprintf("%d. ", i+1))
		parts := make([]string, 0, len(columns))
		for j, col := range columns {
			val := rows[i][j]
			if isMoneyColumn(col) {
				val = formatMoney(val)
			}
			parts = append(parts, fmt.Sprintf("%s: %s", humanizeColumn(col), val))
		}
		sb.WriteString(strings.Join(parts, " | "))
		if i < limit-1 || len(rows) > limit {
			sb.WriteString("\n")
		}
	}

	if len(rows) > 5 {
		sb.WriteString(fmt.Sprintf("...y %d más", len(rows)-5))
	}

	return sb.String()
}

// isMoneyColumn returns true if the column name suggests a monetary value.
func isMoneyColumn(col string) bool {
	lower := strings.ToLower(col)
	for _, mc := range moneyColumns {
		if strings.Contains(lower, mc) {
			return true
		}
	}
	return false
}

// formatMoney formats a numeric string as currency ($X,XXX.XX).
func formatMoney(val string) string {
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return val // not a number, return as-is
	}
	// Format with 2 decimals and thousands separator
	negative := f < 0
	if negative {
		f = -f
	}

	intPart := int64(f)
	decPart := int64((f - float64(intPart)) * 100 + 0.5)

	intStr := formatWithThousands(intPart)

	result := fmt.Sprintf("$%s.%02d", intStr, decPart)
	if negative {
		result = "-" + result
	}
	return result
}

// formatWithThousands adds comma separators to an integer.
func formatWithThousands(n int64) string {
	s := strconv.FormatInt(n, 10)
	if len(s) <= 3 {
		return s
	}
	var sb strings.Builder
	remainder := len(s) % 3
	if remainder > 0 {
		sb.WriteString(s[:remainder])
		if len(s) > remainder {
			sb.WriteString(",")
		}
	}
	for i := remainder; i < len(s); i += 3 {
		sb.WriteString(s[i : i+3])
		if i+3 < len(s) {
			sb.WriteString(",")
		}
	}
	return sb.String()
}

// humanizeColumn converts a snake_case column name to a more readable format.
func humanizeColumn(col string) string {
	return strings.ReplaceAll(col, "_", " ")
}
