package handlers

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/application/use_cases"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/i18n"
)

var stockDetails = regexp.MustCompile(`\(disponible: ([0-9.]+), solicitado: ([0-9.]+)\)$`)

// localizedError translates known domain failures without changing business data.
// Unknown diagnostics remain visible rather than being mislabeled as validation errors.
func localizedError(r *http.Request, err error) string {
	for _, entry := range []struct {
		err error
		key string
	}{
		{entities.ErrProductNameEmpty, "error.product.name"},
		{entities.ErrClientNameRequired, "error.client.name"},
		{entities.ErrProductPriceInvalid, "error.product.price"},
		{entities.ErrProductCostNegative, "error.product.cost"},
		{entities.ErrProductStockNegative, "error.product.stock"},
		{use_cases.ErrProductNotFound, "error.product.not_found"},
		{use_cases.ErrSaleNoItemsProvided, "sales.error.empty"},
		{entities.ErrSaleNoItems, "sales.error.empty"},
		{use_cases.ErrSaleInvalidQuantity, "sales.error.quantity"},
		{entities.ErrSaleItemQtyInvalid, "sales.error.quantity"},
		{entities.ErrSaleInvalidPayment, "sales.error.payment"},
		{entities.ErrSaleInvalidTotal, "sales.error.total"},
		{entities.ErrSaleItemPriceInvalid, "sales.error.item_price"},
		{use_cases.ErrSaleProductNotFound, "error.product.not_found"},
		{use_cases.ErrSaleInsufficientStock, "sales.error.stock"},
	} {
		if errors.Is(err, entry.err) {
			message := tr(r, entry.key)
			// Preserve item IDs, names and quantities accompanying stock failures.
			if entry.err == use_cases.ErrSaleProductNotFound || entry.err == use_cases.ErrSaleInsufficientStock {
				_, detail, found := strings.Cut(err.Error(), entry.err.Error())
				if found {
					if i18n.FromRequest(r).IsEN() {
						detail = stockDetails.ReplaceAllString(detail, "(available: $1, requested: $2)")
					}
					message += detail
				}
			}
			return message
		}
	}
	return formatUseCaseError(err)
}
