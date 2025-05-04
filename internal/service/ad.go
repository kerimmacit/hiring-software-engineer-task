package service

import (
	"go.uber.org/zap"

	"sweng-task/internal/model"
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
