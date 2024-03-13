package main

import (
	"context"
	"fmt"
	"github.com/eniac/mucache/internal/boutique"
	"github.com/eniac/mucache/pkg/cm"
	"github.com/eniac/mucache/pkg/wrappers"
	"net/http"
	"runtime"
)

func heartbeat(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Heartbeat\n"))
	if err != nil {
		return
	}
}

func addProduct(ctx context.Context, req *boutique.AddProductRequest) *boutique.AddProductResponse {
	productId := boutique.AddProduct(ctx, req.Product)
	resp := boutique.AddProductResponse{ProductId: productId}
	return &resp
}

func getProduct(ctx context.Context, req *boutique.GetProductRequest) *boutique.GetProductResponse {
	product := boutique.GetProduct(ctx, req.ProductId)
	//fmt.Printf("Product read: %+v\n", product)
	resp := boutique.GetProductResponse{Product: product}
	return &resp
}

func searchProducts(ctx context.Context, req *boutique.SearchProductsRequest) *boutique.SearchProductsResponse {
	products := boutique.SearchProducts(ctx, req.Query)
	//fmt.Printf("Products read: %+v\n", products)
	resp := boutique.SearchProductsResponse{Products: products}
	return &resp
}

func fetchCatalog(ctx context.Context, req *boutique.FetchCatalogRequest) *boutique.FetchCatalogResponse {
	products := boutique.FetchCatalog(ctx, req.CatalogSize)
	resp := boutique.FetchCatalogResponse{Catalog: products}
	return &resp
}

func addProducts(ctx context.Context, req *boutique.AddProductsRequest) *boutique.AddProductsResponse {
	boutique.AddProducts(ctx, req.Products)
	resp := boutique.AddProductsResponse{OK: "OK"}
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	go cm.ZmqProxy()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/add_product", wrappers.NonROWrapper[boutique.AddProductRequest, boutique.AddProductResponse](addProduct))
	http.HandleFunc("/add_products", wrappers.NonROWrapper[boutique.AddProductsRequest, boutique.AddProductsResponse](addProducts))
	http.HandleFunc("/ro_get_product", wrappers.ROWrapper[boutique.GetProductRequest, boutique.GetProductResponse](getProduct))
	http.HandleFunc("/ro_search_products", wrappers.ROWrapper[boutique.SearchProductsRequest, boutique.SearchProductsResponse](searchProducts))
	http.HandleFunc("/ro_fetch_catalog", wrappers.ROWrapper[boutique.FetchCatalogRequest, boutique.FetchCatalogResponse](fetchCatalog))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
