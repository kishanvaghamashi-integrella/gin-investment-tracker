package casparsermodel

type StatementPeriod struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Valuation struct {
	Date  string `json:"date"`
	NAV   string `json:"nav"`
	Cost  string `json:"cost"`
	Value string `json:"value"`
}

type Transaction struct {
	Date         string  `json:"date"`
	Description  string  `json:"description"`
	Amount       string  `json:"amount"`
	Units        string  `json:"units"`
	NAV          string  `json:"nav"`
	Balance      string  `json:"balance"`
	Type         string  `json:"type"`
	DividendRate *string `json:"dividend_rate"`
}

type Scheme struct {
	Scheme          string        `json:"scheme"`
	Advisor         string        `json:"advisor"`
	RTACode         string        `json:"rta_code"`
	RTA             string        `json:"rta"`
	Type            string        `json:"type"`
	ISIN            string        `json:"isin"`
	AMFI            string        `json:"amfi"`
	Nominees        []string      `json:"nominees"`
	Open            string        `json:"open"`
	Close           string        `json:"close"`
	CloseCalculated string        `json:"close_calculated"`
	Valuation       Valuation     `json:"valuation"`
	Transactions    []Transaction `json:"transactions"`
}

type Folio struct {
	Folio   string   `json:"folio"`
	AMC     string   `json:"amc"`
	PAN     string   `json:"PAN"`
	KYC     string   `json:"KYC"`
	PANKYC  string   `json:"PANKYC"`
	Schemes []Scheme `json:"schemes"`
}

type InvestorInfo struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Mobile  string `json:"mobile"`
}

type CASStatement struct {
	StatementPeriod StatementPeriod `json:"statement_period"`
	Folios          []Folio         `json:"folios"`
	InvestorInfo    InvestorInfo    `json:"investor_info"`
	CASType         string          `json:"cas_type"`
	FileType        string          `json:"file_type"`
}
