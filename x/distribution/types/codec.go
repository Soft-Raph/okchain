package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
)

// RegisterCodec registers concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&MsgWithdrawValidatorCommission{}, "okchain/distribution/MsgWithdrawReward", nil)
	cdc.RegisterConcrete(&MsgSetWithdrawAddress{}, "okchain/distribution/MsgModifyWithdrawAddress", nil)
	cdc.RegisterConcrete(CommunityPoolSpendProposal{}, "okchain/distribution/CommunityPoolSpendProposal", nil)
}

// ModuleCdc generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	cryptocodec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
