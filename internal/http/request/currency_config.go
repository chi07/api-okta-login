package request

type ExclusiveCurrency struct {
	Currency    string `json:"currency" validate:"required"`
	IsExclusive bool   `json:"is_exclusive" validate:"required"`
}

type BulkExclusiveCurrency struct {
	Currencies []*ExclusiveCurrency `json:"currencies" validate:"required"`
}
