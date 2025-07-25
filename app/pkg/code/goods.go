package code

//go:generate codegen -type=int -doc -output ./error_code_generated.md
const (
  // ErrGoodsNotFound - 404: Goods not found.
	ErrGoodsNotFound int = iota + 100501

	// ErrCategoryNotFound - 404: Category not found.
	ErrCategoryNotFound 

	// ErrEsUnmarshal - 500: ES unmarshal error.
	ErrEsUnmarshal

	// ErrElasticSearch - 500: ElasticSearch error.
	ErrElasticsearch
)