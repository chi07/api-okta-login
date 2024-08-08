package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CurrencyConfig struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Currency     string             `bson:"currency" json:"currency"`
	CurrencyName string             `bson:"currency_name" json:"currency_name"`

	Binance       bool   `bson:"binance" json:"binance"`
	BinanceSymbol string `bson:"binance_symbol" json:"binance_symbol"`

	GateIO       bool   `bson:"gateio" json:"gateio"`
	GateIOSymbol string `bson:"gateio_symbol" json:"gateio_symbol"`

	Okx       bool   `bson:"okx" json:"okx"`
	OkxSymbol string `bson:"okx_symbol" json:"okx_symbol"`

	Coinbase       bool   `bson:"coinbase" json:"coinbase"`
	CoinbaseSymbol string `bson:"coinbase_symbol" json:"coinbase_symbol"`

	Bitget       bool   `bson:"bitget" json:"bitget"`
	BitgetSymbol string `bson:"bitget_symbol" json:"bitget_symbol"`

	Bybit       bool   `bson:"bybit" json:"bybit"`
	BybitSymbol string `bson:"bybit_symbol" json:"bybit_symbol"`

	IsExclusive   bool `bson:"is_exclusive" json:"is_exclusive"`
	MinForWarning int  `bson:"min_for_warning" json:"min_for_warning"`

	CreatedAt time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
