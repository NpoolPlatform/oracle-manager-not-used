package reward

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/oraclemgr"
	"github.com/NpoolPlatform/oracle-manager/pkg/test-init" //nolint

	constant "github.com/NpoolPlatform/oracle-manager/pkg/const"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func init() {
	if runByGithubAction, err := strconv.ParseBool(os.Getenv("RUN_BY_GITHUB_ACTION")); err == nil && runByGithubAction {
		return
	}
	if err := testinit.Init(); err != nil {
		fmt.Printf("cannot init test stub: %v\n", err)
	}
}

func TestCRUD(t *testing.T) {
	reward := npool.Reward{
		CoinTypeID:  uuid.New().String(),
		DailyReward: 1000,
	}

	schema, err := New(context.Background(), nil)
	assert.Nil(t, err)

	info, err := schema.Create(context.Background(), &reward)
	if assert.Nil(t, err) {
		if assert.NotEqual(t, info.ID, uuid.UUID{}.String()) {
			reward.ID = info.ID
		}
		assert.Equal(t, info, &reward)
	}

	reward.ID = info.ID

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	info, err = schema.Update(context.Background(), &reward)
	if assert.Nil(t, err) {
		assert.Equal(t, info, &reward)
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	info, err = schema.Row(context.Background(), uuid.MustParse(info.ID))
	if assert.Nil(t, err) {
		assert.Equal(t, info, &reward)
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	infos, total, err := schema.Rows(context.Background(),
		cruder.NewConds().WithCond(constant.FieldID, cruder.EQ, info.ID),
		0, 0)
	if assert.Nil(t, err) {
		assert.Equal(t, total, 1)
		assert.Equal(t, infos[0], &reward)
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	info, err = schema.RowOnly(context.Background(),
		cruder.NewConds().WithCond(constant.FieldID, cruder.EQ, info.ID))
	if assert.Nil(t, err) {
		assert.Equal(t, info, &reward)
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	exist, err := schema.ExistConds(context.Background(),
		cruder.NewConds().WithCond(constant.FieldID, cruder.EQ, info.ID),
	)
	if assert.Nil(t, err) {
		assert.Equal(t, exist, true)
	}

	reward1 := &npool.Reward{
		CoinTypeID:  uuid.New().String(),
		DailyReward: 1000,
	}
	reward2 := &npool.Reward{
		CoinTypeID:  uuid.New().String(),
		DailyReward: 1000,
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	infos, err = schema.CreateBulk(context.Background(), []*npool.Reward{reward1, reward2})
	if assert.Nil(t, err) {
		assert.Equal(t, len(infos), 2)
		assert.NotEqual(t, infos[0].ID, uuid.UUID{}.String())
		assert.NotEqual(t, infos[1].ID, uuid.UUID{}.String())
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	count, err := schema.Count(context.Background(),
		cruder.NewConds().WithCond(constant.FieldID, cruder.EQ, info.ID),
	)
	if assert.Nil(t, err) {
		assert.Equal(t, count, uint32(1))
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	info, err = schema.Delete(context.Background(), uuid.MustParse(info.ID))
	if assert.Nil(t, err) {
		assert.Equal(t, info, &reward)
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	count, err = schema.Count(context.Background(),
		cruder.NewConds().WithCond(constant.FieldID, cruder.EQ, info.ID),
	)
	if assert.Nil(t, err) {
		assert.Equal(t, count, uint32(0))
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	_, err = schema.Row(context.Background(), uuid.MustParse(info.ID))
	assert.NotNil(t, err)
}
