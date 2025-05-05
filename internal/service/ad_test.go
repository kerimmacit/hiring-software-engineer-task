package service

import (
	"reflect"
	"testing"

	"sweng-task/internal/model"
)

func TestAdService_winningAdCalculator(t *testing.T) {
	li1 := &model.LineItem{ID: "1", Name: "1", AdvertiserID: "1", Bid: 10, Budget: 0, Placement: "pl1", Categories: nil, Keywords: nil}
	li2 := &model.LineItem{ID: "2", Name: "2", AdvertiserID: "2", Bid: 20, Budget: 0, Placement: "pl2", Categories: nil, Keywords: nil}

	type args struct {
		q         AdQuery
		lineItems []*model.LineItem
	}
	tests := []struct {
		name string
		args args
		want []*model.LineItem
	}{
		{
			name: "for nil lineItems, winningAdCalculator returns empty ads",
			args: args{
				q:         AdQuery{Placement: "", Category: "", Keyword: "", Limit: 2},
				lineItems: nil,
			},
			want: []*model.LineItem{},
		},
		{
			name: "for empty lineItems, winningAdCalculator returns empty ads",
			args: args{
				q:         AdQuery{Placement: "", Category: "", Keyword: "", Limit: 2},
				lineItems: []*model.LineItem{},
			},
			want: []*model.LineItem{},
		},
		{
			name: "if other fields rating same, winningAdCalculator should return higher bid first",
			args: args{
				q:         AdQuery{Placement: "", Category: "", Keyword: "", Limit: 2},
				lineItems: []*model.LineItem{li1, li2},
			},
			want: []*model.LineItem{li2, li1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AdService{}
			got := s.winningAdCalculator(tt.args.q, tt.args.lineItems)
			if !isAdSlicesEqual(got, tt.want) {
				t.Errorf("winningAdCalculator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func isAdSlicesEqual(got []*model.LineItem, want []*model.LineItem) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if !reflect.DeepEqual(*got[i], *want[i]) {
			return false
		}
	}
	return true
}
