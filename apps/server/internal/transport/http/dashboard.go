package httptransport

import (
	"context"
	"errors"
	"net/http"

	heatservices "github.com/Hughzu/trackstack/apps/server/internal/contexts/heat/application/services"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/calories"
	"github.com/Hughzu/trackstack/apps/server/internal/modules/expenses"
)

type HeatDashboardFunc func(context.Context, heatservices.GetDashboardRequest) (heatservices.DashboardViewModel, error)

type DashboardHandler struct {
	heat     HeatDashboardFunc
	calories *calories.Service
	expenses *expenses.Service
}

type OverarchingDashboardViewModel struct {
	Expenses expenses.ViewDashboard          `json:"expenses"`
	Calories calories.DashboardViewModel     `json:"calories"`
	Heat     heatservices.DashboardViewModel `json:"heat"`
}

func (h *DashboardHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireAuthUserID(w, r)
	if !ok {
		return
	}

	ctx := r.Context()

	expensesFuture, expensesErrChan := make(chan expenses.ViewDashboard), make(chan error)
	caloriesFuture, caloriesErrChan := make(chan calories.DashboardViewModel), make(chan error)
	heatFuture, heatErrChan := make(chan heatservices.DashboardViewModel), make(chan error)

	go func() {
		data, err := h.expenses.GetDashboard(ctx, expenses.GetCurrentSheetRequest{UserID: userID})
		if err != nil {
			expensesErrChan <- err
			return
		}
		expensesFuture <- data
	}()

	go func() {
		data, err := h.calories.GetDashboard(ctx, calories.GetDashboardRequest{UserID: userID, RecentLimit: 8})
		if err != nil {
			caloriesErrChan <- err
			return
		}
		caloriesFuture <- data
	}()

	go func() {
		data, err := h.heat(ctx, heatservices.GetDashboardRequest{UserID: userID, Page: 1, Limit: 20})
		if err != nil {
			heatErrChan <- err
			return
		}
		heatFuture <- data
	}()

	var expensesData expenses.ViewDashboard
	var caloriesData calories.DashboardViewModel
	var heatData heatservices.DashboardViewModel

	var errs []error

	select {
	case data := <-expensesFuture:
		expensesData = data
	case err := <-expensesErrChan:
		errs = append(errs, err)
	}

	select {
	case data := <-caloriesFuture:
		caloriesData = data
	case err := <-caloriesErrChan:
		errs = append(errs, err)
	}

	select {
	case data := <-heatFuture:
		heatData = data
	case err := <-heatErrChan:
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		writeError(w, errors.Join(errs...))
		return
	}

	writeJSON(w, http.StatusOK, OverarchingDashboardViewModel{
		Expenses: expensesData,
		Calories: caloriesData,
		Heat:     heatData,
	})
}
