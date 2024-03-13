package boutique

import (
	"context"
	"fmt"
	"github.com/eniac/mucache/pkg/state"
)

func remove(slice []int, s int) []int {
	return append(slice[:s], slice[s+1:]...)
}

func getCartDefault(ctx context.Context, userId string) Cart {
	cart, err := state.GetState[Cart](ctx, userId)
	if fmt.Sprint(err) == "key not found" {
		cart = Cart{
			UserId: userId,
			Items:  []CartItem{},
		}
	} else if err != nil {
		panic(err)
	}
	return cart
}

func AddItem(ctx context.Context, userId string, productId string, quantity int32) bool {
	item := CartItem{
		ProductId: productId,
		Quantity:  quantity,
	}
	cart := getCartDefault(ctx, userId)

	// Append the new item to the cart
	cart.Items = append(cart.Items, item)
	state.SetState(ctx, userId, cart)
	return true
}

func GetCart(ctx context.Context, userId string) Cart {
	cart, err := state.GetState[Cart](ctx, userId)
	if err != nil {
		cart = Cart{
			UserId: userId,
			Items:  []CartItem{},
		}
	}
	return cart
}

func EmptyCart(ctx context.Context, userId string) bool {
	cart := getCartDefault(ctx, userId)
	cart.Items = []CartItem{}
	state.SetState(ctx, userId, cart)
	return true
}
