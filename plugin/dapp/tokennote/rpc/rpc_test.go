// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc

import (
	"testing"

	"github.com/33cn/chain33/client/mocks"
	rpctypes "github.com/33cn/chain33/rpc/types"
	"github.com/33cn/chain33/types"
	tokenty "github.com/33cn/plugin/plugin/dapp/tokennote/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	context "golang.org/x/net/context"
)

func newTestChannelClient() *channelClient {
	api := &mocks.QueueProtocolAPI{}
	return &channelClient{
		ChannelClient: rpctypes.ChannelClient{QueueProtocolAPI: api},
	}
}

func newTestJrpcClient() *Jrpc {
	return &Jrpc{cli: newTestChannelClient()}
}

func testChannelClientGetTokenBalanceToken(t *testing.T) {
	api := new(mocks.QueueProtocolAPI)

	client := &channelClient{
		ChannelClient: rpctypes.ChannelClient{QueueProtocolAPI: api},
	}

	head := &types.Header{StateHash: []byte("sdfadasds")}
	api.On("GetLastHeader").Return(head, nil)

	var acc = &types.Account{Addr: "1Jn2qu84Z1SUUosWjySggBS9pKWdAP3tZt", Balance: 100}
	accv := types.Encode(acc)
	storevalue := &types.StoreReplyValue{}
	storevalue.Values = append(storevalue.Values, accv)
	api.On("StoreGet", mock.Anything).Return(storevalue, nil)

	var addrs = make([]string, 1)
	addrs = append(addrs, "1Jn2qu84Z1SUUosWjySggBS9pKWdAP3tZt")
	var in = &tokenty.ReqTokennoteBalance{
		Execer:      types.ExecName(tokenty.TokenX),
		Addresses:   addrs,
		TokennoteSymbol: "xxx",
	}
	data, err := client.GetTokennoteBalance(context.Background(), in)
	assert.Nil(t, err)
	accounts := data.Acc
	assert.Equal(t, acc.Addr, accounts[0].Addr)

}

func testChannelClientGetTokenBalanceOther(t *testing.T) {
	api := new(mocks.QueueProtocolAPI)
	client := &channelClient{
		ChannelClient: rpctypes.ChannelClient{QueueProtocolAPI: api},
	}

	head := &types.Header{StateHash: []byte("sdfadasds")}
	api.On("GetLastHeader").Return(head, nil)

	var acc = &types.Account{Addr: "1Jn2qu84Z1SUUosWjySggBS9pKWdAP3tZt", Balance: 100}
	accv := types.Encode(acc)
	storevalue := &types.StoreReplyValue{}
	storevalue.Values = append(storevalue.Values, accv)
	api.On("StoreGet", mock.Anything).Return(storevalue, nil)

	var addrs = make([]string, 1)
	addrs = append(addrs, "1Jn2qu84Z1SUUosWjySggBS9pKWdAP3tZt")
	var in = &tokenty.ReqTokennoteBalance{
		Execer:      types.ExecName("trade"),
		Addresses:   addrs,
		TokennoteSymbol: "xxx",
	}
	data, err := client.GetTokennoteBalance(context.Background(), in)
	assert.Nil(t, err)
	accounts := data.Acc
	assert.Equal(t, acc.Addr, accounts[0].Addr)

}

func TestChannelClientGetTokenBalance(t *testing.T) {
	testChannelClientGetTokenBalanceToken(t)
	testChannelClientGetTokenBalanceOther(t)

}
