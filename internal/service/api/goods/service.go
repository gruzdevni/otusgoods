package goods

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/samber/lo"

	"otusgoods/internal/models"
	query "otusgoods/internal/repo"
)

var ErrNoAvailableGoods = errors.New("Not enough goods for reserve")

type repo interface {
	CheckOrderReserve(ctx context.Context, orderID string) ([]query.GoodsReservation, error)
	DecreaseAvailableGoods(ctx context.Context, arg query.DecreaseAvailableGoodsParams) error
	GetAvailableGoods(ctx context.Context, id []int32) ([]query.AvailableQuantity, error)
	IncreaseAvailableGoods(ctx context.Context, arg query.IncreaseAvailableGoodsParams) error
	ReserveGoodsForOrder(ctx context.Context, arg query.ReserveGoodsForOrderParams) error
	UnreserveGoodsForOrder(ctx context.Context, orderID string) error
	WithTx(tx *sql.Tx) *query.Queries
}

type service struct {
	dbRW *sql.DB
	repo repo
}

type Service interface {
	CheckOrderReserve(ctx context.Context, orderID string) (*models.ReserveStatus, error)
	ReserveGoodsForOrder(ctx context.Context, arg models.NewReserveParams) error
	UnreserveGoodsForOrder(ctx context.Context, orderID string) error
}

func NewService(dbRW *sql.DB, repo repo) Service {
	return &service{
		dbRW: dbRW,
		repo: repo,
	}
}

func (s *service) CheckOrderReserve(ctx context.Context, orderID string) (*models.ReserveStatus, error) {
	res, err := s.repo.CheckOrderReserve(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("checking order reserve: %w", err)
	}

	if len(res) == 0 {
		return &models.ReserveStatus{
			OrderID:       orderID,
			ReserveStatus: "cancelled",
		}, nil
	}

	return &models.ReserveStatus{
		OrderID:       orderID,
		ReserveStatus: "confirmed",
	}, nil
}

func (s *service) ReserveGoodsForOrder(ctx context.Context, arg models.NewReserveParams) error {
	tx, err := s.dbRW.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}

	isCommitted := false
	defer func() {
		if !isCommitted {
			if err := tx.Rollback(); err != nil {
				zerolog.Ctx(ctx).Err(err).Msg("failed to rollback the transaction")
			}
		}
	}()

	goodsIds := make([]int32, 0, len(arg.Goods))

	for _, item := range arg.Goods {
		id, err := strconv.Atoi(item.Nomenclature)
		if err != nil {
			return fmt.Errorf("parsing goods id: %w", err)
		}

		goodsIds = append(goodsIds, int32(id))
	}

	res2, err := s.repo.WithTx(tx).GetAvailableGoods(ctx, goodsIds)
	if err != nil {
		return fmt.Errorf("getting available goods: %w", err)
	}

	if len(res2) == 0 {
		return ErrNoAvailableGoods
	}

	availableMap := lo.Associate(res2, func(item query.AvailableQuantity) (int32, int32) {
		return item.ID, item.AvailableQuantity
	})

	for _, item := range arg.Goods {
		id, err := strconv.Atoi(item.Nomenclature)
		if err != nil {
			return fmt.Errorf("parsing goods id: %w", err)
		}

		if availableMap[int32(id)]-int32(item.Quantity) < 0 {
			return ErrNoAvailableGoods
		}

		if err = s.repo.WithTx(tx).DecreaseAvailableGoods(ctx, query.DecreaseAvailableGoodsParams{
			RequestQuantity: int32(item.Quantity),
			ID:              int32(id),
		}); err != nil {
			return fmt.Errorf("decreasing available quantity: %w", err)
		}

		if err = s.repo.WithTx(tx).ReserveGoodsForOrder(ctx, query.ReserveGoodsForOrderParams{
			OrderID:          arg.OrderID,
			NomenclatureID:   int32(id),
			QuantityReserved: int32(item.Quantity),
		}); err != nil {
			return fmt.Errorf("reserving goods for order: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	isCommitted = true

	return nil
}

func (s *service) UnreserveGoodsForOrder(ctx context.Context, orderID string) error {
	tx, err := s.dbRW.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}

	isCommitted := false
	defer func() {
		if !isCommitted {
			if err := tx.Rollback(); err != nil {
				zerolog.Ctx(ctx).Err(err).Msg("failed to rollback the transaction")
			}
		}
	}()

	res, err := s.repo.CheckOrderReserve(ctx, orderID)
	if err != nil {
		return fmt.Errorf("checking order reserve: %w", err)
	}

	if len(res) == 0 {
		return nil
	}

	for _, item := range res {
		if err = s.repo.WithTx(tx).IncreaseAvailableGoods(ctx, query.IncreaseAvailableGoodsParams{
			RequestQuantity: item.QuantityReserved,
			ID:              item.NomenclatureID,
		}); err != nil {
			return fmt.Errorf("increasing available goods: %w", err)
		}
	}

	if err = s.repo.WithTx(tx).UnreserveGoodsForOrder(ctx, orderID); err != nil {
		return fmt.Errorf("unreserving goods for order: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	isCommitted = true

	return nil
}
