package application

import (
	"context"

	"github.com/Hughzu/trackstack/apps/server/internal/contexts/heat/application/ports"
	"github.com/Hughzu/trackstack/apps/server/internal/contexts/heat/application/services"
	"github.com/Hughzu/trackstack/apps/server/internal/contexts/heat/domain"
)

var ErrInvalidInput = domain.ErrInvalidInput

type Refill = domain.Refill

type ListRefillsRequest = services.ListRefillsRequest

type CreateRefillRequest = services.CreateRefillRequest

type DeleteRefillRequest = services.DeleteRefillRequest

type SeasonSnapshot = services.SeasonSnapshot

type DashboardViewModel = services.DashboardViewModel

type GetDashboardRequest = services.GetDashboardRequest

type RefillStore interface {
	ports.RefillRangeLister
	ports.RecentRefillLister
	ports.LatestRefillGetter
	ports.RefillRangeSummarizer
	ports.RefillCreator
	ports.RefillDeleter
}

type Service struct {
	listRefills  services.ListRefillsService
	createRefill services.CreateRefillService
	deleteRefill services.DeleteRefillService
	dashboard    services.GetDashboardService
}

func NewService(store RefillStore) *Service {
	return &Service{
		listRefills:  services.NewListRefillsService(store),
		createRefill: services.NewCreateRefillService(store),
		deleteRefill: services.NewDeleteRefillService(store),
		dashboard:    services.NewGetDashboardService(store, store, store),
	}
}

func (s *Service) ListRefills(ctx context.Context, req services.ListRefillsRequest) ([]domain.Refill, error) {
	return s.listRefills.Execute(ctx, req)
}

func (s *Service) CreateRefill(ctx context.Context, req services.CreateRefillRequest) (domain.Refill, error) {
	return s.createRefill.Execute(ctx, req)
}

func (s *Service) DeleteRefill(ctx context.Context, req services.DeleteRefillRequest) (bool, error) {
	return s.deleteRefill.Execute(ctx, req)
}

func (s *Service) GetDashboard(ctx context.Context, req services.GetDashboardRequest) (services.DashboardViewModel, error) {
	return s.dashboard.Execute(ctx, req)
}
