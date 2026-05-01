package dto

type DashboardDataDto struct {
	CurrentInvestmentValue     *float64                 `json:"current_investment_value,omitempty"`
	PreviousDayInvestmentValue *float64                 `json:"previous_day_investment_value,omitempty"`
	TotalInvestedValue         *float64                 `json:"total_invested_value,omitempty"`
	TopHoldings                []HoldingResponseDto     `json:"top_holdings,omitempty"`
	RecentTransactions         []TransactionResponseDto `json:"recent_transactions,omitempty"`
}
