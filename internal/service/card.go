package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"bank-api/internal/model"
	"bank-api/internal/repository"
	"bank-api/pkg/crypto"
	"bank-api/pkg/luhn"

	"golang.org/x/crypto/bcrypt"
)

type CardService struct {
	cardRepo    repository.CardRepository
	accountRepo repository.AccountRepository
	encrypt     func(string) (string, error) // PGP-шифрование
	hmacSecret  []byte
}

func NewCardService(
	cardRepo repository.CardRepository,
	accountRepo repository.AccountRepository,
	encrypt func(string) (string, error),
	hmacSecret string,
) *CardService {
	return &CardService{
		cardRepo:    cardRepo,
		accountRepo: accountRepo,
		encrypt:     encrypt,
		hmacSecret:  []byte(hmacSecret),
	}
}

func (s *CardService) CreateCard(ctx context.Context, userID, accountID int) (*model.Card, error) {
	acc, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil || acc.UserID != userID {
		return nil, errors.New("account not found or access denied")
	}

	pan := luhn.Generate("4")
	expDate := time.Now().AddDate(3, 0, 0).Format("01/06")
	cvv := fmt.Sprintf("%03d", secureRandomInt(100, 999))

	encryptedPAN, err := s.encrypt(pan)
	if err != nil {
		return nil, fmt.Errorf("encrypt PAN: %w", err)
	}
	hmac := crypto.ComputeHMAC(pan, s.hmacSecret)
	encryptedExp, err := s.encrypt(expDate)
	if err != nil {
		return nil, fmt.Errorf("encrypt expiration: %w", err)
	}
	cvvHash, err := bcrypt.GenerateFromPassword([]byte(cvv), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash CVV: %w", err)
	}

	card := &model.Card{
		AccountID:    accountID,
		EncryptedPAN: encryptedPAN,
		PANHMAC:      hmac,
		EncryptedExp: encryptedExp,
		CVVHash:      string(cvvHash),
		MaskedPAN:    pan[len(pan)-4:],
	}
	if err := s.cardRepo.Create(ctx, card); err != nil {
		return nil, fmt.Errorf("create card in repo: %w", err)
	}
	return card, nil
}

func secureRandomInt(min, max int) int {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	if err != nil {
		return min
	}
	return min + int(nBig.Int64())
}
