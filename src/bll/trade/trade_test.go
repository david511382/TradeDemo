package trade

import (
	"testing"
	"time"
	"zerologix-homework/src/pkg/timeutil"
	"zerologix-homework/src/pkg/util"

	gomock "github.com/golang/mock/gomock"
)

func TestTradeLogic_Match(t *testing.T) {
	t.Parallel()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	type args struct {
		order Order
	}
	type migrations struct {
		timeUtilCreatorFn     func() timeutil.ITime
		tradeStorageCreatorFn func() ITradeStorage
	}
	type wants struct {
		err bool
	}
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			"匹配",
			args{
				Order{
					OrderType: ORDER_TYPE_SELL,
					Quantity:  10,
					Price:     5,
				},
			},
			migrations{
				timeUtilCreatorFn: func() timeutil.ITime {
					mockObj := timeutil.NewMockITime(mockCtl)
					return mockObj
				},
				tradeStorageCreatorFn: func() ITradeStorage {
					mockObj := NewMockITradeStorage(mockCtl)

					mockObj.EXPECT().LockOrders(float64(5)).Return(nil)

					var storeOrders = []*Order{
						// 跳過相同type
						{
							ID:        4,
							OrderType: ORDER_TYPE_SELL,
							Quantity:  10,
							Price:     5,
							Timestamp: 4,
						},
						// 數量不夠
						{
							ID:        1,
							OrderType: ORDER_TYPE_BUY,
							Quantity:  1,
							Price:     5,
							Timestamp: 1,
						},
						// 數量不夠
						{
							ID:        2,
							OrderType: ORDER_TYPE_BUY,
							Quantity:  1,
							Price:     5,
							Timestamp: 2,
						},
						// 扣完有剩
						{
							ID:        3,
							OrderType: ORDER_TYPE_BUY,
							Quantity:  9,
							Price:     5,
							Timestamp: 3,
						},
						// 其它的
						{
							ID:        5,
							OrderType: ORDER_TYPE_BUY,
							Quantity:  10,
							Price:     5,
							Timestamp: 5,
						},
					}
					mockObj.EXPECT().Load(float64(5)).Return(storeOrders, nil)

					// 被排序
					var (
						price  float64 = 5
						orders         = []*Order{
							// 扣完有剩
							{
								ID:        3,
								OrderType: ORDER_TYPE_BUY,
								Quantity:  1,
								Price:     5,
								Timestamp: 3,
							},
							// 跳過相同type
							{
								ID:        4,
								OrderType: ORDER_TYPE_SELL,
								Quantity:  10,
								Price:     5,
								Timestamp: 4,
							},
							// 其它的
							{
								ID:        5,
								OrderType: ORDER_TYPE_BUY,
								Quantity:  10,
								Price:     5,
								Timestamp: 5,
							},
						}
					)
					mockObj.EXPECT().Set(price, orders).Return(nil)

					mockObj.EXPECT().ReleaseOrdersLock(float64(5))

					return mockObj
				},
			},
			wants{
				err: false,
			},
		},
		{
			"新增",
			args{
				Order{
					OrderType: ORDER_TYPE_SELL,
					Quantity:  10,
					Price:     5,
				},
			},
			migrations{
				timeUtilCreatorFn: func() timeutil.ITime {
					mockObj := timeutil.NewMockITime(mockCtl)

					mockObj.EXPECT().Now().Return(*util.GetTimePLoc(time.Local, 2013, 8, 2, 1, 2, 3, 4567))

					return mockObj
				},
				tradeStorageCreatorFn: func() ITradeStorage {
					mockObj := NewMockITradeStorage(mockCtl)

					mockObj.EXPECT().LockOrders(float64(5)).Return(nil)

					mockObj.EXPECT().Load(float64(5)).Return(nil, nil)

					mockObj.EXPECT().LockID().Return(nil)

					mockObj.EXPECT().LoadID().Return(uint(1), nil)

					mockObj.EXPECT().UpdateID(uint(2)).Return(nil)

					mockObj.EXPECT().ReleaseIDLock()

					var (
						price      float64 = 5
						wantOrders         = []*Order{
							// 新增
							{
								ID:        2,
								OrderType: ORDER_TYPE_SELL,
								Quantity:  10,
								Price:     5,
								Timestamp: 1375376523000004567,
							},
						}
					)
					mockObj.EXPECT().Set(price, wantOrders).Return(nil)

					mockObj.EXPECT().ReleaseOrdersLock(float64(5))

					return mockObj
				},
			},
			wants{
				err: false,
			},
		},
		{
			"不匹配",
			args{
				Order{
					OrderType: ORDER_TYPE_SELL,
					Quantity:  10,
					Price:     5,
				},
			},
			migrations{
				timeUtilCreatorFn: func() timeutil.ITime {
					mockObj := timeutil.NewMockITime(mockCtl)

					mockObj.EXPECT().Now().Return(*util.GetTimePLoc(time.Local, 2013, 8, 2, 1, 2, 3, 4567))

					return mockObj
				},
				tradeStorageCreatorFn: func() ITradeStorage {
					mockObj := NewMockITradeStorage(mockCtl)

					mockObj.EXPECT().LockOrders(float64(5)).Return(nil)

					var storeOrders = []*Order{
						// 不匹配
						{
							ID:        4,
							OrderType: ORDER_TYPE_SELL,
							Quantity:  10,
							Price:     5,
							Timestamp: 4,
						},
					}
					mockObj.EXPECT().Load(float64(5)).Return(storeOrders, nil)

					mockObj.EXPECT().LockID().Return(nil)

					mockObj.EXPECT().LoadID().Return(uint(1), nil)

					mockObj.EXPECT().UpdateID(uint(2)).Return(nil)

					mockObj.EXPECT().ReleaseIDLock()

					var (
						price      float64 = 5
						wantOrders         = []*Order{
							// 不匹配
							{
								ID:        4,
								OrderType: ORDER_TYPE_SELL,
								Quantity:  10,
								Price:     5,
								Timestamp: 4,
							},
							// 新增
							{
								ID:        2,
								OrderType: ORDER_TYPE_SELL,
								Quantity:  10,
								Price:     5,
								Timestamp: 1375376523000004567,
							},
						}
					)
					mockObj.EXPECT().Set(price, wantOrders).Return(nil)

					mockObj.EXPECT().ReleaseOrdersLock(float64(5))

					return mockObj
				},
			},
			wants{
				err: false,
			},
		},
		{
			"匹配完",
			args{
				Order{
					OrderType: ORDER_TYPE_SELL,
					Quantity:  10,
					Price:     5,
				},
			},
			migrations{
				timeUtilCreatorFn: func() timeutil.ITime {
					mockObj := timeutil.NewMockITime(mockCtl)
					return mockObj
				},
				tradeStorageCreatorFn: func() ITradeStorage {
					mockObj := NewMockITradeStorage(mockCtl)

					mockObj.EXPECT().LockOrders(float64(5)).Return(nil)

					var storeOrders = []*Order{
						// 匹配
						{
							ID:        4,
							OrderType: ORDER_TYPE_BUY,
							Quantity:  10,
							Price:     5,
							Timestamp: 4,
						},
					}
					mockObj.EXPECT().Load(float64(5)).Return(storeOrders, nil)

					var (
						price float64 = 5
					)
					mockObj.EXPECT().Delete(price).Return(nil)

					mockObj.EXPECT().ReleaseOrdersLock(float64(5))

					return mockObj
				},
			},
			wants{
				err: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeUtil := tt.migrations.timeUtilCreatorFn()
			iTradeStorage := tt.migrations.tradeStorageCreatorFn()

			l := NewTradeLogic(
				timeUtil,
				iTradeStorage,
			)
			err := l.Match(tt.args.order)
			if (err != nil) != tt.wants.err {
				t.Errorf("TradeLogic.Match() error = %v, wantErr %v", err, tt.wants.err)
				return
			}
		})
	}
}
