package crm

type Customer struct {
	ID   string `json:"id"`
	Code string `json:"code"`
}

type FileMakerCustomer struct {
	ID           string `json:"Id_customer"`
	CustomerCode string `json:"CustomerCode"`
}

func (fc FileMakerCustomer) ToCustomer() Customer {
	return Customer{
		ID:   fc.ID,
		Code: fc.CustomerCode,
	}
}
