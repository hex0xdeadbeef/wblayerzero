package entities

type (
	Order struct {
		UID         string `json:"order_uid"`
		TrackNumber string `json:"track_number"`

		Entry             string `json:"entry"`
		Locale            string `json:"locale"`
		InternalSignature string `json:"internal_signature"`

		CustomerID      string `json:"customer_id"`
		DeliveryService string `json:"delivery_service"`

		DateCreated string `json:"date_created"`

		ShardKey int `json:"shardkey"`
		SmID     int `json:"sm_id"`
		OffShard int `json:"oof_shard"`

		Items    []Item `json:"items"`
		Payment  `json:"payment"`
		Delivery `json:"delivery"`
	}

	Item struct {
		OrderUID string `json:"order_uid"`

		TrackNumber string `json:"track_number"`

		Status int `json:"status"`

		ChrtID int `json:"chrt_id"`
		NmID   int `json:"nm_id"`
		RID    int `json:"rid"`

		Brand string `json:"brand"`
		Name  string `json:"name"`
		Size  string `json:"size"`

		Price      float64 `json:"price"`
		Sale       float64 `json:"sale"`
		TotalPrice float64 `json:"total_price"`
	}

	Payment struct {
		OrderUID string `json:"order_uid"`

		TxID      string `json:"transaction_id"`
		RequestID string `json:"request_id"`

		Bank     string `json:"bank"`
		Currency string `json:"currency"`
		Provider string `json:"provider"`

		PaymentDt string `json:"payment_dt"`

		Amount       float64 `json:"amount"`
		DeliveryCost float64 `json:"delivery_cost"`
		GoodsTotal   float64 `json:"goods_total"`
		CustomFee    float64 `json:"custom_fee"`
	}

	Delivery struct {
		OrderUID string `json:"order_uid"`

		Zip   string `json:"zip"`
		Name  string `json:"name"`
		Phone string `json:"phone"`
		Email string `json:"email"`

		City    string `json:"city"`
		Address string `json:"address"`
		Region  string `json:"region"`
	}
)
