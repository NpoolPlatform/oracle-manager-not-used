package client

import (
	"context"
	"fmt"
	"time"

	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/oraclemgr"

	constant "github.com/NpoolPlatform/oracle-manager/pkg/message/const"
)

func do(ctx context.Context, fn func(_ctx context.Context, cli npool.OracleManagerClient) (cruder.Any, error)) (cruder.Any, error) {
	_ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	conn, err := grpc2.GetGRPCConn(constant.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get oracle connection: %v", err)
	}
	defer conn.Close()

	cli := npool.NewOracleManagerClient(conn)

	return fn(_ctx, cli)
}

func GetCurrencyOnly(ctx context.Context, conds cruder.FilterConds) (*npool.Currency, error) {
	info, err := do(ctx, func(_ctx context.Context, cli npool.OracleManagerClient) (cruder.Any, error) {
		resp, err := cli.GetCurrencyOnly(ctx, &npool.GetCurrencyOnlyRequest{
			Conds: conds,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get oracle: %v", err)
		}
		return resp.Info, nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail get oracle: %v", err)
	}
	return info.(*npool.Currency), nil
}
