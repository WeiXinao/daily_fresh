package code

const (
	// ErrOrderNotFound - 404: Order not found.
	ErrOrderNotFound int = iota + 100801

	// ErrShopCartItemNotFound - 404: ShopCart item not found.
	ErrShopCartItemNotFound

	// ErrSubmitOrder - 500: Submit order error.
	ErrSubmitOrder

	// ErrNoGoodsSelected - 404: No goods selected.
	ErrNoGoodsSelected
)
