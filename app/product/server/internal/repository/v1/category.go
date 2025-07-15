package repository

import "github.com/WeiXinao/daily_fresh/app/product/server/internal/repository/v1/dao/category"

type CategoryRepo interface {	
}

var _ CategoryRepo = (*categoryRepo)(nil)

type categoryRepo struct {
	categoryDao category.CategoryDao
}

func NewCategoryRepo(categoryDao category.CategoryDao) CategoryRepo {
	return &categoryRepo{
		categoryDao: categoryDao,
	}
}