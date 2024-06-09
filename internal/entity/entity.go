package entity

type (
	Order struct {
		Uid         string
		TrackNumber string

		Entry             string
		Locale            string
		InternalSignature string

		CustomerID      string
		DeliveryService string

		DateCreated string

		ShardKey int
		SmID     int
		OffShard int

		Items []Item
		Payment
		Delivery
	}

	Item struct {
		OrderUID string

		TrackNumber string

		Status int

		ChrtID int
		NmID   int
		RID    int

		Brand string
		Name  string
		Size  string

		Price      float64
		Sale       float64
		TotalPrice float64
	}

	Payment struct {
		OrderUID string

		TxID      string
		RequestID string

		Bank     string
		Currency string
		Provider string

		PaymentDt string

		Amount       float64
		DeliveryCost float64
		GoodsTotal   float64
		CustomFee    float64
	}

	Delivery struct {
		OrderUID string

		Name  string
		Phome string
		Email string

		ZIP     string
		City    string
		Address string
		Region  string
	}
)
