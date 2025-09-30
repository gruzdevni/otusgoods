package restapi

import (
	"otusgoods/internal/models"
	"otusgoods/internal/restapi/operations/other"
	"otusgoods/internal/restapi/operations/warehouse"
	"otusgoods/internal/service/api/goods"

	"github.com/go-openapi/runtime/middleware"
)

type Handler struct {
	goodsSrv goods.Service
}

func NewHandler(goodsSrv goods.Service) *Handler {
	return &Handler{
		goodsSrv: goodsSrv,
	}
}

func (h *Handler) GetHealth(_ other.GetHealthParams) middleware.Responder {
	return other.NewGetHealthOK().WithPayload(&models.DefaultStatusResponse{Code: "01", Message: "OK"})
}

func (h *Handler) CheckOrderReserve(params warehouse.GetCheckReserveStatusOrderIDParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()

	res, err := h.goodsSrv.CheckOrderReserve(ctx, params.OrderID)
	if err != nil {
		return warehouse.NewGetCheckReserveStatusOrderIDInternalServerError()
	}

	return warehouse.NewGetCheckReserveStatusOrderIDOK().WithPayload(res)
}

func (h *Handler) ReserveGoodsForOrder(params warehouse.PostReserveOrderGoodsParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()

	err := h.goodsSrv.ReserveGoodsForOrder(ctx, *params.Request)
	if err != nil {
		return warehouse.NewPostReserveOrderGoodsInternalServerError().WithPayload(&models.DefaultStatusResponse{Code: "03", Message: err.Error()})
	}

	return warehouse.NewPostReserveOrderGoodsOK()
}

func (h *Handler) UnreserveGoodsForOrder(params warehouse.DeleteUnreserveGoodsOrderIDParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()

	err := h.goodsSrv.UnreserveGoodsForOrder(ctx, params.OrderID)
	if err != nil {
		return warehouse.NewDeleteUnreserveGoodsOrderIDInternalServerError().WithPayload(&models.DefaultStatusResponse{Code: "03", Message: err.Error()})
	}

	return warehouse.NewDeleteUnreserveGoodsOrderIDOK()
}
