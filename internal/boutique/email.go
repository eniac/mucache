package boutique

import (
	"context"
)

func SendConfirmation(ctx context.Context, email string, order OrderResult) bool {
	// fmt.Println("Sent confirmation email to % for order %s", email, order.OrderId)
	return true
}
