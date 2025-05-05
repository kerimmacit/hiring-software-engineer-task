package service

import (
	"time"

	"sweng-task/internal/domain_errors"
	"sweng-task/internal/model"
	"sweng-task/internal/repo"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// LineItemService provides operations for line items
type LineItemService struct {
	repo repo.LineItemRepository
	log  *zap.SugaredLogger
}

// NewLineItemService creates a new LineItemService
func NewLineItemService(repo repo.LineItemRepository, log *zap.SugaredLogger) *LineItemService {
	return &LineItemService{
		repo: repo,
		log:  log,
	}
}

// Create creates a new line item
func (s *LineItemService) Create(item model.LineItemCreate) (*model.LineItem, error) {
	now := time.Now()

	lineItem := &model.LineItem{
		ID:           "li_" + uuid.New().String(),
		Name:         item.Name,
		AdvertiserID: item.AdvertiserID,
		Bid:          item.Bid,
		Budget:       item.Budget,
		Placement:    item.Placement,
		Categories:   item.Categories,
		Keywords:     item.Keywords,
		Status:       model.LineItemStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err := s.repo.CreateLineItem(lineItem)
	if err != nil {
		return nil, err
	}
	s.log.Infow("Line item created",
		"id", lineItem.ID,
		"name", lineItem.Name,
		"advertiser_id", lineItem.AdvertiserID,
		"placement", lineItem.Placement,
	)
	return lineItem, nil
}

// GetByID retrieves a line item by ID
func (s *LineItemService) GetByID(id string) (*model.LineItem, error) {
	lineItem, err := s.repo.GetLineItemById(id)
	if err != nil {
		return nil, err
	}
	if lineItem == nil {
		return nil, domain_errors.ErrLineItemNotFound
	}
	return lineItem, nil
}

// GetAll retrieves all line items, optionally filtered by advertiser ID and placement
func (s *LineItemService) GetAll(advertiserID, placement string) ([]*model.LineItem, error) {
	return s.repo.GetLineItems(repo.GetLineItemsFilter{AdvertiserID: advertiserID, Placement: placement})
}

// FindMatchingLineItems finds line items matching the given placement and filters
// This method will be used by the AdService when implementing the ad selection logic
func (s *LineItemService) FindMatchingLineItems(placement string, category, keyword string) ([]*model.LineItem, error) {
	lineItems, err := s.repo.GetLineItems(repo.GetLineItemsFilter{
		Placement: placement,
		Status:    model.LineItemStatusActive,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*model.LineItem, 0)
	for _, item := range lineItems {
		// Apply category filter if specified
		if category != "" {
			categoryFound := false
			for _, cat := range item.Categories {
				if cat == category {
					categoryFound = true
					break
				}
			}
			if !categoryFound {
				continue
			}
		}

		// Apply keyword filter if specified
		if keyword != "" {
			keywordFound := false
			for _, kw := range item.Keywords {
				if kw == keyword {
					keywordFound = true
					break
				}
			}
			if !keywordFound {
				continue
			}
		}

		result = append(result, item)
	}

	return result, nil
}
