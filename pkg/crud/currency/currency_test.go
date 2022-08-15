package currency

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
	currency := npool.Currency{
		AppID:          uuid.New().String(),
		CoinTypeID:     uuid.New().String(),
		PriceVSUSDT:    1000,
		AppPriceVSUSDT: 1200,
	}

	schema, err := New(context.Background(), nil)
	assert.Nil(t, err)

	info, err := schema.Create(context.Background(), &currency)
	if assert.Nil(t, err) {
		if assert.NotEqual(t, info.ID, uuid.UUID{}.String()) {
			currency.ID = info.ID
		}
		assert.Equal(t, info, &currency)
	}

	currency.ID = info.ID

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	info, err = schema.Update(context.Background(), &currency)
	if assert.Nil(t, err) {
		assert.Equal(t, info, &currency)
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	info, err = schema.Row(context.Background(), uuid.MustParse(info.ID))
	if assert.Nil(t, err) {
		assert.Equal(t, info, &currency)
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	infos, total, err := schema.Rows(context.Background(),
		cruder.NewConds().WithCond(constant.FieldID, cruder.EQ, info.ID),
		0, 0)
	if assert.Nil(t, err) {
		assert.Equal(t, total, 1)
		assert.Equal(t, infos[0], &currency)
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	info, err = schema.RowOnly(context.Background(),
		cruder.NewConds().WithCond(constant.FieldID, cruder.EQ, info.ID))
	if assert.Nil(t, err) {
		assert.Equal(t, info, &currency)
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	info, err = schema.Delete(context.Background(), uuid.MustParse(info.ID))
	if assert.Nil(t, err) {
		assert.Equal(t, info, &currency)
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	info, err = schema.Row(context.Background(), uuid.MustParse(info.ID))
	assert.NotNil(t, err)
	assert.Nil(t, info)
}
