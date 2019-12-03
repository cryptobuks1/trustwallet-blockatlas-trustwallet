package blockatlas

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTicker_ApplyRate(t *testing.T) {
	tests := []struct {
		name        string
		price       float64
		percent     float64
		rate        float64
		wantPrice   float64
		wantPercent float64
	}{
		{"apply rate 1", 443, 0.3, 344, 152392, -45717.299999999996},
		{"apply rate 2", 1111.22, 1.2, 0.88, 977.8736, -1172.24832},
		{"apply rate 3", 3.33, -1, 22.3, 74.259, 73.259},
		{"apply rate 4", 9.3332, -0.88, 22, 205.3304, 179.810752},
		{"apply rate 5", 0.00000001, 4, 0.2, 0.000000002, 3.999999992},
		{"apply rate 6", 0.00000001, -2, 10, 0.0000001, -1.9999998},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			ticker := &Ticker{
				Price: TickerPrice{Value: tt.price, Change24h: tt.percent},
			}
			ticker.ApplyRate(tt.rate, DefaultCurrency)
			assert.Equal(t, tt.wantPrice, ticker.Price.Value)
			assert.Equal(t, tt.wantPercent, ticker.Price.Change24h)
		})
	}
}
