package boutique

import (
	"context"
	"fmt"
	"github.com/eniac/mucache/pkg/state"
	"math"
	"strconv"
	"strings"
)

func _carry(units float64, nanos float64) carry {
	fractionSize := math.Pow(10, 9)
	nanos += fractionSize
	_units := math.Floor(units) + math.Floor(nanos/fractionSize)
	_nanos := float64(int64(math.Round(nanos)) % int64(fractionSize))
	return carry{Units: _units, Nanos: _nanos}
}

func _MoneyToString(m Money) string {
	nanosStr := strconv.FormatInt(int64(m.Nanos), 10)
	nanosStr = strings.Repeat("0", 9-len(nanosStr)) + nanosStr
	return fmt.Sprintf("%v.%v %v", m.Units, nanosStr, m.Currency)
}

func SetCurrencySupport(ctx context.Context, currency Currency) bool {
	state.SetState(ctx, currency.CurrencyCode, currency)
	return true
}

func InitCurrencies(ctx context.Context, currencies []Currency) {
	currencyCodes := make([]string, len(currencies))
	for i, currency := range currencies {
		currencyCodes[i] = currency.CurrencyCode
		SetCurrencySupport(ctx, currency)
	}
	state.SetState(ctx, "CURRENCIES", currencyCodes)
}

func ConvertCurrency(ctx context.Context, amount Money, toCurrency string) Money {
	fromRate, err := state.GetState[Currency](ctx, amount.Currency)
	if err != nil {
		panic(err)
	}

	fromRate64, err := strconv.ParseFloat(fromRate.Rate, 64)
	if err != nil {
		panic(err)
	}

	// Convert: from_currency --> EUR
	euros := _carry(float64(amount.Units)/fromRate64, float64(amount.Nanos)/fromRate64)

	euros.Nanos = math.Round(euros.Nanos)

	// Convert: EUR --> to_currency
	toRate, err := state.GetState[Currency](ctx, toCurrency)
	if err != nil {
		panic(err)
	}

	toRate64, err := strconv.ParseFloat(toRate.Rate, 64)
	if err != nil {
		panic(err)
	}
	_result := _carry(euros.Units*toRate64, euros.Nanos*toRate64)

	_result.Units = math.Floor(_result.Units)
	_result.Nanos = math.Floor(_result.Nanos)
	result := Money{Currency: toCurrency, Units: int32(_result.Units), Nanos: int64(_result.Nanos)}

	return result
}

func GetSupportedCurrencies(ctx context.Context) []Currency {
	keys, err := state.GetState[[]string](ctx, "CURRENCIES")
	if err != nil {
		panic(err)
	}

	// Bulk
	var currencies []Currency
	if len(keys) > 0 {
		currencies = state.GetBulkStateDefault[Currency](ctx, keys, Currency{})
	} else {
		currencies = make([]Currency, len(keys))
	}
	//for _, currencyCode := range keys {
	//	currency, err := state.GetState[Currency](ctx, currencyCode)
	//	if err != nil {
	//		panic(err)
	//	}
	//	currencies = append(currencies, currency)
	//}
	return currencies
}
