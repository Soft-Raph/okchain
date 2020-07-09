package rest

import (
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/okex/okchain/x/gov"
	govrest "github.com/okex/okchain/x/gov/client/rest"
	paramscutils "github.com/okex/okchain/x/params/client/utils"
)

// ProposalRESTHandler returns a ProposalRESTHandler that exposes the param change REST handler with a given sub-route
func ProposalRESTHandler(cliCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "param_change",
		Handler:  postProposalHandlerFn(cliCtx),
	}
}

func postProposalHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req paramscutils.ParamChangeProposalReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		content := proposal.NewParameterChangeProposal(req.Title, req.Description, req.Changes.ToParamChanges())

		msg, err := gov.NewMsgSubmitProposal(content, req.Deposit, req.Proposer)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		authclient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
