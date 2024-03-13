package boutique

import (
	"context"
	"github.com/eniac/mucache/pkg/invoke"
	"math/rand"
)

func min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func GetRecommendations(ctx context.Context, productIds []string) []string {
	_productIds := make([]string, 0)
	req1 := FetchCatalogRequest{}
	catalog := invoke.Invoke[FetchCatalogResponse](ctx, "productcatalog", "ro_fetch_catalog", req1)
	for _, x := range catalog.Catalog {
		_productIds = append(_productIds, x.Id)
	}
	filteredProducts := make([]string, 0, len(_productIds)-len(productIds))
	for _, id := range _productIds {
		if contains(productIds, id) == false {
			filteredProducts = append(filteredProducts, id)
		}
	}
	numProducts := len(filteredProducts)
	numReturn := min(5, numProducts)
	// sample list of indices to return
	indices := rand.Perm(numProducts)[:numReturn]
	// fetch product ids from indices
	prodList := make([]string, numReturn)
	for i, idx := range indices {
		prodList[i] = filteredProducts[idx]
	}
	return prodList
}
