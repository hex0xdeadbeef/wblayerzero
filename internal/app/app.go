package app

import (
	"log"
	"wblayerzero/internal/config"
	"wblayerzero/internal/database/postgres"
)

func Run() error {
	err := config.Load(config.CfgFilePath)
	if err != nil {
		return err
	}

	storage, err := postgres.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := storage.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}

/*
// generateRandomString создает случайную строку указанной длины
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// generateOrders генерирует слайс объектов Order с уникальными UID
func generateOrders(count int) []entities.Order {
	orders := make([]entities.Order, count)
	for i := 0; i < count; i++ {
		uid := generateRandomString(8)
		orders[i] = entities.Order{
			UID:               uid,
			TrackNumber:       fmt.Sprintf("track%d", i+1),
			Entry:             fmt.Sprintf("entry%d", i+1),
			Locale:            "en",
			InternalSignature: fmt.Sprintf("signature%d", i+1),
			CustomerID:        fmt.Sprintf("customer%d", i+1),
			DeliveryService:   fmt.Sprintf("delivery%d", i+1),
			DateCreated:       time.Now().Format(time.RFC3339),
			ShardKey:          i + 1,
			SmID:              i + 2,
			OffShard:          i + 3,
			Items: []entities.Item{
				{
					OrderUID:    uid,
					TrackNumber: fmt.Sprintf("track%d", i+1),
					Status:      1,
					ChrtID:      rand.Intn(100000),
					NmID:        rand.Intn(100000),
					RID:         rand.Intn(100000),
					Brand:       fmt.Sprintf("brand%d", i+1),
					Name:        fmt.Sprintf("item%d", i+1),
					Size:        "M",
					Price:       float64(rand.Intn(10000)) / 100,
					Sale:        float64(rand.Intn(1000)) / 100,
					TotalPrice:  float64(rand.Intn(10000)) / 100,
				},
			},
			Payment: entities.Payment{
				OrderUID:     uid,
				TxID:         fmt.Sprintf("tx%d", i+1),
				RequestID:    fmt.Sprintf("req%d", i+1),
				Bank:         fmt.Sprintf("bank%d", i+1),
				Currency:     "USD",
				Provider:     fmt.Sprintf("provider%d", i+1),
				PaymentDt:    time.Now().Format(time.RFC3339),
				Amount:       float64(rand.Intn(10000)) / 100,
				DeliveryCost: float64(rand.Intn(1000)) / 100,
				GoodsTotal:   float64(rand.Intn(10000)) / 100,
				CustomFee:    float64(rand.Intn(1000)) / 100,
			},
			Delivery: entities.Delivery{
				OrderUID: uid,
				Name:     fmt.Sprintf("John Doe%d", i+1),
				Phone:    fmt.Sprintf("+12345678%03d", i+1),
				Email:    fmt.Sprintf("johndoe%d@example.com", i+1),
				Zip:      fmt.Sprintf("%06d", rand.Intn(100000)),
				City:     fmt.Sprintf("City%d", i+1),
				Address:  fmt.Sprintf("123 Main St Apt %d", i+1),
				Region:   "NY",
			},
		}
	}
	return orders
}
**/
