package rpcClient_test

import (
	"fmt"
	rpcClient "hexa_micro/pkg/shortenservice/interface/rpcClient"
	"hexa_micro/pkg/shortenservice/mocks"
	"hexa_micro/pkg/shortenservice/model"
	"testing"
)

func TestRPCHandler_Find(t *testing.T) {
	mockRedirect := &model.Redirect{
		Code:     "XH2ey9WR",
		URL:      "http://vnexpress.net",
		CreateAt: 123456,
	}

	mockShortenUCase := new(mocks.ShortenUseCaseFake)
	h := rpcClient.NewRPCHandler(mockShortenUCase)
	type args struct {
		code  string
		reply *model.Redirect
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"Find Successfully",
			args{
				code:  mockRedirect.Code,
				reply: mockRedirect,
			},
			false,
		},
		{
			"Find Unsuccessfully 'cause of invalid code",
			args{
				code:  "",
				reply: mockRedirect,
			},
			true,
		},
		{
			"Not found",
			args{
				code:  "abc",
				reply: mockRedirect,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := h.Find(tt.args.code, tt.args.reply); (err != nil) != tt.wantErr {
				t.Errorf("RPCHandler.Find() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRPCHandler_Store(t *testing.T) {
	redirect := &model.Redirect{
		URL: "http://vnexpress.net",
	}

	invalidRedirect := &model.Redirect{
		URL: "http/vnexpress.net",
	}

	mockShortenUCase := new(mocks.ShortenUseCaseFake)
	h := rpcClient.NewRPCHandler(mockShortenUCase)

	type args struct {
		item  *model.Redirect
		reply *model.Redirect
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"Store Successfully",
			args{
				item:  redirect,
				reply: redirect,
			},
			false,
		},
		{
			"Store Unsuccessfully 'cause of invalid URL",
			args{
				item:  invalidRedirect,
				reply: invalidRedirect,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := h.Store(tt.args.item, tt.args.reply); (err != nil) != tt.wantErr {
				t.Errorf("RPCHandler.Store() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Printf("%v\n", tt.args.reply)
		})
	}
}
