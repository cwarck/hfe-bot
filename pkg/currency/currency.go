package currency

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// AmountRegexp is the regular expression for parsing the amount and currency from the message.
var AmountRegexp = regexp.MustCompile(`^(\d+(?:[.,]\d+)?)\s*([A-Za-z]{3})?$`)

// ParseFromMessage parses the amount and currency from the message.
func ParseFromMessage(text string, defaultCurrency string) (float64, string, error) {
	// Assume the amount is in the default currency if the message is a number.
	if amount, err := strconv.Atoi(text); err == nil {
		return float64(amount), defaultCurrency, nil
	}

	items := strings.Split(text, " ")
	if len(items) != 2 {
		return 0.0, "", fmt.Errorf("unable to parse amount and currency from the message")
	}

	amount, err := strconv.ParseFloat(items[0], 64)
	if err != nil {
		return 0.0, "", fmt.Errorf("invalid amount format: %w", err)
	}
	currency := strings.ToUpper(items[1])

	return amount, currency, nil
}

// GetAmount gets the amount in the default currency.
func GetAmount(text string, oxrAppId string, defaultCurrency string) (float64, error) {
	var result float64

	amount, currency, err := ParseFromMessage(text, defaultCurrency)
	if err != nil {
		return 0.0, err
	}

	rates, err := GetRates(oxrAppId)
	if err != nil {
		return 0.0, fmt.Errorf("failed to get exchange rates: %w", err)
	}

	if currency != rates.Base {
		// Openexchangerates API free tier provides rates in USD, so we have to normalize the amount to USD first.
		result, err = ConvertToBase(amount, rates.Rates, currency)
		if err != nil {
			return 0.0, err
		}
	}
	result, err = ConvertToDefault(result, rates.Rates, defaultCurrency)
	if err != nil {
		return 0.0, err
	}

	return math.Round(result), nil
}

// ConvertToBase converts the given amount to the base currency. first.
func ConvertToBase(amount float64, rates map[string]float64, currency string) (float64, error) {
	if rate, ok := rates[currency]; ok {
		return amount / rate, nil
	}
	return 0.0, fmt.Errorf("unsupported currency: %s", currency)
}

// ConvertToDefault converts the given amount to the default currency.
func ConvertToDefault(amount float64, rates map[string]float64, defaultCurrency string) (float64, error) {
	if rate, ok := rates[defaultCurrency]; ok {
		return amount * rate, nil
	}
	return 0.0, fmt.Errorf("unsupported default currency: %s", defaultCurrency)
}
