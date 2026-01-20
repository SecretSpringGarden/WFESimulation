package economic

import (
	"math"
	"workforce-ai-transition-simulator/internal/types"
)

// EconomicModel manages budget constraints and revenue calculations
type EconomicModel struct {
	fixedBudget     float64
	revenueScenario types.RevenueScenario
	revenueHistory  []float64
}

// NewEconomicModel creates a new EconomicModel instance
func NewEconomicModel(fixedBudget float64, revenueScenario types.RevenueScenario) *EconomicModel {
	return &EconomicModel{
		fixedBudget:     fixedBudget,
		revenueScenario: revenueScenario,
		revenueHistory:  make([]float64, 0),
	}
}

// GetFixedBudget returns the fixed budget value
func (em *EconomicModel) GetFixedBudget() float64 {
	return em.fixedBudget
}

// GetRevenueHistory returns the revenue history
func (em *EconomicModel) GetRevenueHistory() []float64 {
	return em.revenueHistory
}

// CalculateWorkforceCost sums costs of all humans and AI agents
func (em *EconomicModel) CalculateWorkforceCost(humans []*types.HumanWorker, agents []*types.AIAgent) float64 {
	totalCost := 0.0
	
	// Sum human costs
	for _, human := range humans {
		totalCost += human.BaseCost
	}
	
	// Sum AI agent costs
	for _, agent := range agents {
		totalCost += agent.GetCost()
	}
	
	return totalCost
}

// GetAvailableBudget calculates remaining budget after current workforce costs
func (em *EconomicModel) GetAvailableBudget(humans []*types.HumanWorker, agents []*types.AIAgent) float64 {
	currentCost := em.CalculateWorkforceCost(humans, agents)
	return em.fixedBudget - currentCost
}

// CanAfford checks if a cost fits within available budget
func (em *EconomicModel) CanAfford(cost float64, humans []*types.HumanWorker, agents []*types.AIAgent) bool {
	availableBudget := em.GetAvailableBudget(humans, agents)
	return cost <= availableBudget
}

// CalculateRevenue calculates revenue based on productivity and time step
// Handles Flat_Revenue and Explosive_Growth scenarios
func (em *EconomicModel) CalculateRevenue(productivity float64, timeStep int) float64 {
	var revenue float64
	
	switch em.revenueScenario {
	case types.FlatRevenue:
		// Flat revenue: constant multiplier of productivity
		revenue = productivity * 100000.0 // Base revenue multiplier
		
	case types.ExplosiveGrowth:
		// Explosive growth: exponential increase over time
		// Revenue = productivity * base_multiplier * (1 + growth_rate)^timeStep
		baseMultiplier := 100000.0
		growthRate := 0.05 // 5% growth per time step
		revenue = productivity * baseMultiplier * math.Pow(1.0+growthRate, float64(timeStep))
		
	default:
		// Default to flat revenue
		revenue = productivity * 100000.0
	}
	
	// Record revenue in history
	em.revenueHistory = append(em.revenueHistory, revenue)
	
	return revenue
}

// GetCostPerProductivityUnit calculates cost-effectiveness metric for workers
// Returns the cost per unit of productivity
func (em *EconomicModel) GetCostPerProductivityUnit(cost float64, productivity float64) float64 {
	if productivity == 0.0 {
		// Avoid division by zero - return a very high cost to indicate inefficiency
		return math.Inf(1)
	}
	return cost / productivity
}
