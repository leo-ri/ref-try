package api

import (
	. "coupon_service/internal/api/entity"
	serviceEntity "coupon_service/internal/service/entity"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Service interface {
	ApplyCoupon(int, string) (*serviceEntity.Basket, error)
	CreateCoupon(int, string, int) any
	GetCoupons([]string) ([]serviceEntity.Coupon, error)
}

// TODO response with message error (c.JSON)

func (a *API) Apply(c *gin.Context) {
	apiReq := ApplicationRequest{}
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	basket, err := a.svc.ApplyCoupon(apiReq.BasketValue, apiReq.Code)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, basket)
}

func (a *API) Create(c *gin.Context) {
	apiReq := Coupon{}
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	err := a.svc.CreateCoupon(apiReq.Discount, apiReq.Code, apiReq.MinBasketValue)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	_, _ = c.Writer.Write([]byte(apiReq.Code))
	log.Print(apiReq.Code)
	c.Status(http.StatusOK)
}

func (a *API) Get(c *gin.Context) {
	codes := c.QueryArray("code")
	if len(codes) < 1 {
		c.Status(http.StatusBadRequest)
		return
	}
	log.Print("Requested: ", codes)

	coupons, err := a.svc.GetCoupons(codes)
	if err != nil {
		log.Printf("Couldnot find thoose: %v", err)
	}
	c.JSON(http.StatusOK, coupons)
}
