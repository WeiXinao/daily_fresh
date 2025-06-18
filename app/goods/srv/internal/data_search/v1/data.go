package search

type SearchFactory interface {
	Goods() GoodsStore
}