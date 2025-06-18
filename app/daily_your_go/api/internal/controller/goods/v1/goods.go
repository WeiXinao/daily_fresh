package goods

import (
	proto "github.com/WeiXinao/daily_your_go/api/goods/v1"
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/domain/request"
	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/service"
	"github.com/WeiXinao/daily_your_go/app/pkg/translator/ginx"
	"github.com/WeiXinao/daily_your_go/pkg/common/core"
	"github.com/WeiXinao/daily_your_go/pkg/log"
	"github.com/WeiXinao/xkit/slice"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
)

type GoodsController struct {
	srv service.ServiceFactory
	ut ut.Translator
}

func (g *GoodsController) List(ctx *gin.Context) {
	log.Info("goods list  function called")

	var r request.GoodsFilter
	if err := ctx.ShouldBindQuery(&r); err != nil {
		ginx.HandleValidatorError(ctx, g.ut, err)
		return
	}

	gfr := &proto.GoodsFilterRequest{
		PriceMin:    r.PriceMin,
		PriceMax:    r.PriceMax,
		IsHot:       r.IsHot,
		IsNew:       r.IsNew,
		IsTab:       r.IsTab,
		TopCategory: r.TopCategory,
		Pages:       r.Pages,
		PagePerNums: r.PagePerNums,
		KeyWords:    r.KeyWords,
		Brand:       r.Brand,
	}
	glr, err := g.srv.Goods().List(ctx, gfr)
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	core.WriteResponse(ctx, nil, gin.H{
		"total": glr.Total,
		"data": slice.Map(glr.Data, func(index int, src *proto.GoodsInfoResponse) map[string]any {
			return map[string]any{
				"id": src.Id,
				"name":        src.Name,
				"goods_brief": src.GoodsBrief,
				"desc":        src.GoodsDesc,
				"ship_free":   src.ShipFree,
				"images":      src.Images,
				"desc_images": src.DescImages,
				"front_image": src.GoodsFrontImage,
				"shop_price":  src.ShopPrice,
				"category": map[string]any{
						"id":   src.Category.Id,
						"name": src.Category.Name,
				},
				"brand": map[string]any{
						"id":   src.Brand.Id,
						"name": src.Brand.Name,
						"logo": src.Brand.Logo,
				},
				"is_hot":  src.IsHot,
				"is_new":  src.IsNew,
				"on_sale": src.OnSale,
			}
		}),
	})
}