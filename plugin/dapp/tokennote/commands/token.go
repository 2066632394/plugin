// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"github.com/33cn/chain33/rpc/jsonclient"
	rpctypes "github.com/33cn/chain33/rpc/types"
	"github.com/33cn/chain33/system/dapp/commands"
	"github.com/33cn/chain33/types"
	tokennotety "github.com/33cn/plugin/plugin/dapp/tokennote/types"
	"github.com/spf13/cobra"
)

var (
	tokennoteSymbol string
)

// TokennoteCmd tokennote 命令行
func TokennoteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokennote",
		Short: "Tokennote management",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(
		CreateTokennoteTransferCmd(),
		CreateTokennoteWithdrawCmd(),
		GetTokennotesFinishCreatedCmd(),
		//GetTokennoteAssetsCmd(),
		//GetTokennoteBalanceCmd(),
		CreateRawTokennoteCreateTxCmd(),
		CreateTokennoteTransferExecCmd(),
	)

	return cmd
}

// CreateTokennoteTransferCmd create raw transfer tx
func CreateTokennoteTransferCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer",
		Short: "Create a tokennote transfer transaction",
		Run:   createTokennoteTransfer,
	}
	addCreateTokennoteTransferFlags(cmd)
	return cmd
}

func addCreateTokennoteTransferFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("to", "t", "", "receiver account address")
	cmd.MarkFlagRequired("to")
	cmd.Flags().Float64P("amount", "a", 0, "transaction amount")
	cmd.MarkFlagRequired("amount")
	cmd.Flags().StringP("note", "n", "", "transaction note info")
	cmd.Flags().StringP("symbol", "s", "", "tokennote symbol")
	cmd.MarkFlagRequired("symbol")
}

func createTokennoteTransfer(cmd *cobra.Command, args []string) {
	commands.CreateAssetTransfer(cmd, args, tokennotety.TokennoteX)
}

// CreateTokennoteTransferExecCmd create raw transfer tx
func CreateTokennoteTransferExecCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send_exec",
		Short: "Create a tokennote send to executor transaction",
		Run:   createTokennoteSendToExec,
	}
	addCreateTokennoteSendToExecFlags(cmd)
	return cmd
}

func addCreateTokennoteSendToExecFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("exec", "e", "", "receiver executor address")
	cmd.MarkFlagRequired("exec")

	cmd.Flags().Float64P("amount", "a", 0, "transaction amount")
	cmd.MarkFlagRequired("amount")

	cmd.Flags().StringP("note", "n", "", "transaction note info")

	cmd.Flags().StringP("symbol", "s", "", "tokennote symbol")
	cmd.MarkFlagRequired("symbol")
}

func createTokennoteSendToExec(cmd *cobra.Command, args []string) {
	commands.CreateAssetSendToExec(cmd, args, tokennotety.TokennoteX)
}

// CreateTokennoteWithdrawCmd create raw withdraw tx
func CreateTokennoteWithdrawCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw",
		Short: "Create a tokennote withdraw transaction",
		Run:   createTokennoteWithdraw,
	}
	addCreateTokennoteWithdrawFlags(cmd)
	return cmd
}

func addCreateTokennoteWithdrawFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("exec", "e", "", "execer withdrawn from")
	cmd.MarkFlagRequired("exec")

	cmd.Flags().Float64P("amount", "a", 0, "withdraw amount")
	cmd.MarkFlagRequired("amount")

	cmd.Flags().StringP("note", "n", "", "transaction note info")

	cmd.Flags().StringP("symbol", "s", "", "tokennote symbol")
	cmd.MarkFlagRequired("symbol")
}

func createTokennoteWithdraw(cmd *cobra.Command, args []string) {
	commands.CreateAssetWithdraw(cmd, args, tokennotety.TokennoteX)
}


// GetTokennotesFinishCreatedCmd get finish created tokennotes
func GetTokennotesFinishCreatedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get_finish_created",
		Short: "Get finish created tokennotes",
		Run:   getFinishCreatedTokennotes,
	}
	return cmd
}

func getFinishCreatedTokennotes(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	paraName, _ := cmd.Flags().GetString("paraName")
	var reqtokennotes tokennotety.ReqTokennotes
	reqtokennotes.Status = tokennotety.TokennoteStatusCreated
	reqtokennotes.QueryAll = true
	var params rpctypes.Query4Jrpc
	params.Execer = getRealExecName(paraName, "tokennote")
	params.FuncName = "GetTokennotes"
	params.Payload = types.MustPBToJSON(&reqtokennotes)
	rpc, err := jsonclient.NewJSONClient(rpcLaddr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	var res tokennotety.ReplyTokennotes
	err = rpc.Call("Chain33.Query", params, &res)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	for _, createdTokennote := range res.Tokens {
		//createdTokennote.Price = createdTokennote.Price / types.Coin
		//createdTokennote.Total = createdTokennote.Total / types.TokennotePrecision

		//fmt.Printf("---The %dth Finish Created tokennote is below--------------------\n", i)
		data, err := json.MarshalIndent(createdTokennote, "", "    ")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		fmt.Println(string(data))
	}
}

//// GetTokennoteAssetsCmd get tokennote assets
//func GetTokennoteAssetsCmd() *cobra.Command {
//	cmd := &cobra.Command{
//		Use:   "tokennote_assets",
//		Short: "Get tokennote assets",
//		Run:   tokennoteAssets,
//	}
//	addTokennoteAssetsFlags(cmd)
//	return cmd
//}
//
//func addTokennoteAssetsFlags(cmd *cobra.Command) {
//	cmd.Flags().StringP("exec", "e", "", "execer name")
//	cmd.MarkFlagRequired("exec")
//
//	cmd.Flags().StringP("addr", "a", "", "account address")
//	cmd.MarkFlagRequired("addr")
//}
//
//func tokennoteAssets(cmd *cobra.Command, args []string) {
//	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
//	paraName, _ := cmd.Flags().GetString("paraName")
//	addr, _ := cmd.Flags().GetString("addr")
//	execer, _ := cmd.Flags().GetString("exec")
//	execer = getRealExecName(paraName, execer)
//	req := tokennotety.ReqAccountTokennoteAssets{
//		Address: addr,
//		Execer:  execer,
//	}
//	var params rpctypes.Query4Jrpc
//	params.Execer = getRealExecName(paraName, "tokennote")
//	params.FuncName = "GetAccountTokennoteAssets"
//	params.Payload = types.MustPBToJSON(&req)
//
//	var res tokennotety.ReplyAccountTokennoteAssets
//	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.Query", params, &res)
//	ctx.SetResultCb(parseTokennoteAssetsRes)
//	ctx.Run()
//}
//
//func parseTokennoteAssetsRes(arg interface{}) (interface{}, error) {
//	res := arg.(*tokennotety.ReplyAccountTokennoteAssets)
//	var result []*tokennotety.TokennoteAccountResult
//	for _, ta := range res.TokennoteAssets {
//		balanceResult := strconv.FormatFloat(float64(ta.Account.Balance)/float64(types.TokennotePrecision), 'f', 4, 64)
//		frozenResult := strconv.FormatFloat(float64(ta.Account.Frozen)/float64(types.TokennotePrecision), 'f', 4, 64)
//		tokennoteAccount := &tokennotety.TokennoteAccountResult{
//			Tokennote:    ta.Symbol,
//			Addr:     ta.Account.Addr,
//			Currency: ta.Account.Currency,
//			Balance:  balanceResult,
//			Frozen:   frozenResult,
//		}
//		result = append(result, tokennoteAccount)
//	}
//	return result, nil
//}
//
//// GetTokennoteBalanceCmd get tokennote balance
//func GetTokennoteBalanceCmd() *cobra.Command {
//	cmd := &cobra.Command{
//		Use:   "tokennote_balance",
//		Short: "Get tokennote balance of one or more addresses",
//		Run:   tokennoteBalance,
//	}
//	addTokennoteBalanceFlags(cmd)
//	return cmd
//}
//
//func addTokennoteBalanceFlags(cmd *cobra.Command) {
//	cmd.Flags().StringVarP(&tokennoteSymbol, "symbol", "s", "", "tokennote symbol")
//	cmd.MarkFlagRequired("symbol")
//
//	cmd.Flags().StringP("exec", "e", "", "execer name")
//	cmd.MarkFlagRequired("exec")
//
//	cmd.Flags().StringP("address", "a", "", "account addresses, separated by space")
//	cmd.MarkFlagRequired("address")
//}
//
//func tokennoteBalance(cmd *cobra.Command, args []string) {
//	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
//	addr, _ := cmd.Flags().GetString("address")
//	tokennote, _ := cmd.Flags().GetString("symbol")
//	execer, _ := cmd.Flags().GetString("exec")
//	paraName, _ := cmd.Flags().GetString("paraName")
//	execer = getRealExecName(paraName, execer)
//	addresses := strings.Split(addr, " ")
//	params := tokennotety.ReqTokennoteBalance{
//		Addresses:   addresses,
//		TokennoteSymbol: tokennote,
//		Execer:      execer,
//	}
//	var res []*rpctypes.Account
//	ctx := jsonclient.NewRPCCtx(rpcLaddr, "tokennote.GetTokennoteBalance", params, &res)
//	ctx.SetResultCb(parseTokennoteBalanceRes)
//	ctx.Run()
//}
//
//func parseTokennoteBalanceRes(arg interface{}) (interface{}, error) {
//	res := arg.(*[]*rpctypes.Account)
//	var result []*tokennotety.TokennoteAccountResult
//	for _, one := range *res {
//		balanceResult := strconv.FormatFloat(float64(one.Balance)/float64(types.TokennotePrecision), 'f', 4, 64)
//		frozenResult := strconv.FormatFloat(float64(one.Frozen)/float64(types.TokennotePrecision), 'f', 4, 64)
//		tokennoteAccount := &tokennotety.TokennoteAccountResult{
//			Tokennote:    tokennoteSymbol,
//			Addr:     one.Addr,
//			Currency: one.Currency,
//			Balance:  balanceResult,
//			Frozen:   frozenResult,
//		}
//		result = append(result, tokennoteAccount)
//	}
//	return result, nil
//}

// CreateRawTokennoteCreateTxCmd create raw tokennote precreate transaction
func CreateRawTokennoteCreateTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a tokennote transaction",
		Run:   tokennoteCreated,
	}
	addTokennoteCreatedFlags(cmd)
	return cmd
}

func addTokennoteCreatedFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("issuer", "isa", "", "address of tokennote owner")
	cmd.MarkFlagRequired("issuer")

	cmd.Flags().StringP("issuer_name", "isn", "", "tokennote Issuer name")
	cmd.MarkFlagRequired("issuer_name")

	cmd.Flags().StringP("currency", "s", "", "tokennote currency")
	cmd.MarkFlagRequired("currency")

	cmd.Flags().StringP("issuer_phone", "isp", "", "tokennote issuer phone")

	cmd.Flags().StringP("issuer_id", "isi", "", "tokennote issuer id ")

	cmd.Flags().Int64P("rate", "r", 0, "tokennote rate")
	cmd.MarkFlagRequired("rate")

	cmd.Flags().Int64P("total", "t", 0, "total amount of the tokennote")
	cmd.MarkFlagRequired("total")

	cmd.Flags().Int64P("acceptdate", "d", 0, "tokennote accept date")
	cmd.MarkFlagRequired("total")

	cmd.Flags().StringP("introduction", "intr", "", "tokennote introduction")

}

func tokennoteCreated(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	issuer, _ := cmd.Flags().GetString("issuer")
	issuerName, _ := cmd.Flags().GetString("issuer_name")
	issuerPhone, _ := cmd.Flags().GetString("issuer_phone")
	issuerId, _ := cmd.Flags().GetString("issuer_id")
	currency, _ := cmd.Flags().GetString("currency")
	introduction, _ := cmd.Flags().GetString("introduction")
	rate, _ := cmd.Flags().GetInt64("rate")
	total, _ := cmd.Flags().GetInt64("total")
	acceptdate, _ := cmd.Flags().GetInt64("acceptdate")

	params := &tokennotety.TokennoteCreate{
		Issuer:        issuer,
		IssuerPhone:   issuerPhone,
		IssuerId:      issuerId,
		IssuerName:         issuerName,
		Acceptor:       issuer,
		AcceptanceDate: acceptdate,
		Currency:       currency,
		Balance:        total,
		Introduction: introduction,
		Rate:        rate,
		Repayamount:     0,
	}
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "tokennote.CreateRawTokennotePreCreateTx", params, nil)
	ctx.RunWithoutMarshal()
}

