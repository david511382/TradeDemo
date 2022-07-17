package util

import (
	"testing"
)

func TestAscTimeRanges_Append(t *testing.T) {
	type args struct {
		newInsertTimeRange TimeRange
	}
	tests := []struct {
		name              string
		trs               AscTimeRanges
		args              args
		wantNewTimeRanges AscTimeRanges
	}{
		{
			"insert first",
			AscTimeRanges{
				{
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 15),
				},
			},
			args{
				newInsertTimeRange: TimeRange{
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 14),
				},
			},
			AscTimeRanges{
				{
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 14),
				},
				{
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 15),
				},
			},
		},
		{
			"insert mid",
			AscTimeRanges{
				{
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 3),
				},
				{
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 15),
				},
			},
			args{
				newInsertTimeRange: TimeRange{
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 14),
				},
			},
			AscTimeRanges{
				{
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 3),
				},
				{
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 14),
				},
				{
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 15),
				},
			},
		},
		{
			"insert last",
			AscTimeRanges{
				{
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 3),
				},
			},
			args{
				newInsertTimeRange: TimeRange{
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 14),
				},
			},
			AscTimeRanges{
				{
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 3),
				},
				{
					From: GetUTCTime(2013, 8, 2),
					To:   GetUTCTime(2013, 8, 14),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNewTimeRanges := tt.trs.Append(tt.args.newInsertTimeRange)
			if ok, msg := Comp(gotNewTimeRanges, tt.wantNewTimeRanges); !ok {
				t.Fatal(msg)
			}
		})
	}
}

func TestAscTimeRanges_CombineByCount(t *testing.T) {
	type migrations struct {
		trs AscTimeRanges
	}
	type wants struct {
		countAscTimeRangesMap map[int]AscTimeRanges
	}
	tests := []struct {
		name string
		migrations
		wants
	}{
		{
			"總測試",
			migrations{
				trs: NewAscTimeRanges(
					[]TimeRange{
						{
							From: *GetTimePLoc(nil, 2013, 8, 2, 0),
							To:   *GetTimePLoc(nil, 2013, 8, 2, 5),
						},
						{
							From: *GetTimePLoc(nil, 2013, 8, 2, 0),
							To:   *GetTimePLoc(nil, 2013, 8, 2, 6),
						},
						{
							From: *GetTimePLoc(nil, 2013, 8, 2, 1),
							To:   *GetTimePLoc(nil, 2013, 8, 2, 4),
						},
						{
							From: *GetTimePLoc(nil, 2013, 8, 2, 2),
							To:   *GetTimePLoc(nil, 2013, 8, 2, 3),
						},
					},
				),
			},
			wants{
				map[int]AscTimeRanges{
					1: {
						{
							From: *GetTimePLoc(nil, 2013, 8, 2, 5),
							To:   *GetTimePLoc(nil, 2013, 8, 2, 6),
						},
					},
					2: {
						{
							From: *GetTimePLoc(nil, 2013, 8, 2, 0),
							To:   *GetTimePLoc(nil, 2013, 8, 2, 1),
						},
						{
							From: *GetTimePLoc(nil, 2013, 8, 2, 4),
							To:   *GetTimePLoc(nil, 2013, 8, 2, 5),
						},
					},
					3: {
						{
							From: *GetTimePLoc(nil, 2013, 8, 2, 1),
							To:   *GetTimePLoc(nil, 2013, 8, 2, 2),
						},
						{
							From: *GetTimePLoc(nil, 2013, 8, 2, 3),
							To:   *GetTimePLoc(nil, 2013, 8, 2, 4),
						},
					},
					4: {
						{
							From: *GetTimePLoc(nil, 2013, 8, 2, 2),
							To:   *GetTimePLoc(nil, 2013, 8, 2, 3),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCountAscTimeRangesMap := tt.migrations.trs.CombineByCount()
			if ok, msg := Comp(gotCountAscTimeRangesMap, tt.wants.countAscTimeRangesMap); !ok {
				t.Error(msg)
				return
			}
		})
	}
}

func TestAscTimeRanges_Sub(t *testing.T) {
	t.Parallel()

	type args struct {
		sub TimeRange
	}
	type migrations struct {
		trs AscTimeRanges
	}
	type wants struct {
		newTimeRanges          AscTimeRanges
		intersectionTimeRanges AscTimeRanges
		leftTimeRanges         AscTimeRanges
	}
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			"not enough, over",
			args{
				sub: TimeRange{
					From: GetUTCTime(2013, 2),
					To:   GetUTCTime(2013, 6),
				},
			},
			migrations{
				AscTimeRanges{
					{
						From: GetUTCTime(2013, 1),
						To:   GetUTCTime(2013, 3),
					},
					{
						From: GetUTCTime(2013, 4),
						To:   GetUTCTime(2013, 5),
					},
				},
			},
			wants{
				newTimeRanges: AscTimeRanges{
					{
						From: GetUTCTime(2013, 1),
						To:   GetUTCTime(2013, 2),
					},
				},
				intersectionTimeRanges: AscTimeRanges{
					{
						From: GetUTCTime(2013, 2),
						To:   GetUTCTime(2013, 3),
					},
					{
						From: GetUTCTime(2013, 4),
						To:   GetUTCTime(2013, 5),
					},
				},
				leftTimeRanges: AscTimeRanges{
					{
						From: GetUTCTime(2013, 3),
						To:   GetUTCTime(2013, 4),
					},
					{
						From: GetUTCTime(2013, 5),
						To:   GetUTCTime(2013, 6),
					},
				},
			},
		},
		{
			"combine to sub",
			args{
				sub: TimeRange{
					From: GetUTCTime(2013, 1),
					To:   GetUTCTime(2013, 4),
				},
			},
			migrations{
				AscTimeRanges{
					{
						From: GetUTCTime(2013, 1),
						To:   GetUTCTime(2013, 3),
					},
					{
						From: GetUTCTime(2013, 2),
						To:   GetUTCTime(2013, 4),
					},
				},
			},
			wants{
				newTimeRanges: AscTimeRanges{
					{
						From: GetUTCTime(2013, 2),
						To:   GetUTCTime(2013, 3),
					},
				},
				intersectionTimeRanges: AscTimeRanges{
					{
						From: GetUTCTime(2013, 1),
						To:   GetUTCTime(2013, 3),
					},
					{
						From: GetUTCTime(2013, 3),
						To:   GetUTCTime(2013, 4),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNewTimeRanges, intersectionTimeRanges, leftTimeRanges := tt.migrations.trs.Sub(tt.args.sub)
			if ok, msg := Comp(gotNewTimeRanges, tt.wants.newTimeRanges); !ok {
				t.Error(msg)
				return
			}
			if ok, msg := Comp(intersectionTimeRanges, tt.wants.intersectionTimeRanges); !ok {
				t.Error(msg)
				return
			}
			if ok, msg := Comp(leftTimeRanges, tt.wants.leftTimeRanges); !ok {
				t.Error(msg)
				return
			}
		})
	}
}

func BenchmarkSubCourtDetails(b *testing.B) {
	type args struct {
		sub TimeRange
	}
	type migrations struct {
		trs AscTimeRanges
	}
	arg := args{
		sub: TimeRange{
			From: GetUTCTime(2013, 2),
			To:   GetUTCTime(2013, 6),
		},
	}
	migration := migrations{
		AscTimeRanges{
			{
				From: GetUTCTime(2013, 1),
				To:   GetUTCTime(2013, 3),
			},
			{
				From: GetUTCTime(2013, 4),
				To:   GetUTCTime(2013, 5),
			},
		},
	}
	for i := 0; i < b.N; i++ {
		_, _, _ = migration.trs.Sub(arg.sub)
	}
}
