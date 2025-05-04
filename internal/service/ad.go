package service

import (
	"go.uber.org/zap"
	"slices"
	"sort"

	"sweng-task/internal/model"
)

// Scoring Weights
const (
	wCategory = 0.4
	wKeyword  = 0.4
	wBid      = 0.2
)

// AdService selects winning ads for placement
type AdService struct {
	lineItemService *LineItemService
	log             *zap.SugaredLogger
}

func NewAdService(lineItemService *LineItemService, log *zap.SugaredLogger) *AdService {
	return &AdService{
		lineItemService: lineItemService,
		log:             log,
	}
}

type AdQuery struct {
	Placement string
	Category  string
	Keyword   string
	Limit     int
}

func (s *AdService) GetWinningAds(q AdQuery) ([]*model.Ad, error) {
	lineItems, err := s.lineItemService.FindMatchingLineItems(q.Placement, q.Category, q.Keyword)
	if err != nil {
		return nil, err
	}
	if len(lineItems) == 0 {
		return []*model.Ad{}, nil
	}

	maxBid := lineItems[0].Bid
	for _, li := range lineItems {
		if li.Bid > maxBid {
			maxBid = li.Bid
		}
	}

	scores := make([]float64, len(lineItems))
	for i, li := range lineItems {
		var score float64
		if slices.Contains(li.Categories, q.Category) {
			score += wCategory
		}
		if slices.Contains(li.Keywords, q.Keyword) {
			score += wKeyword
		}
		if maxBid > 0 {
			score = (li.Bid / maxBid) * wBid
		}
		scores[i] = score
	}

	// Sort line items by descending score
	sort.Slice(lineItems, func(i, j int) bool {
		return scores[i] > scores[j]
	})

	var result []*model.Ad
	for _, li := range lineItems {
		result = append(result, &model.Ad{
			ID:           li.ID,
			Name:         li.Name,
			AdvertiserID: li.AdvertiserID,
			Bid:          li.Bid,
			Placement:    li.Placement,
			ServeURL:     "/ad/serve/" + li.ID,
		})
		if len(result) >= q.Limit {
			break
		}
	}
	return result, nil
}
