package constant

type OracleField string

const (
	FieldID         = "ID"
	FieldAppID      = "AppID"
	FieldCoinTypeID = "CoinTypeID"

	RewardFieldDailyReward = "DailyReward"

	CurrencyFieldPriceVSUSDT    = "PriceVSUSDT"
	CurrencyFieldAppPriceVSUSDT = "AppPriceVSUSDT"
	CurrencyFieldOverPercent    = "OverPercent"
	CurrencyFieldCurrencyMethod = "CurrencyMethod"

	CurrencyFixAmount   = "fix-amount"
	CurrencyOverPercent = "over-percent"
)
