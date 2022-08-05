package reward

import (
	"context"
	"fmt"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/price"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/oracle-manager/pkg/db/ent"
	"github.com/NpoolPlatform/oracle-manager/pkg/db/ent/reward"

	constant "github.com/NpoolPlatform/oracle-manager/pkg/const"
	"github.com/NpoolPlatform/oracle-manager/pkg/db"

	npool "github.com/NpoolPlatform/message/npool/oraclemgr"

	"github.com/google/uuid"
)

type Reward struct {
	*db.Entity
}

func New(ctx context.Context, tx *ent.Tx) (*Reward, error) {
	e, err := db.NewEntity(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("fail create entity: %v", err)
	}

	return &Reward{
		Entity: e,
	}, nil
}

func (s *Reward) rowToObject(row *ent.Reward) *npool.Reward {
	if row == nil {
		return nil
	}

	return &npool.Reward{
		ID:          row.ID.String(),
		CoinTypeID:  row.CoinTypeID.String(),
		DailyReward: price.DBPriceToVisualPrice(row.DailyReward),
	}
}

func (s *Reward) Create(ctx context.Context, in *npool.Reward) (*npool.Reward, error) {
	var info *ent.Reward
	var err error

	err = db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		info, err = s.Tx.Reward.Create().
			SetCoinTypeID(uuid.MustParse(in.GetCoinTypeID())).
			SetDailyReward(price.VisualPriceToDBPrice(in.GetDailyReward())).
			Save(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail create reward: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *Reward) CreateBulk(ctx context.Context, in []*npool.Reward) ([]*npool.Reward, error) {
	rows := []*ent.Reward{}
	var err error

	err = db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		bulk := make([]*ent.RewardCreate, len(in))
		for i, info := range in {
			bulk[i] = s.Tx.Reward.Create().
				SetCoinTypeID(uuid.MustParse(info.GetCoinTypeID())).
				SetDailyReward(price.VisualPriceToDBPrice(info.GetDailyReward()))
		}
		rows, err = s.Tx.Reward.CreateBulk(bulk...).Save(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail create rewards: %v", err)
	}

	infos := []*npool.Reward{}
	for _, row := range rows {
		infos = append(infos, s.rowToObject(row))
	}

	return infos, nil
}

func (s *Reward) Update(ctx context.Context, in *npool.Reward) (*npool.Reward, error) {
	var info *ent.Reward
	var err error

	err = db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		info, err = s.Tx.Reward.UpdateOneID(uuid.MustParse(in.GetID())).
			SetDailyReward(price.VisualPriceToDBPrice(in.GetDailyReward())).
			Save(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail update reward: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *Reward) Row(ctx context.Context, id uuid.UUID) (*npool.Reward, error) {
	var info *ent.Reward
	var err error

	err = db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		info, err = s.Tx.Reward.Query().Where(reward.ID(id)).Only(_ctx)
		if ent.IsNotFound(err) {
			return nil
		}
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail get reward: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *Reward) queryFromConds(conds cruder.Conds) (*ent.RewardQuery, error) { //nolint
	stm := s.Tx.Reward.Query()
	for k, v := range conds {
		switch k {
		case constant.FieldID:
			id, err := cruder.AnyTypeUUID(v.Val)
			if err != nil {
				return nil, fmt.Errorf("invalid id: %v", err)
			}
			stm = stm.Where(reward.ID(id))
		case constant.FieldCoinTypeID:
			id, err := cruder.AnyTypeUUID(v.Val)
			if err != nil {
				return nil, fmt.Errorf("invalid coin type id: %v", err)
			}
			stm = stm.Where(reward.CoinTypeID(id))
		case constant.RewardFieldDailyReward:
			fvalue, err := cruder.AnyTypeFloat64(v.Val)
			if err != nil {
				return nil, fmt.Errorf("invalid daily reward value: %v", err)
			}

			value := price.VisualPriceToDBPrice(fvalue)

			switch v.Op {
			case cruder.EQ:
				stm = stm.Where(reward.DailyRewardEQ(value))
			case cruder.GT:
				stm = stm.Where(reward.DailyRewardGT(value))
			case cruder.LT:
				stm = stm.Where(reward.DailyRewardLT(value))
			}
		default:
			return nil, fmt.Errorf("invalid reward field")
		}
	}

	return stm, nil
}

func (s *Reward) Rows(ctx context.Context, conds cruder.Conds, offset, limit int) ([]*npool.Reward, int, error) {
	rows := []*ent.Reward{}
	var total int

	err := db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		stm, err := s.queryFromConds(conds)
		if err != nil {
			return fmt.Errorf("fail construct stm: %v", err)
		}

		total, err = stm.Count(_ctx)
		if err != nil {
			return fmt.Errorf("fail count reward: %v", err)
		}

		rows, err = stm.Order(ent.Desc(reward.FieldUpdatedAt)).Offset(offset).Limit(limit).All(_ctx)
		if err != nil {
			return fmt.Errorf("fail query reward: %v", err)
		}

		return nil
	})
	if err != nil {
		return nil, 0, fmt.Errorf("fail get reward: %v", err)
	}

	infos := []*npool.Reward{}
	for _, row := range rows {
		infos = append(infos, s.rowToObject(row))
	}

	return infos, total, nil
}

func (s *Reward) RowOnly(ctx context.Context, conds cruder.Conds) (*npool.Reward, error) {
	var info *ent.Reward

	err := db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		stm, err := s.queryFromConds(conds)
		if err != nil {
			return fmt.Errorf("fail construct stm: %v", err)
		}

		info, err = stm.Only(_ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil
			}
			return fmt.Errorf("fail query reward: %v", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail get reward: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *Reward) Count(ctx context.Context, conds cruder.Conds) (uint32, error) {
	var err error
	var total int

	err = db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		stm, err := s.queryFromConds(conds)
		if err != nil {
			return fmt.Errorf("fail construct stm: %v", err)
		}

		total, err = stm.Count(_ctx)
		if err != nil {
			return fmt.Errorf("fail check rewards: %v", err)
		}

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("fail count rewards: %v", err)
	}

	return uint32(total), nil
}

func (s *Reward) Exist(ctx context.Context, id uuid.UUID) (bool, error) {
	var err error
	exist := false

	err = db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		exist, err = s.Tx.Reward.Query().Where(reward.ID(id)).Exist(_ctx)
		return err
	})
	if err != nil {
		return false, fmt.Errorf("fail check reward: %v", err)
	}

	return exist, nil
}

func (s *Reward) ExistConds(ctx context.Context, conds cruder.Conds) (bool, error) {
	var err error
	exist := false

	err = db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		stm, err := s.queryFromConds(conds)
		if err != nil {
			return fmt.Errorf("fail construct stm: %v", err)
		}

		exist, err = stm.Exist(_ctx)
		if err != nil {
			return fmt.Errorf("fail check rewards: %v", err)
		}

		return nil
	})
	if err != nil {
		return false, fmt.Errorf("fail check rewards: %v", err)
	}

	return exist, nil
}

func (s *Reward) Delete(ctx context.Context, id uuid.UUID) (*npool.Reward, error) {
	var info *ent.Reward
	var err error

	err = db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		info, err = s.Tx.Reward.UpdateOneID(id).
			SetDeletedAt(uint32(time.Now().Unix())).
			Save(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail delete reward: %v", err)
	}

	return s.rowToObject(info), nil
}
