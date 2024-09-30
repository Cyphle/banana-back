package account

import (
	"banana-back/domain"
	"time"
)

type Budget struct {
	ID            int64
	frequency     domain.Frequency
	initialAmount float64
	actualAmount  float64
	name          string
	startDate     time.Time
}
