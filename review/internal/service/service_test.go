package service

import (
	"coupon_service/internal/repository/memdb"
	sEntity "coupon_service/internal/service/entity"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		repo CouponRepository
	}
	tests := []struct {
		name string
		args args
		want Service
	}{
		{"initialize service", args{repo: nil}, Service{repo: nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_ApplyCoupon(t *testing.T) {
	presetupRepo := memdb.New()
	coupon := &sEntity.Coupon{
		Discount:       10,
		Code:           "code10",
		MinBasketValue: 100,
	}
	presetupRepo.Save(*coupon)

	type args struct {
		value int
		code  string
	}
	tests := []struct {
		name    string
		args    args
		wantB   *sEntity.Basket
		wantErr bool
	}{
		{"apply 10% succes", args{1000, "code10"}, &sEntity.Basket{
			Value:                 900,
			AppliedDiscount:       10,
			ApplicationSuccessful: true,
		}, false},
		{"apply 10% fail", args{50, "code10"}, &sEntity.Basket{
			Value:                 50,
			AppliedDiscount:       0,
			ApplicationSuccessful: false,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				repo: presetupRepo,
			}
			gotB, err := s.ApplyCoupon(tt.args.value, tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyCoupon() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotB, tt.wantB) {
				t.Errorf("ApplyCoupon() gotB = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}

func TestService_CreateCoupon(t *testing.T) {
	type fields struct {
		repo *memdb.Repository
	}
	type args struct {
		discount       int
		code           string
		minBasketValue int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   any
	}{
		{"Apply 10%", fields{memdb.New()}, args{10, "Superdiscount", 55}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				repo: tt.fields.repo,
			}

			result := s.CreateCoupon(tt.args.discount, tt.args.code, tt.args.minBasketValue)
			if result != nil {
				t.Errorf("Should go without errors: %v", result)
			}
		})
	}
}
