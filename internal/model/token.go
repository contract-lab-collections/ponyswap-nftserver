package model

type AssetToken struct {
	Model
	TokenName string           `gorm:"size:20; not null;" json:"name"`
	Symbol    string           `gorm:"size:10; not null;" json:"symbol"`
	Decimals  uint8            `gorm:"size:255;not null;" json:"decimals"`
	Address   string           `gorm:"size:255;not null;index" json:"address"`
	ChainID   string           `gorm:"size:255;not null;" json:"chainId"`
	Icon      string           `gorm:"size:255;" json:"icon"`
	Stable    bool             `gorm:"" json:"stable"`
	Verify    bool             `gorm:"" json:"verify"`
	Stats     *AssetTokenStats `json:"stats,omitempty"`
}

type AssetTokenStats struct {
	Model
	AssetTokenID       uint64  `json:"assetTokenId"`
	LiquidityToken     float64 `json:"liquidityToken"`
	LiquidityUSD       float64 `json:"liquidityUSD"`
	LiquidityUSDChange float64 `json:"liquidityUSDChange"`
	PriceUSD           float64 `json:"priceUSD"`
	PriceUSDChange     float64 `json:"priceUSDChange"`
	PriceUSDChangeWeek float64 `json:"priceUSDChangeWeek"`
	VolumeUSD          float64 `json:"volumeUSD"`
	VolumeUSDChange    float64 `json:"volumeUSDChange"`
	VolumeUSDWeek      float64 `json:"volumeUSDWeek"`
	TxCount            int     `json:"txCount"`
}
