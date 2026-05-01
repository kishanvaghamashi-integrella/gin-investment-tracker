package dto

type CreateUserAssetRequest struct {
	AssetID int64 `json:"asset_id,omitempty" binding:"required,gt=0"`
}
