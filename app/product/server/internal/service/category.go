package service

import "github.com/WeiXinao/daily_fresh/app/product/server/internal/repository/v1"


type CategoryService interface {

}

type categoryService struct {
	catRepo repository.CategoryRepo
}

func NewCategoryService(catRepo repository.CategoryRepo) CategoryService {
	return &categoryService{
		catRepo: catRepo,
	}
}