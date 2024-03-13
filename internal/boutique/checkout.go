package boutique

import (
	"context"
	"github.com/eniac/mucache/pkg/invoke"
	"github.com/google/uuid"
)

const (
	nanosMod = 1000000000
)

func Sum(l, r Money) Money {
	units := l.Units + r.Units
	nanos := l.Nanos + r.Nanos

	if (units == 0 && nanos == 0) || (units > 0 && nanos >= 0) || (units < 0 && nanos <= 0) {
		// same sign <units, nanos>
		units += int32(nanos / nanosMod)
		nanos = nanos % nanosMod
	} else {
		// different sign. nanos guaranteed to not to go over the limit
		if units > 0 {
			units--
			nanos += nanosMod
		} else {
			units++
			nanos -= nanosMod
		}
	}

	return Money{
		Units:    units,
		Nanos:    nanos,
		Currency: l.Currency}
}

func MultiplySlow(m Money, n uint32) Money {
	out := m
	for n > 1 {
		out = Sum(out, m)
		n--
	}
	return out
}

type orderPrep struct {
	orderItems            []OrderItem
	cartItems             []CartItem
	shippingCostLocalized Money
}

func PlaceOrder(ctx context.Context, userId string, userCurrency string, address Address, email string, creditCard CreditCard) OrderResult {

	orderID := uuid.New().String()

	prep := prepareOrderItemsAndShippingQuoteFromCart(ctx, userId, userCurrency, address)

	total := Money{
		Currency: userCurrency,
		Units:    0,
		Nanos:    0,
	}
	total = Sum(total, prep.shippingCostLocalized)
	for _, it := range prep.orderItems {
		multPrice := MultiplySlow(*it.Cost, uint32(it.Item.Quantity))
		total = Sum(total, multPrice)
	}

	chargeCard(ctx, total, creditCard)
	//fmt.Println("payment went through with txID %s", txID)
	shippingTrackingID := shipOrder(ctx, address, prep.cartItems)

	_ = emptyUserCart(ctx, userId)

	orderResult := OrderResult{
		OrderId:            orderID,
		ShippingTrackingId: shippingTrackingID,
		ShippingCost:       &prep.shippingCostLocalized,
		ShippingAddress:    address,
		Items:              prep.orderItems,
	}

	//sendOrderConfirmation(ctx, req.Email, orderResult)

	return orderResult
}

func prepOrderItems(ctx context.Context, items []CartItem, userCurrency string) []OrderItem {
	out := make([]OrderItem, len(items))
	for i, item := range items {
		req1 := GetProductRequest{ProductId: item.ProductId}
		product := invoke.Invoke[GetProductResponse](ctx, "productcatalog", "ro_get_product", req1)
		req2 := ConvertCurrencyRequest{Amount: *product.Product.PriceUsd, ToCurrency: userCurrency}
		price := invoke.Invoke[ConvertCurrencyResponse](ctx, "currency", "ro_convert_currency", req2)
		out[i] = OrderItem{
			Item: &item,
			Cost: &price.Amount}
	}
	return out
}

func prepareOrderItemsAndShippingQuoteFromCart(ctx context.Context, userID, userCurrency string, address Address) orderPrep {
	var out orderPrep
	cartItems := getUserCart(ctx, userID)
	orderItems := prepOrderItems(ctx, cartItems.Items, userCurrency)
	shippingUSD := quoteShipping(ctx, address, cartItems.Items)
	shippingPrice := convertCurrency(ctx, shippingUSD, userCurrency)
	out.shippingCostLocalized = shippingPrice
	out.cartItems = cartItems.Items
	out.orderItems = orderItems
	return out
}

func quoteShipping(ctx context.Context, address Address, items []CartItem) Money {
	req1 := GetQuoteRequest{Items: items}
	shippingQuote := invoke.Invoke[GetQuoteResponse](ctx, "shipping", "ro_get_quote", req1)
	return shippingQuote.CostUsd
}

func getUserCart(ctx context.Context, userID string) Cart {
	req1 := GetCartRequest{UserId: userID}
	cart := invoke.Invoke[GetCartResponse](ctx, "cart", "ro_get_cart", req1)
	return cart.Cart
}

func emptyUserCart(ctx context.Context, userID string) bool {
	req1 := EmptyCartRequest{UserId: userID}
	res := invoke.Invoke[EmptyCartResponse](ctx, "cart", "empty_cart", req1)
	return res.Ok
}

func convertCurrency(ctx context.Context, from Money, toCurrency string) Money {
	req1 := ConvertCurrencyRequest{Amount: from, ToCurrency: toCurrency}
	result := invoke.Invoke[ConvertCurrencyResponse](ctx, "currency", "ro_convert_currency", req1)
	return result.Amount
}

func chargeCard(ctx context.Context, amount Money, paymentInfo CreditCard) string {
	req1 := ChargeRequest{Amount: amount, CreditCard: paymentInfo}
	paymentResp := invoke.Invoke[ChargeResponse](ctx, "payment", "charge", req1)
	if paymentResp.Error != "" {
		panic(paymentResp.Error)
	}
	return paymentResp.Uuid
}

func sendOrderConfirmation(ctx context.Context, email string, order OrderResult) bool {
	req1 := SendOrderConfirmationRequest{
		Email: email,
		Order: order,
	}
	resp := invoke.Invoke[SendOrderConfirmationResponse](ctx, "email", "ro_send_email", req1)
	return resp.Ok
}

func shipOrder(ctx context.Context, address Address, items []CartItem) string {
	req1 := ShipOrderRequest{Address: address, Items: items}
	resp := invoke.Invoke[ShipOrderResponse](ctx, "shipping", "ship_order", req1)
	return resp.TrackingId
}
