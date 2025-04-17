package format

import (
	"fmt"
	"github.com/leekchan/accounting"
	"github.com/shopspring/decimal"
)

var defaultFormat = accounting.Accounting{
	Symbol:         "THB",
	Precision:      2,
	Thousand:       ",",
	Decimal:        ".",
	Format:         "%s %v ",
	FormatNegative: "(%s %v)",
}

func DecimalToMoney(amount decimal.Decimal) string {
	return defaultFormat.FormatMoneyDecimal(amount)
}

func NullDecimalToMoney(amount decimal.NullDecimal) string {
	var amt decimal.Decimal
	if amount.Valid {
		amt = amount.Decimal
	} else {
		amt = decimal.Zero
	}
	return defaultFormat.FormatMoneyDecimal(amt)
}

func Percent(amount decimal.Decimal) string {
	return fmt.Sprintf("%s %%", amount.StringFixedBank(2))
}
