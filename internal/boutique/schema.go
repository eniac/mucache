package boutique

type Money struct {
	Currency string `json:"currencyCode"`
	Units    int32  `json:"units"`
	Nanos    int64  `json:"nanos"`
}

type carry struct {
	Units float64 `json:"units"`
	Nanos float64 `json:"nanos"`
}

type Currency struct {
	CurrencyCode string `json:"currencyCode"`
	Rate         string `json:"rate"`
}

type Address struct {
	StreetAddress string `json:"street_address"`
	City          string `json:"city"`
	State         string `json:"state"`
	Country       string `json:"country"`
	ZipCode       int32  `json:"zip_code"`
}

type CartItem struct {
	ProductId string `json:"product_id"`
	Quantity  int32  `json:"quantity"`
}

type Cart struct {
	UserId string     `json:"user_id"`
	Items  []CartItem `json:"cart"`
}

type Product struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Picture     string   `json:"picture"`
	PriceUsd    *Money   `json:"priceUSD"`
	Categories  []string `json:"categories"`
}

type OrderItem struct {
	Item *CartItem `json:"item"`
	Cost *Money    `json:"cost"`
}

/*
 * Frontend
 */

type HomeRequest struct {
	Userid      string `json:"user_id"`
	CatalogSize int    `json:"catalog_size"`
}

type HomeResponse struct {
	Products   []Product  `json:"products"`
	UserCart   Cart       `json:"userCart"`
	Currencies []Currency `json:"currencies"`
}

type FrontendSetCurrencyRequest struct {
	Cur Currency `json:"cur"`
}

type FrontendSetCurrencyResponse struct {
	OK string `json:"ok"`
}

type BrowseProductRequest struct {
	ProductId string `json:"product_id"`
}

type BrowseProductResponse struct {
	Prod Product `json:"prod"`
}

type AddToCartRequest struct {
	UserId    string `json:"user_id"`
	ProductId string `json:"product_id"`
	Quantity  int32  `json:"quantity"`
}

type AddToCartResponse struct {
	OK bool `json:"ok"`
}

type ViewCartRequest struct {
	UserId string `json:"user_id"`
}

type ViewCartResponse struct {
	C Cart `json:"c"`
}

type CheckoutRequest struct {
	UserId       string     `json:"user_id"`
	UserCurrency string     `json:"user_currency"`
	Address      Address    `json:"address"`
	Email        string     `json:"email"`
	CreditCard   CreditCard `json:"credit_card"`
}

type CheckoutResponse struct {
	Res OrderResult `json:"res"`
}

/*
 *	Product Catalog structs
 */

type AddProductRequest struct {
	Product Product `json:"product"`
}

type AddProductResponse struct {
	ProductId string `json:"product_id"`
}

type AddProductsRequest struct {
	Products []Product `json:"products"`
}

type AddProductsResponse struct {
	OK string `json:"ok"`
}

type GetProductRequest struct {
	ProductId string `json:"id"`
}

type GetProductResponse struct {
	Product Product `json:"product"`
}

type SearchProductsRequest struct {
	Query string `json:"name"`
}

type SearchProductsResponse struct {
	Products []Product `json:"products"`
}

type FetchCatalogRequest struct {
	CatalogSize int `json:"catalog_size"`
}

type FetchCatalogResponse struct {
	Catalog []Product `json:"catalog"`
}

/*
 *	Recommendations structs
 */

type GetRecommendationsRequest struct {
	ProductIds []string `json:"product_ids"`
}

type GetRecommendationsResponse struct {
	ProductIds []string `json:"product_ids"`
}

/*
 *	Currency structs
 */

type SetCurrencySupportRequest struct {
	Currency Currency `json:"currency"`
}

type SetCurrencySupportResponse struct {
	Ok bool `json:"ok"`
}

type ConvertCurrencyRequest struct {
	Amount     Money  `json:"amount"`
	ToCurrency string `json:"to_currency"`
}

type ConvertCurrencyResponse struct {
	Amount Money `json:"amount"`
}

type GetSupportedCurrenciesRequest struct {
}

type GetSupportedCurrenciesResponse struct {
	Currencies []Currency `json:"currencies"`
}

type InitCurrencyRequest struct {
	Currencies []Currency `json:"currencies"`
}

type InitCurrencyResponse struct {
	Ok string `json:"ok"`
}

/*
 *	Shipping structs
 */

type Quote struct {
	Dollars int32
	Cents   int64
}

type GetQuoteRequest struct {
	Items []CartItem `json:"items"`
}

type GetQuoteResponse struct {
	CostUsd Money `json:"cost_usd"`
}

type ShipOrderRequest struct {
	Address Address    `json:"address"`
	Items   []CartItem `json:"items"`
}

type ShipOrderResponse struct {
	TrackingId string `json:"tracking_id"`
}

/*
 *	Cart structs
 */

type AddItemRequest struct {
	UserId    string `json:"user_id"`
	ProductId string `json:"product_id"`
	Quantity  int32  `json:"quantity"`
}

type AddItemResponse struct {
	Ok bool `json:"ok"`
}

type GetCartRequest struct {
	UserId string `json:"user_id"`
}

type GetCartResponse struct {
	Cart Cart `json:"cart"`
}

type EmptyCartRequest struct {
	UserId string `json:"user_id"`
}

type EmptyCartResponse struct {
	Ok bool `json:"ok"`
}

/*
 *	Payment struct
 */

type CreditCard struct {
	CardNumber      string `json:"card_number"`
	CardType        string `json:"card_type"`
	ExpirationMonth int    `json:"expiration_month"`
	ExpirationYear  int    `json:"expiration_year"`
}

type ChargeRequest struct {
	Amount     Money      `json:"amount"`
	CreditCard CreditCard `json:"credit_cart"`
}

type ChargeResponse struct {
	Uuid  string `json:"uuid"`
	Error string `json:"error"`
}

/*
 * Checkout struct
 */

type OrderResult struct {
	OrderId            string      `json:"order_id"`
	ShippingTrackingId string      `json:"shipping_tracking_id"`
	ShippingCost       *Money      `json:"shipping_cost"`
	ShippingAddress    Address     `json:"shipping_address"`
	Items              []OrderItem `json:"items"`
}

type PlaceOrderRequest struct {
	UserId       string     `json:"user_id"`
	UserCurrency string     `json:"user_currency"`
	Address      Address    `json:"address"`
	Email        string     `json:"email"`
	CreditCard   CreditCard `json:"credit_card"`
}

type PlaceOrderResponse struct {
	Order OrderResult `json:"order"`
}

/*
 * Email structs
 */

type SendOrderConfirmationRequest struct {
	Email string      `json:"email"`
	Order OrderResult `json:"order"`
}

type SendOrderConfirmationResponse struct {
	Ok bool `json:"ok"`
}
