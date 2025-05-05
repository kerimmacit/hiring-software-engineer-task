package service

import (
	"reflect"
	"sweng-task/internal/model"
	"testing"
)

func TestAdService_winningAdCalculator(t *testing.T) {
	type args struct {
		q         AdQuery
		lineItems []*model.LineItem
	}
	tests := []struct {
		name string
		args args
		want []*model.Ad
	}{
		{
			name: "for nil lineItems, winningAdCalculator returns empty ads",
			args: args{
				q:         AdQuery{Placement: "", Category: "", Keyword: "", Limit: 2},
				lineItems: nil,
			},
			want: []*model.Ad{},
		},
		{
			name: "for empty lineItems, winningAdCalculator returns empty ads",
			args: args{
				q:         AdQuery{Placement: "", Category: "", Keyword: "", Limit: 2},
				lineItems: []*model.LineItem{},
			},
			want: []*model.Ad{},
		},
		{
			name: "if other fields rating same, winningAdCalculator should return higher bid first",
			args: args{
				q: AdQuery{Placement: "", Category: "", Keyword: "", Limit: 2},
				lineItems: []*model.LineItem{
					{ID: "1", Name: "1", AdvertiserID: "1", Bid: 10, Budget: 0, Placement: "pl1", Categories: nil, Keywords: nil},
					{ID: "2", Name: "2", AdvertiserID: "2", Bid: 20, Budget: 0, Placement: "pl2", Categories: nil, Keywords: nil},
				},
			},
			want: []*model.Ad{
				{ID: "2", Name: "2", AdvertiserID: "2", Bid: 20, Placement: "pl2", ServeURL: serveUrlGenerator(&model.LineItem{ID: "2"})},
				{ID: "1", Name: "1", AdvertiserID: "1", Bid: 10, Placement: "pl1", ServeURL: serveUrlGenerator(&model.LineItem{ID: "1"})},
			},
		},
		{
			name: "limit parameter should limit winningAdCalculator return slice length",
			args: args{
				q: AdQuery{Placement: "", Category: "", Keyword: "", Limit: 2},
				lineItems: []*model.LineItem{
					{ID: "1", Name: "1", AdvertiserID: "1", Bid: 10, Placement: "pl1"},
					{ID: "2", Name: "2", AdvertiserID: "2", Bid: 20, Placement: "pl2"},
					{ID: "3", Name: "3", AdvertiserID: "3", Bid: 30, Placement: "pl3"},
					{ID: "4", Name: "4", AdvertiserID: "4", Bid: 40, Placement: "pl4"},
				},
			},
			want: []*model.Ad{
				{ID: "4", Name: "4", AdvertiserID: "4", Bid: 40, Placement: "pl4", ServeURL: serveUrlGenerator(&model.LineItem{ID: "4"})},
				{ID: "3", Name: "3", AdvertiserID: "3", Bid: 30, Placement: "pl3", ServeURL: serveUrlGenerator(&model.LineItem{ID: "3"})},
			},
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

func isAdSlicesEqual(got []*model.Ad, want []*model.Ad) bool {
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
