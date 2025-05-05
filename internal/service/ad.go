package service

import (
	"errors"
	"slices"
	"sort"

	"go.uber.org/zap"

	"sweng-task/internal/domain_errors"
	"sweng-task/internal/model"
	"sweng-task/internal/repo"
)

// relevancy scoring helper struct
type scoredItem struct {
	*model.LineItem
	score float64
}

// Scoring Weights
const (
	wCategory = 0.3
	wKeyword  = 0.2
)

// AdService selects winning ads for placement
type AdService struct {
	lineItemService *LineItemService
	lineItemRepo    repo.LineItemRepository
	log             *zap.SugaredLogger
}

func NewAdService(lineItemRepo repo.LineItemRepository, lineItemService *LineItemService, log *zap.SugaredLogger) *AdService {
	return &AdService{
		lineItemService: lineItemService,
		lineItemRepo:    lineItemRepo,
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
	lineItems = s.winningAdCalculator(q, lineItems)

	var result []*model.Ad
	for _, lineItem := range lineItems {
		if lineItem.Budget < lineItem.Bid {
			continue
		}
		updatedLineItem, err := s.lineItemRepo.UpdateBudget(lineItem, lineItem.Budget-lineItem.Bid)
		if err != nil || updatedLineItem == nil {
			if !errors.Is(err, domain_errors.ErrLineItemAlreadyUpdated) {
				s.log.Warnw("line item is already updated before budget spending",
					"id", lineItem.ID)
			} else {
				s.log.Errorw("error in updating line item budget spending",
					"id", lineItem.ID)
			}
			continue
		}
		result = append(result, &model.Ad{
			ID:           updatedLineItem.ID,
			Name:         updatedLineItem.Name,
			AdvertiserID: updatedLineItem.AdvertiserID,
			Bid:          updatedLineItem.Bid,
			Placement:    updatedLineItem.Placement,
			ServeURL:     serveUrlGenerator(updatedLineItem),
		})
		if len(result) >= q.Limit {
			break
		}
	}

	var ids []string
	for _, ad := range result {
		ids = append(ids, ad.ID)
	}
	s.log.Infow("winning ads selected",
		"placement", q.Placement,
		"category", q.Category,
		"keyword", q.Keyword,
		"limit", q.Limit,
		"candidates", len(lineItems),
		"returned", len(result),
		"ad_ids", ids,
	)
	return result, nil
}

func (s *AdService) winningAdCalculator(q AdQuery, lineItems []*model.LineItem) []*model.LineItem {
	if len(lineItems) == 0 {
		return []*model.LineItem{}
	}

	scoredItems := make([]scoredItem, len(lineItems))
	for i, lineItem := range lineItems {
		scoredItems[i] = scoredItem{
			LineItem: lineItem,
			score:    lineItem.Bid * relevancyScore(lineItem, q),
		}
	}

	// Sort line items by descending score
	sort.Slice(lineItems, func(i, j int) bool {
		return scoredItems[i].score > scoredItems[j].score
	})
	return lineItems
}

func serveUrlGenerator(li *model.LineItem) string {
	return "/ad/serve/" + li.ID
}

func relevancyScore(lineItem *model.LineItem, q AdQuery) float64 {
	score := 1.0
	if q.Category != "" && slices.Contains(lineItem.Categories, q.Category) {
		score += wCategory
	}
	if q.Keyword != "" && slices.Contains(lineItem.Keywords, q.Keyword) {
		score += wKeyword
	}
	return score
}
