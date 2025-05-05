package repo

import (
	"sync"

	"go.uber.org/zap"

	"sweng-task/internal/domain_errors"
	"sweng-task/internal/model"
)

type GetLineItemsFilter struct {
	Status       model.LineItemStatus
	AdvertiserID string
	Placement    string
}

type LineItemRepository interface {
	CreateLineItem(event *model.LineItem) error
	GetLineItemById(id string) (*model.LineItem, error)
	GetLineItems(filter GetLineItemsFilter) ([]*model.LineItem, error)
	UpdateBudget(li *model.LineItem, newBudget float64) (*model.LineItem, error)
}

var _ LineItemRepository = (*LineItemRepositoryImp)(nil)

type LineItemRepositoryImp struct {
	items map[string]*model.LineItem
	mu    sync.RWMutex
	log   *zap.SugaredLogger
}

func NewLineItemRepository(log *zap.SugaredLogger) LineItemRepository {
	return &LineItemRepositoryImp{
		items: make(map[string]*model.LineItem),
		log:   log,
	}
}

func (s *LineItemRepositoryImp) CreateLineItem(item *model.LineItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[item.ID] = item
	return nil
}

func (s *LineItemRepositoryImp) GetLineItemById(id string) (*model.LineItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.items[id], nil
}

func (s *LineItemRepositoryImp) GetLineItems(filter GetLineItemsFilter) ([]*model.LineItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*model.LineItem, 0)
	for _, item := range s.items {
		if filter.Status != "" && item.Status != filter.Status {
			continue
		}
		if filter.AdvertiserID != "" && item.AdvertiserID != filter.AdvertiserID {
			continue
		}
		if filter.Placement != "" && item.Placement != filter.Placement {
			continue
		}
		result = append(result, item)
	}
	return result, nil
}

func (s *LineItemRepositoryImp) UpdateBudget(li *model.LineItem, newBudget float64) (*model.LineItem, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.items[li.ID]
	if !ok || item.Budget != li.Budget || !item.UpdatedAt.Equal(li.UpdatedAt) {
		return nil, domain_errors.ErrLineItemAlreadyUpdated
	}
	item.Budget = newBudget
	s.items[item.ID] = item
	return item, nil
}
