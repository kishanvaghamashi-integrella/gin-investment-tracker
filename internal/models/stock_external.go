package model

type StockChartResult struct {
	Chart StockChart `json:"chart"`
}

type StockChart struct {
	Result []StockResult `json:"result"`
	Error  interface{}   `json:"error"`
}

type StockResult struct {
	Meta StockMeta `json:"meta"`
}

type StockMeta struct {
	ChartPreviousClose float64 `json:"chartPreviousClose"`
}
