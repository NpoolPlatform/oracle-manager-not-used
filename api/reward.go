// +build !codeanalysis

package api

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	constant "github.com/NpoolPlatform/oracle-manager/pkg/const"
	crud "github.com/NpoolPlatform/oracle-manager/pkg/crud/reward"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/oraclemgr"

	"github.com/google/uuid"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateReward(ctx context.Context, in *npool.CreateRewardRequest) (*npool.CreateRewardResponse, error) {
	if _, err := uuid.Parse(in.GetInfo().GetCoinTypeID()); err != nil {
		logger.Sugar().Errorf("invalid request coin type id: %v", err)
		return &npool.CreateRewardResponse{}, status.Error(codes.Internal, err.Error())
	}

	if in.GetInfo().GetDailyReward() <= 0 {
		logger.Sugar().Errorf("invalid daily reward")
		return &npool.CreateRewardResponse{}, status.Error(codes.Internal, "invalid reward parameter")
	}

	schema, err := crud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail create schema entity: %v", err)
		return &npool.CreateRewardResponse{}, status.Error(codes.Internal, err.Error())
	}

	info, err := schema.Create(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorf("fail create reward: %v", err)
		return &npool.CreateRewardResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateRewardResponse{
		Info: info,
	}, nil
}

func (s *Server) CreateRewards(ctx context.Context, in *npool.CreateRewardsRequest) (*npool.CreateRewardsResponse, error) {
	for _, info := range in.GetInfos() {
		if _, err := uuid.Parse(info.GetCoinTypeID()); err != nil {
			logger.Sugar().Errorf("invalid request coin type id: %v", err)
			return &npool.CreateRewardsResponse{}, status.Error(codes.Internal, err.Error())
		}
	}

	schema, err := crud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail create schema entity: %v", err)
		return &npool.CreateRewardsResponse{}, status.Error(codes.Internal, err.Error())
	}

	infos, err := schema.CreateBulk(ctx, in.GetInfos())
	if err != nil {
		logger.Sugar().Errorf("fail create rewards: %v", err)
		return &npool.CreateRewardsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateRewardsResponse{
		Infos: infos,
	}, nil
}

func (s *Server) UpdateReward(ctx context.Context, in *npool.UpdateRewardRequest) (*npool.UpdateRewardResponse, error) {
	if _, err := uuid.Parse(in.GetInfo().GetCoinTypeID()); err != nil {
		logger.Sugar().Errorf("invalid request coin type id: %v", err)
		return &npool.UpdateRewardResponse{}, status.Error(codes.Internal, err.Error())
	}

	if _, err := uuid.Parse(in.GetInfo().GetID()); err != nil {
		logger.Sugar().Errorf("invalid reward id: %v", err)
		return &npool.UpdateRewardResponse{}, status.Error(codes.Internal, err.Error())
	}

	schema, err := crud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail create schema entity: %v", err)
		return &npool.UpdateRewardResponse{}, status.Error(codes.Internal, err.Error())
	}

	info, err := schema.Update(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorf("fail update reward: %v", err)
		return &npool.UpdateRewardResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateRewardResponse{
		Info: info,
	}, nil
}

func rewardCondsToConds(conds cruder.FilterConds) (cruder.Conds, error) {
	newConds := cruder.NewConds()

	for k, v := range conds {
		switch v.Op {
		case cruder.EQ:
		case cruder.GT:
		case cruder.LT:
		case cruder.LIKE:
		default:
			return nil, fmt.Errorf("invalid filter condition op")
		}

		switch k {
		case constant.FieldID:
			fallthrough //nolint
		case constant.FieldCoinTypeID:
			newConds = newConds.WithCond(k, v.Op, v.Val.GetStringValue())
		case constant.RewardFieldDailyReward:
			newConds = newConds.WithCond(k, v.Op, v.Val.GetNumberValue())
		default:
			return nil, fmt.Errorf("invalid reward field")
		}
	}

	return newConds, nil
}

func (s *Server) GetReward(ctx context.Context, in *npool.GetRewardRequest) (*npool.GetRewardResponse, error) {
	id, err := uuid.Parse(in.GetID())
	if err != nil {
		return &npool.GetRewardResponse{}, fmt.Errorf("invalid reward id: %v", err)
	}

	schema, err := crud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail create schema entity: %v", err)
		return &npool.GetRewardResponse{}, status.Error(codes.Internal, err.Error())
	}

	info, err := schema.Row(ctx, id)
	if err != nil {
		logger.Sugar().Errorf("fail get reward: %v", err)
		return &npool.GetRewardResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetRewardResponse{
		Info: info,
	}, nil
}

func (s *Server) GetRewardOnly(ctx context.Context, in *npool.GetRewardOnlyRequest) (*npool.GetRewardOnlyResponse, error) {
	conds, err := rewardCondsToConds(in.GetConds())
	if err != nil {
		logger.Sugar().Errorf("invalid reward fields: %v", err)
		return &npool.GetRewardOnlyResponse{}, status.Error(codes.Internal, err.Error())
	}

	schema, err := crud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail create schema entity: %v", err)
		return &npool.GetRewardOnlyResponse{}, status.Error(codes.Internal, err.Error())
	}

	info, err := schema.RowOnly(ctx, conds)
	if err != nil {
		logger.Sugar().Errorf("fail get rewards: %v", err)
		return &npool.GetRewardOnlyResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetRewardOnlyResponse{
		Info: info,
	}, nil
}

func (s *Server) GetRewards(ctx context.Context, in *npool.GetRewardsRequest) (*npool.GetRewardsResponse, error) {
	conds, err := rewardCondsToConds(in.GetConds())
	if err != nil {
		logger.Sugar().Errorf("invalid reward fields: %v", err)
		return &npool.GetRewardsResponse{}, status.Error(codes.Internal, err.Error())
	}

	schema, err := crud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail create schema entity: %v", err)
		return &npool.GetRewardsResponse{}, status.Error(codes.Internal, err.Error())
	}

	infos, total, err := schema.Rows(ctx, conds, int(in.GetOffset()), int(in.GetLimit()))
	if err != nil {
		logger.Sugar().Errorf("fail get rewards: %v", err)
		return &npool.GetRewardsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetRewardsResponse{
		Infos: infos,
		Total: int32(total),
	}, nil
}

func (s *Server) ExistReward(ctx context.Context, in *npool.ExistRewardRequest) (*npool.ExistRewardResponse, error) {
	id, err := uuid.Parse(in.GetID())
	if err != nil {
		return &npool.ExistRewardResponse{}, fmt.Errorf("invalid reward id: %v", err)
	}

	schema, err := crud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail create schema entity: %v", err)
		return &npool.ExistRewardResponse{}, status.Error(codes.Internal, err.Error())
	}

	exist, err := schema.Exist(ctx, id)
	if err != nil {
		logger.Sugar().Errorf("fail check reward: %v", err)
		return &npool.ExistRewardResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.ExistRewardResponse{
		Result: exist,
	}, nil
}

func (s *Server) ExistRewardConds(ctx context.Context, in *npool.ExistRewardCondsRequest) (*npool.ExistRewardCondsResponse, error) {
	conds, err := rewardCondsToConds(in.GetConds())
	if err != nil {
		logger.Sugar().Errorf("invalid reward fields: %v", err)
		return &npool.ExistRewardCondsResponse{}, status.Error(codes.Internal, err.Error())
	}

	if len(conds) == 0 {
		logger.Sugar().Errorf("empty reward fields: %v", err)
		return &npool.ExistRewardCondsResponse{}, status.Error(codes.Internal, "empty reward fields")
	}

	schema, err := crud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail create schema entity: %v", err)
		return &npool.ExistRewardCondsResponse{}, status.Error(codes.Internal, err.Error())
	}

	exist, err := schema.ExistConds(ctx, conds)
	if err != nil {
		logger.Sugar().Errorf("fail check reward: %v", err)
		return &npool.ExistRewardCondsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.ExistRewardCondsResponse{
		Result: exist,
	}, nil
}

func (s *Server) CountRewards(ctx context.Context, in *npool.CountRewardsRequest) (*npool.CountRewardsResponse, error) {
	conds, err := rewardCondsToConds(in.GetConds())
	if err != nil {
		logger.Sugar().Errorf("invalid reward fields: %v", err)
		return &npool.CountRewardsResponse{}, status.Error(codes.Internal, err.Error())
	}

	if len(conds) == 0 {
		logger.Sugar().Errorf("empty reward fields: %v", err)
		return &npool.CountRewardsResponse{}, status.Error(codes.Internal, "empty reward fields")
	}

	schema, err := crud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail create schema entity: %v", err)
		return &npool.CountRewardsResponse{}, status.Error(codes.Internal, err.Error())
	}

	total, err := schema.Count(ctx, conds)
	if err != nil {
		logger.Sugar().Errorf("fail count rewards: %v", err)
		return &npool.CountRewardsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CountRewardsResponse{
		Result: total,
	}, nil
}

func (s *Server) DeleteReward(ctx context.Context, in *npool.DeleteRewardRequest) (*npool.DeleteRewardResponse, error) {
	id, err := uuid.Parse(in.GetID())
	if err != nil {
		return &npool.DeleteRewardResponse{}, fmt.Errorf("invalid reward id: %v", err)
	}

	schema, err := crud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail create schema entity: %v", err)
		return &npool.DeleteRewardResponse{}, status.Error(codes.Internal, err.Error())
	}

	info, err := schema.Delete(ctx, id)
	if err != nil {
		logger.Sugar().Errorf("fail delete reward: %v", err)
		return &npool.DeleteRewardResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteRewardResponse{
		Info: info,
	}, nil
}
