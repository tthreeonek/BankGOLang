package service

import (
	"context"
	"errors"

	"bank-api/internal/model"
	"bank-api/internal/repository"
)

type AccountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) CreateForUser(ctx context.Context, userID int) (*model.Account, error) {
	acc := &model.Account{UserID: userID, Balance: 0}
	if err := s.repo.Create(ctx, acc); err != nil {
		return nil, err
	}
	return acc, nil
}

func (s *AccountService) GetUserAccounts(ctx context.Context, userID int) ([]model.Account, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *AccountService) Deposit(ctx context.Context, accountID, userID int, amount float64) error {
	acc, err := s.repo.GetByID(ctx, accountID)
	if err != nil {
		return err
	}
	if acc.UserID != userID {
		return errors.New("not your account")
	}
	return s.repo.UpdateBalance(ctx, accountID, amount)
}

func (s *AccountService) PredictBalance(ctx context.Context, userID, accountID, days int) (map[string]float64, error) {
	acc, err := s.repo.GetByID(ctx, accountID)
	if err != nil || acc.UserID != userID {
		return nil, errors.New("account not found or access denied")
	}
	// Упрощённый прогноз: вычитаем виртуальные будущие расходы (можно развить)
	predicted := acc.Balance - 500.0*float64(days)/30.0
	return map[string]float64{
		"current_balance":   acc.Balance,
		"predicted_balance": predicted,
		"days":              float64(days),
	}, nil
}
