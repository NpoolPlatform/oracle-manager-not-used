package currency

import (
	"context"
	"fmt"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/price"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/oracle-manager/pkg/db/ent"
	"github.com/NpoolPlatform/oracle-manager/pkg/db/ent/currency"

	constant "github.com/NpoolPlatform/oracle-manager/pkg/const"
	"github.com/NpoolPlatform/oracle-manager/pkg/db"

	npool "github.com/NpoolPlatform/message/npool/oraclemgr"

	"github.com/google/uuid"
)

type Currency struct {
	*db.Entity
}

func New(ctx context.Context, tx *ent.Tx) (*Currency, error) {
	e, err := db.NewEntity(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("fail create entity: %v", err)
	}

	return &Currency{
		Entity: e,
	}, nil
}

func (s *Currency) rowToObject(row *ent.Currency) *npool.Currency {
	return &npool.Currency{
		ID:             row.ID.String(),
		AppID:          row.AppID.String(),
		CoinTypeID:     row.CoinTypeID.String(),
		PriceVSUSDT:    price.DBPriceToVisualPrice(row.PriceVsUsdt),
		AppPriceVSUSDT: price.DBPriceToVisualPrice(row.AppPriceVsUsdt),
		OverPercent:    row.OverPercent,
		CurrencyMethod: row.CurrencyMethod,
	}
}

func (s *Currency) Create(ctx context.Context, in *npool.Currency) (*npool.Currency, error) {
	var info *ent.Currency
	var err error

	err = db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		info, err = s.Tx.Currency.Create().
			SetAppID(uuid.MustParse(in.GetAppID())).
			SetCoinTypeID(uuid.MustParse(in.GetCoinTypeID())).
			SetPriceVsUsdt(price.VisualPriceToDBPrice(in.GetPriceVSUSDT())).
			SetAppPriceVsUsdt(price.VisualPriceToDBPrice(in.GetAppPriceVSUSDT())).
			SetOverPercent(in.GetOverPercent()).
			SetCurrencyMethod(in.GetCurrencyMethod()).
			Save(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail create currency: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *Currency) Update(ctx context.Context, in *npool.Currency) (*npool.Currency, error) {
	var info *ent.Currency
	var err error

	err = db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		info, err = s.Tx.Currency.UpdateOneID(uuid.MustParse(in.GetID())).
			SetPriceVsUsdt(price.VisualPriceToDBPrice(in.GetPriceVSUSDT())).
			SetAppPriceVsUsdt(price.VisualPriceToDBPrice(in.GetAppPriceVSUSDT())).
			SetOverPercent(in.GetOverPercent()).
			SetCurrencyMethod(in.GetCurrencyMethod()).
			Save(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail update currency: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *Currency) Row(ctx context.Context, id uuid.UUID) (*npool.Currency, error) {
	var info *ent.Currency
	var err error

	err = db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		info, err = s.Tx.Currency.Query().Where(currency.ID(id)).Only(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail get currency: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *Currency) queryFromConds(conds cruder.Conds) (*ent.CurrencyQuery, error) { //nolint
	stm := s.Tx.Currency.Query()
	for k, v := range conds {
		switch k {
		case constant.FieldID:
			id, err := cruder.AnyTypeUUID(v.Val)
			if err != nil {
				return nil, fmt.Errorf("invalid id: %v", err)
			}
			stm = stm.Where(currency.ID(id))
		case constant.FieldAppID:
			id, err := cruder.AnyTypeUUID(v.Val)
			if err != nil {
				return nil, fmt.Errorf("invalid app id: %v", err)
			}
			stm = stm.Where(currency.AppID(id))
		case constant.FieldCoinTypeID:
			id, err := cruder.AnyTypeUUID(v.Val)
			if err != nil {
				return nil, fmt.Errorf("invalid coin type id: %v", err)
			}
			stm = stm.Where(currency.CoinTypeID(id))
		case constant.CurrencyFieldCurrencyMethod:
			method, err := cruder.AnyTypeString(v.Val)
			if err != nil {
				return nil, fmt.Errorf("invalid currency method: %v", err)
			}
			stm = stm.Where(currency.CurrencyMethod(method))
		default:
			return nil, fmt.Errorf("invalid currency field")
		}
	}

	return stm, nil
}

func (s *Currency) Rows(ctx context.Context, conds cruder.Conds, offset, limit int) ([]*npool.Currency, int, error) {
	rows := []*ent.Currency{}
	var total int

	err := db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		stm, err := s.queryFromConds(conds)
		if err != nil {
			return fmt.Errorf("fail construct stm: %v", err)
		}

		total, err = stm.Count(_ctx)
		if err != nil {
			return fmt.Errorf("fail count currency: %v", err)
		}

		rows, err = stm.Order(ent.Desc(currency.FieldUpdatedAt)).Offset(offset).Limit(limit).All(_ctx)
		if err != nil {
			return fmt.Errorf("fail query currency: %v", err)
		}

		return nil
	})
	if err != nil {
		return nil, 0, fmt.Errorf("fail get currency: %v", err)
	}

	infos := []*npool.Currency{}
	for _, row := range rows {
		infos = append(infos, s.rowToObject(row))
	}

	return infos, total, nil
}

func (s *Currency) RowOnly(ctx context.Context, conds cruder.Conds) (*npool.Currency, error) {
	var info *ent.Currency

	err := db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		stm, err := s.queryFromConds(conds)
		if err != nil {
			return fmt.Errorf("fail construct stm: %v", err)
		}

		info, err = stm.Only(_ctx)
		if err != nil {
			return fmt.Errorf("fail query currency: %v", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail get currency: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *Currency) Delete(ctx context.Context, id uuid.UUID) (*npool.Currency, error) {
	var info *ent.Currency
	var err error

	err = db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		info, err = s.Tx.Currency.UpdateOneID(id).
			SetDeletedAt(uint32(time.Now().Unix())).
			Save(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail delete currency: %v", err)
	}

	return s.rowToObject(info), nil
}
