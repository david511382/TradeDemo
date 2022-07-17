package api

import (
	"fmt"
	"testing"
	tradeBll "zerologix-homework/src/bll/trade"
	"zerologix-homework/src/pkg/util"
	"zerologix-homework/src/server/reqs"
	"zerologix-homework/src/server/resp"

	gomock "github.com/golang/mock/gomock"
)

func TestOrder_PostBuy(t *testing.T) {
	t.Parallel()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	type args struct {
		req *reqs.OrderPostBuy
	}
	type migrations struct {
		tradeLogicCreatorFn func() tradeBll.ITradeLogic
	}
	type wants struct {
		result resp.OrderPostBuy
		err    bool
	}
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			"執行成功",
			args{
				&reqs.OrderPostBuy{
					Quantity: 1,
					Price:    2,
				},
			},
			migrations{
				tradeLogicCreatorFn: func() tradeBll.ITradeLogic {
					mockObj := tradeBll.NewMockITradeLogic(mockCtl)

					var (
						order = tradeBll.Order{
							ID:        0,
							OrderType: tradeBll.ORDER_TYPE_BUY,
							Quantity:  1,
							Price:     2,
							Timestamp: 0,
						}
					)
					mockObj.EXPECT().Match(
						order,
					).Return(nil)

					return mockObj
				},
			},
			wants{
				result: resp.OrderPostBuy{
					Message: "OK",
				},
				err: false,
			},
		},
		{
			"執行失敗",
			args{
				&reqs.OrderPostBuy{},
			},
			migrations{
				tradeLogicCreatorFn: func() tradeBll.ITradeLogic {
					mockObj := tradeBll.NewMockITradeLogic(mockCtl)

					var (
						order = tradeBll.Order{
							OrderType: tradeBll.ORDER_TYPE_BUY,
						}
					)
					mockObj.EXPECT().Match(
						order,
					).Return(fmt.Errorf("error"))

					return mockObj
				},
			},
			wants{
				result: resp.OrderPostBuy{
					Message: "Error",
				},
				err: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iTradeLogic := tt.migrations.tradeLogicCreatorFn()

			l := NewOrder(
				iTradeLogic,
			)
			gotResult, err := l.PostBuy(tt.args.req)
			if (err != nil) != tt.wants.err {
				t.Errorf("Order.PostBuy() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if ok, msg := util.Comp(gotResult, tt.wants.result); !ok {
				t.Error(msg)
				return
			}
		})
	}
}

func TestOrder_PostSell(t *testing.T) {
	t.Parallel()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	type args struct {
		req *reqs.OrderPostSell
	}
	type migrations struct {
		tradeLogicCreatorFn func() tradeBll.ITradeLogic
	}
	type wants struct {
		result resp.OrderPostSell
		err    bool
	}
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			"執行成功",
			args{
				&reqs.OrderPostSell{
					OrderPostBuy: reqs.OrderPostBuy{
						Quantity: 1,
						Price:    2,
					},
				},
			},
			migrations{
				tradeLogicCreatorFn: func() tradeBll.ITradeLogic {
					mockObj := tradeBll.NewMockITradeLogic(mockCtl)

					var (
						order = tradeBll.Order{
							ID:        0,
							OrderType: tradeBll.ORDER_TYPE_SELL,
							Quantity:  1,
							Price:     2,
							Timestamp: 0,
						}
					)
					mockObj.EXPECT().Match(
						order,
					).Return(nil)

					return mockObj
				},
			},
			wants{
				result: resp.OrderPostSell{
					OrderPostBuy: resp.OrderPostBuy{
						Message: "OK",
					},
				},
				err: false,
			},
		},
		{
			"執行失敗",
			args{
				&reqs.OrderPostSell{},
			},
			migrations{
				tradeLogicCreatorFn: func() tradeBll.ITradeLogic {
					mockObj := tradeBll.NewMockITradeLogic(mockCtl)

					var (
						order = tradeBll.Order{
							OrderType: tradeBll.ORDER_TYPE_SELL,
						}
					)
					mockObj.EXPECT().Match(
						order,
					).Return(fmt.Errorf("error"))

					return mockObj
				},
			},
			wants{
				result: resp.OrderPostSell{
					OrderPostBuy: resp.OrderPostBuy{
						Message: "Error",
					},
				},
				err: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iTradeLogic := tt.migrations.tradeLogicCreatorFn()

			l := NewOrder(
				iTradeLogic,
			)
			gotResult, err := l.PostSell(tt.args.req)
			if (err != nil) != tt.wants.err {
				t.Errorf("Order.PostBuy() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
			if ok, msg := util.Comp(gotResult, tt.wants.result); !ok {
				t.Error(msg)
				return
			}
		})
	}
}
