// services/payment_service.go
package services

import (
	"backend-go-gin/models"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type PaymentService struct {
	snapClient snap.Client
	coreClient coreapi.Client
}

func NewPaymentService() *PaymentService {
	// Setup Midtrans config
	var s = PaymentService{}
	s.snapClient.New("SB-Mid-server-49nmsGZlgAXoJQkYrI9bRTlP", midtrans.Sandbox)
	s.coreClient.New("SB-Mid-server-49nmsGZlgAXoJQkYrI9bRTlP", midtrans.Sandbox)
	return &s
}

func (ps *PaymentService) CreateSnapTransaction(order *models.Order, customer *models.User, detail *models.UserDetail) (*snap.Response, error) {
	// Buat request Snap
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  order.Invoice,
			GrossAmt: int64(order.TotalHarga),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: detail.Nama,
			Email: customer.Email,
			Phone: detail.Telepon,
		},
		EnabledPayments: snap.AllSnapPaymentType,
	}

	// Buat transaksi Snap
	snapResp, err := ps.snapClient.CreateTransaction(req)
	if err != nil {
		return nil, err
	}

	return snapResp, nil
}

func (ps *PaymentService) VerifyPayment(orderID string) (*coreapi.TransactionStatusResponse, error) {
	// Verifikasi status pembayaran
	transactionStatusResp, err := ps.coreClient.CheckTransaction(orderID)
	if err != nil {
		return nil, err
	}
	
	return transactionStatusResp, nil
}