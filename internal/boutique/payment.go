package boutique

import (
	"context"
	"github.com/google/uuid"
	"time"
)

func Charge(ctx context.Context, amount Money, creditCard CreditCard) (string, string) {
	if !(creditCard.CardType == "visa" || creditCard.CardType == "mastercard") {
		return "-1", "invalid card type"
	}

	now := time.Now()
	currentMonth := int(now.Month())
	currentYear := now.Year()
	year := creditCard.ExpirationYear
	month := creditCard.ExpirationMonth
	if (currentYear*12 + int(currentMonth)) > (year*12 + month) {
		return "-1", "expired credit card"
	}

	uid := uuid.New().String()

	return uid, ""
}
