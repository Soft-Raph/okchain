package cli

import (
	"bufio"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"strings"

	"github.com/cosmos/cosmos-sdk/version"
	"github.com/okex/okchain/x/gov"

	"github.com/pkg/errors"


	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/okex/okchain/x/common"
	dexUtils "github.com/okex/okchain/x/dex/client/utils"
	"github.com/okex/okchain/x/dex/types"
	"github.com/spf13/cobra"
)

// Dex tags
const (
	FlagBaseAsset  = "base-asset"
	FlagQuoteAsset = "quote-asset"
	FlagInitPrice  = "init-price"
	FlagProduct    = "product"
	FlagFrom       = "from"
	FlagTo         = "to"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "dex",
		Short: "Decentralized exchange management subcommands",
	}

	txCmd.AddCommand(flags.PostCommands(
		getCmdList(cdc),
		getCmdDeposit(cdc),
		getCmdWithdraw(cdc),
		getCmdTransferOwnership(cdc),
		getMultiSignsCmd(cdc),
	)...)

	return txCmd
}

func getCmdList(cdc *codec.Codec) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list a trading pair",
		Args:  cobra.ExactArgs(0),
		Long: strings.TrimSpace(`List a trading pair:

$ okchaincli tx dex list --base-asset mytoken --quote-asset okt --from mykey
`),
		RunE: func(cmd *cobra.Command, _ []string) error {
			cliCtx := client.NewContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			if err := authtypes.NewAccountRetriever(authclient.Codec).EnsureExists(cliCtx, cliCtx.FromAddress); err != nil {
				return err
			}

			flags := cmd.Flags()
			baseAsset, err := flags.GetString(FlagBaseAsset)
			if err != nil {
				return err
			}
			quoteAsset, err := flags.GetString(FlagQuoteAsset)
			if err != nil {
				return err
			}
			strInitPrice, err := flags.GetString(FlagInitPrice)
			if err != nil {
				return err
			}
			initPrice := sdk.MustNewDecFromStr(strInitPrice)
			owner := cliCtx.GetFromAddress()
			listMsg := types.NewMsgList(owner, baseAsset, quoteAsset, initPrice)
			return authclient.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{&listMsg})
		},
	}

	cmd.Flags().StringP(FlagBaseAsset, "", "btc", FlagBaseAsset+" should be issued before listed to opendex")
	cmd.Flags().StringP(FlagQuoteAsset, "", common.NativeToken, FlagQuoteAsset+" should be issued before listed to opendex")
	cmd.Flags().StringP(FlagInitPrice, "", "0.01", FlagInitPrice+" should be valid price")

	return cmd
}

// getCmdDeposit implements depositing tokens for a product.
func getCmdDeposit(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "deposit [product] [amount]",
		Args:  cobra.ExactArgs(2),
		Short: "deposit an amount of token on a product",
		Long: strings.TrimSpace(`Deposit an amount of token on a product:

$ okchaincli tx dex deposit mytoken_okt 1000okt --from mykey

The 'product' is a trading pair in full name of the tokens: ${base-asset-symbol}_${quote-asset-symbol}, for example 'mytoken_okt'.
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := client.NewContext().WithCodec(cdc)

			product := args[0]

			// Get depositor address
			from := cliCtx.GetFromAddress()

			// Get amount of coins
			amount, err := sdk.ParseDecCoin(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgDeposit(product, amount, from)
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{&msg})
		},
	}
}

// getCmdWithdraw implements withdrawing tokens from a product.
func getCmdWithdraw(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "withdraw [product] [amount]",
		Args:  cobra.ExactArgs(2),
		Short: "withdraw an amount of token from a product",
		Long: strings.TrimSpace(`Withdraw an amount of token from a product:

$ okchaincli tx dex withdraw mytoken_okt 1000okt --from mykey

The 'product' is a trading pair in full name of the tokens: ${base-asset-symbol}_${quote-asset-symbol}, for example 'mytoken_okt'.
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := client.NewContext().WithCodec(cdc)

			product := args[0]

			// Get depositor address
			from := cliCtx.GetFromAddress()

			// Get amount of coins
			amount, err := sdk.ParseDecCoin(args[1])
			if err != nil {
				return err
			}
			msg := types.NewMsgWithdraw(product, amount, from)
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{&msg})
		},
	}
}

// getCmdTransferOwnership is the CLI command for transfer ownership of product
func getCmdTransferOwnership(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-ownership",
		Short: "change the owner of the product",
		RunE: func(cmd *cobra.Command, _ []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := client.NewContext().WithCodec(cdc)
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			if err := authtypes.NewAccountRetriever(authclient.Codec).EnsureExists(cliCtx, cliCtx.FromAddress); err != nil {
				return err
			}
			flags := cmd.Flags()

			product, err := flags.GetString(FlagProduct)
			if err != nil || product == "" {
				return fmt.Errorf("invalid product:%s", product)
			}

			to, err := flags.GetString(FlagTo)
			if err != nil {
				return fmt.Errorf("invalid to:%s", to)
			}

			toAddr, err := sdk.AccAddressFromBech32(to)
			if err != nil {
				return fmt.Errorf("invalid to:%s", to)
			}

			from := cliCtx.GetFromAddress()
			msg := types.NewMsgTransferOwnership(from, toAddr, product,nil, nil)
			return authclient.PrintUnsignedStdTx(txBldr, cliCtx, []sdk.Msg{&msg})
		},
	}
	cmd.Flags().StringP(FlagProduct, "p", "", "product to be transferred")
	cmd.Flags().String(FlagTo, "", "the user to be transferred")
	return cmd
}

func getMultiSignsCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multisign",
		Short: "append signature to the unsigned tx file of transfer-ownership",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := client.NewContext().WithCodec(cdc)
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))

			stdTx, err := authclient.ReadTxFromFile(cliCtx, args[0])
			if err != nil {
				return err
			}

			if len(stdTx.GetMsgs()) == 0 {
				return errors.New("msg is empty")
			}

			msg, ok := stdTx.GetMsgs()[0].(*types.MsgTransferOwnership)
			if !ok {
				return errors.New("invalid msg type")
			}

			flags := cmd.Flags()
			_, err = flags.GetString(FlagFrom)
			if err != nil {
				return fmt.Errorf("invalid from:%s", err.Error())
			}

			signature, pub, err := txBldr.Keybase().Sign(cliCtx.GetFromName(), msg.GetSignBytes())
			if err != nil {
				return fmt.Errorf("sign failed:%s", err.Error())
			}

			msg.Pubkey = pub.Bytes()
			msg.ToSignature = signature
			return authclient.PrintUnsignedStdTx(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
	return cmd
}

// GetCmdSubmitDelistProposal implememts a command handler for submitting a dex delist proposal transaction
func GetCmdSubmitDelistProposal(cliCtx client.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "delist-proposal [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a dex delist proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a dex delist proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal delist-proposal <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
 "title": "delist xxx/%s",
 "description": "delist asset from dex",
 "base_asset": "xxx",
 "quote_asset": "%s",
 "deposit": [
   {
     "denom": "%s",
     "amount": "100"
   }
 ]
}
`, version.ClientName, sdk.DefaultBondDenom, sdk.DefaultBondDenom, sdk.DefaultBondDenom,
			)),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cliCtx.Codec))
			cliCtx := client.NewContext().WithCodec(cliCtx.Codec)

			proposal, err := dexUtils.ParseDelistProposalJSON(cliCtx.Codec, args[0])
			if err != nil {
				return err
			}

			from := cliCtx.GetFromAddress()
			content := types.NewDelistProposal(proposal.Title, proposal.Description, from, proposal.BaseAsset, proposal.QuoteAsset)
			msg, err := gov.NewMsgSubmitProposal(content, proposal.Deposit, from)
			if err != nil {
				return err
			}
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

}
