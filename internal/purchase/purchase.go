package purchase

import (
	coffeeco "coffeeco/internal"
	"coffeeco/internal/loyalty"
	"coffeeco/internal/payment"
	"coffeeco/internal/store"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
)

// ------------------- PURCHASE AGGREGATE ----------------------
//
// -------------------------------------------------------------
type Purchase struct {
	id                 uuid.UUID
	Store              store.Store
	ProductsToPurchase []coffeeco.Product
	total              money.Money
	PaymentMeans       payment.Means
	timeofPurchase     time.Time
	CardToken          *string
}

func (p *Purchase) validateAndEnrich() error {
	if len(p.ProductsToPurchase) == 0 {
		return errors.New("purchase must consist of at least one product")
	}
	p.total = *money.New(0, "USD")

	for _, v := range p.ProductsToPurchase {
		newTotal, _ := p.total.Add(&v.BasePrice)
		p.total = *newTotal
	}
	if p.total.IsZero() {
		return errors.New("likely mistake; purchase should never  be 0. Please validate")
	}
	p.id = uuid.New()
	p.timeofPurchase = time.Now()

	return nil
}

// ---------------- CARD CHARGE SERVICE INTERFACE ----------------------
//
// ---------------------------------------------------------------------

type CardChargeService interface {
	ChargeCard(ctx context.Context, amount money.Money, cardToken string) error
}

// ---------------- StoreService SERVICE INTERFACE ----------------------
//
// ---------------------------------------------------------------------
type StoreService interface {
	GetStoreSpecificDiscount(ctx context.Context, storeID uuid.UUID) (float32, error)
}

// ---------------------------- SERVICE -------------------------------
//
// --------------------------------------------------------------------

type Service struct {
	cardService  CardChargeService
	purchaseRepo Repository
	storeService StoreService
}

func NewService(cardService CardChargeService, purchaseRepo Repository, storeService StoreService) *Service {
	return &Service{cardService: cardService, purchaseRepo: purchaseRepo, storeService: storeService}
}

func (s Service) CompletePurchase(ctx context.Context, storeID uuid.UUID, purchase *Purchase, coffeeBuxCard *loyalty.Coffeebux) error {
	if err := purchase.validateAndEnrich(); err != nil {
		return err
	}

	if err := s.calculateStoreSpecificDiscount(ctx, storeID, purchase); err != nil {
		return err
	}

	switch purchase.PaymentMeans {
	case payment.MEANS_CARD:
		if err := s.cardService.ChargeCard(ctx, purchase.total, *purchase.CardToken); err != nil {
			return errors.New("card charge failed,cancelling purchase")
		}
	case payment.MEANS_CASH:
	// For the reader to add :)
	case payment.MEANS_COFFEEBUX:
		if err := coffeeBuxCard.Pay(ctx, purchase.ProductsToPurchase); err != nil {
			return fmt.Errorf("failed to charge loyatlycard: %w", err)
		}
	default:
		return errors.New("unknown payment type")
	}

	if err := s.purchaseRepo.Store(ctx, *purchase); err != nil {
		return errors.New("failed to store purchase")
	}
	if coffeeBuxCard != nil {
		coffeeBuxCard.AddStamp()
	}
	return nil
}

func (s *Service) calculateStoreSpecificDiscount(ctx context.Context, storeID uuid.UUID, purchase *Purchase) error {
	discount, err := s.storeService.GetStoreSpecificDiscount(ctx, storeID)

	if err != nil && err != store.ErrNoDiscount {
		return fmt.Errorf("failed to get discount: %w", err)
	}
	purchasePrice := purchase.total
	if discount > 0 {
		purchase.total = *purchasePrice.Multiply(int64(100 - discount))
	}
	return nil
}
