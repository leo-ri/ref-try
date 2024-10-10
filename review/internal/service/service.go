package service

import (
	. "coupon_service/internal/service/entity"
	"fmt"

	"github.com/google/uuid"
)

type CouponRepository interface {
	FindByCode(string) (*Coupon, error)
	Save(Coupon) error
}

// ---------

type Service struct {
	repo CouponRepository
}

func New(repo CouponRepository) Service {
	return Service{
		repo: repo,
	}
}

func (s Service) ApplyCoupon(totalValue int, code string) (b *Basket, e error) {
	b = &Basket{
		Value:                 totalValue,
		AppliedDiscount:       0,
		ApplicationSuccessful: false,
	}
	coupon, err := s.repo.FindByCode(code)
	if err != nil {
		return nil, err
	}
	if b.Value > coupon.MinBasketValue {
		b.AppliedDiscount = coupon.Discount
		b.ApplicationSuccessful = true
		// b.Value = b.Value * b.AppliedDiscount / 100 // total discount? or
		b.Value = totalValue - b.Value*b.AppliedDiscount/100 // total discount? or
	}
	if b.Value < 0 {
		return nil, fmt.Errorf("Tried to apply discount to negative value")
	}
	return
}

func (s Service) CreateCoupon(discount int, code string, minBasketValue int) any {
	coupon := Coupon{
		Discount:       discount,
		Code:           code,
		MinBasketValue: minBasketValue,
		ID:             uuid.NewString(),
	}

	if err := s.repo.Save(coupon); err != nil {
		return err
	}
	return nil
}

// ?? is this should be strict ? or we need to show all
// lets assume show that we found
func (s Service) GetCoupons(codes []string) ([]Coupon, error) {
	coupons := make([]Coupon, 0, len(codes))
	var e error = nil

	for idx, code := range codes {
		coupon, err := s.repo.FindByCode(code)
		if err != nil {
			if e == nil {
				e = fmt.Errorf("code: %s, index: %d", code, idx)
			} else {
				e = fmt.Errorf("%w; code: %s, index: %d", e, code, idx)
			}
			continue
		}
		if coupon != nil {
			coupons = append(coupons, *coupon)
		}
	}
	return coupons, e
}
