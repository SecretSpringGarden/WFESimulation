package economic

import (
	"math"
	"testing"
	"workforce-ai-transition-simulator/internal/types"
)

func TestGetCostPerProductivityUnit(t *testing.T) {
	em := NewEconomicModel(1000000.0, types.FlatRevenue)

	tests := []struct {
		name         string
		cost         float64
		productivity float64
		expected     float64
	}{
		{
			name:         "Normal case",
			cost:         100000.0,
			productivity: 50.0,
			expected:     2000.0,
		},
		{
			name:         "High productivity",
			cost:         100000.0,
			productivity: 100.0,
			expected:     1000.0,
		},
		{
			name:         "Low productivity",
			cost:         100000.0,
			productivity: 10.0,
			expected:     10000.0,
		},
		{
			name:         "Zero productivity",
			cost:         100000.0,
			productivity: 0.0,
			expected:     math.Inf(1),
		},
		{
			name:         "Zero cost",
			cost:         0.0,
			productivity: 50.0,
			expected:     0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := em.GetCostPerProductivityUnit(tt.cost, tt.productivity)
			if math.IsInf(tt.expected, 1) {
				if !math.IsInf(result, 1) {
					t.Errorf("Expected positive infinity, got %f", result)
				}
			} else if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}
