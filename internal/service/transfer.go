package service

import (
	"context"
	"errors"

	"bank-api/internal/repository"
)

// TransferService отвечает за переводы между счетами
type TransferService struct {
	accountRepo repository.AccountRepository
}

// NewTransferService создаёт новый экземпляр TransferService
func NewTransferService(accountRepo repository.AccountRepository) *TransferService {
	return &TransferService{
		accountRepo: accountRepo,
	}
}

// Transfer выполняет перевод средств между счетами с проверкой принадлежности и баланса
func (s *TransferService) Transfer(
	ctx context.Context,
	fromUserID int, // ID пользователя, выполняющего перевод
	fromAccID int, // ID счета отправителя
	toAccID int, // ID счета получателя
	amount float64, // сумма перевода
) error {
	// Проверяем, что счет отправителя существует и принадлежит данному пользователю
	from, err := s.accountRepo.GetByID(ctx, fromAccID)
	if err != nil {
		return errors.New("source account not found")
	}
	if from.UserID != fromUserID {
		return errors.New("access denied: not your account")
	}

	// Проверяем, что счет получателя существует (любой пользователь)
	_, err = s.accountRepo.GetByID(ctx, toAccID)
	if err != nil {
		return errors.New("destination account not found")
	}

	// Проверяем достаточность средств
	if from.Balance < amount {
		return errors.New("insufficient funds")
	}

	// Выполняем перевод (в репозитории уже wrapped in transaction)
	return s.accountRepo.Transfer(ctx, fromAccID, toAccID, amount)
}
