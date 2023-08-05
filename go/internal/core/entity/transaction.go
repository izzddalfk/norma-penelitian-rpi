package entity

type Transaction struct {
	ID            int64
	TotalAmount   float64
	PaymentAmount float64
	ReturnAmount  float64
}

func (t *Transaction) SetPaymentAndReturnAmount(payAmount float64) {
	t.PaymentAmount = payAmount
	t.ReturnAmount = payAmount - t.TotalAmount
}
