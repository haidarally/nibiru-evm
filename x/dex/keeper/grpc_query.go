package keeper

import (
	"context"
	"fmt"

	"github.com/MatrixDao/matrix/x/dex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/gogo/protobuf/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

/*
Handler for the QueryParamsRequest query.

args
  ctx: the cosmos-sdk context
  req: a QueryParamsRequest proto object

ret
  QueryParamsResponse: the QueryParamsResponse proto object response, containing the params
  error: an error if any occurred
*/
func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

/*
Handler for the QueryGetPoolRequest query.

args
  ctx: the cosmos-sdk context
  req: a QueryGetPoolRequest proto object

ret
  QueryGetPoolResponse: the QueryGetPoolResponse proto object response, containing the pool
  error: an error if any occurred
*/
func (k Keeper) GetPool(goCtx context.Context, req *types.QueryGetPoolRequest) (*types.QueryGetPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	poolKey := types.GetKeyPrefixPools(req.PoolId)
	bz := store.Get(poolKey)

	var pool types.Pool
	err := pool.Unmarshal(bz)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetPoolResponse{
		Pool: &pool,
	}, nil
}

/*
Handler for the QueryGetPoolNumberRequest query.

args
  ctx: the cosmos-sdk context
  req: a QueryGetPoolNumberRequest proto object

ret
  QueryGetPoolNumberResponse: the QueryGetPoolNumberResponse proto object response, containing the next pool id number
  error: an error if any occurred
*/
func (k Keeper) GetPoolNumber(goCtx context.Context, req *types.QueryGetPoolNumberRequest) (*types.QueryGetPoolNumberResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var poolNumber uint64

	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyNextGlobalPoolNumber)
	if bz == nil {
		k.Logger(ctx).Error("Could not get pool number. Not initialized.")
		panic(fmt.Errorf("pool number has not been initialized -- Should have been done in InitGenesis"))
	} else {
		val := gogotypes.UInt64Value{}

		err := k.cdc.Unmarshal(bz, &val)
		if err != nil {
			panic(err)
		}

		poolNumber = val.GetValue()
	}

	return &types.QueryGetPoolNumberResponse{
		PoolId: poolNumber,
	}, nil
}