package dto

type CreateUserAssetRequest struct {
	AssetID int64 `json:"asset_id" binding:"required,gt=0"`
}
