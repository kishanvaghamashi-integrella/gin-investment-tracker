package assetprice

type AssetPriceFetcherInterface interface {
	FetchPrice(externalID string) (float64, float64, error)
}
