package rest


//type setNameReq struct {
//	BaseReq rest.BaseReq `json:"base_req"`
//	Name    string       `json:"name"`
//	Value   string       `json:"value"`
//	Owner   string       `json:"owner"`
//}

//func setNameHandler(cliCtx context.CLIContext) http.HandlerFunc {
	//return func(w http.ResponseWriter, r *http.Request) {
	//	var req setNameReq
	//	if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
	//		rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
	//		return
	//	}
	//
	//	baseReq := req.BaseReq.Sanitize()
	//	if !baseReq.ValidateBasic(w) {
	//		return
	//	}
	//
	//	addr, err := sdk.AccAddressFromBech32(req.Owner)
	//	if err != nil {
	//		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
	//		return
	//	}
	//
	//	// create the message
	//	msg := types.NewMsgSetName(req.Name, req.Value, addr)
	//	err = msg.ValidateBasic()
	//	if err != nil {
	//		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
	//		return
	//	}
	//
	//	utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	//}
//}

