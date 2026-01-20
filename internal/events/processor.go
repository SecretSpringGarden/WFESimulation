package events

import (
	"math/rand"
	"workforce-ai-transition-simulator/internal/types"
)

// EventProcessor handles attrition, learning, failures, and workforce optimization
type EventProcessor struct {
	attritionConfig         types.AttritionConfig
	catastrophicFailureRate float64
	aiLearningSpeed         types.AILearningSpeed
	timeZoneInefficiency    float64
	rng                     *rand.Rand
}

// NewEventProcessor creates a new EventProcessor instance
func NewEventProcessor(
	attritionConfig types.AttritionConfig,
	catastrophicFailureRate float64,
	aiLearningSpeed types.AILearningSpeed,
	timeZoneInefficiency float64,
	rng *rand.Rand,
) *EventProcessor {
	return &EventProcessor{
		attritionConfig:         attritionConfig,
		catastrophicFailureRate: catastrophicFailureRate,
		aiLearningSpeed:         aiLearningSpeed,
		timeZoneInefficiency:    timeZoneInefficiency,
		rng:                     rng,
	}
}


// ProcessAttrition handles different types of human worker attrition
// Returns a list of worker IDs to remove
func (ep *EventProcessor) ProcessAttrition(humans []*types.HumanWorker, timeStep int) []string {
	workersToRemove := make([]string, 0)
	
	switch ep.attritionConfig.Type {
	case types.NaturalAttrition:
		// Natural attrition: probabilistically remove workers based on natural rate
		// Convert annual rate to per-time-step probability
		// Assuming each time step represents a month (12 time steps per year)
		monthlyRate := ep.attritionConfig.NaturalRate / 12.0 / 100.0
		
		// Apply forced acceleration
		effectiveRate := monthlyRate * ep.attritionConfig.ForcedAcceleration
		
		for _, human := range humans {
			// Never remove business owner
			if human.IsBusinessOwner {
				continue
			}
			
			// Probabilistically determine if this worker leaves
			if ep.rng.Float64() < effectiveRate {
				workersToRemove = append(workersToRemove, human.ID)
			}
		}
		
	case types.HiringFreeze:
		// Hiring freeze: still allow natural attrition but prevent new hires
		// This is handled by the simulation controller, but we still process natural attrition
		monthlyRate := ep.attritionConfig.NaturalRate / 12.0 / 100.0
		effectiveRate := monthlyRate * ep.attritionConfig.ForcedAcceleration
		
		for _, human := range humans {
			if human.IsBusinessOwner {
				continue
			}
			
			if ep.rng.Float64() < effectiveRate {
				workersToRemove = append(workersToRemove, human.ID)
			}
		}
		
	case types.ReductionInForce:
		// Reduction in force: actively remove workers according to RIF parameters
		// Use forced acceleration as the percentage of workforce to remove
		targetRemovalCount := int(float64(len(humans)) * ep.attritionConfig.ForcedAcceleration / 100.0)
		
		// Select workers to remove (excluding business owner)
		eligibleWorkers := make([]*types.HumanWorker, 0)
		for _, human := range humans {
			if !human.IsBusinessOwner {
				eligibleWorkers = append(eligibleWorkers, human)
			}
		}
		
		// Randomly select workers to remove
		// Shuffle and take the first N workers
		ep.rng.Shuffle(len(eligibleWorkers), func(i, j int) {
			eligibleWorkers[i], eligibleWorkers[j] = eligibleWorkers[j], eligibleWorkers[i]
		})
		
		// Take up to targetRemovalCount workers
		removalCount := targetRemovalCount
		if removalCount > len(eligibleWorkers) {
			removalCount = len(eligibleWorkers)
		}
		
		for i := 0; i < removalCount; i++ {
			workersToRemove = append(workersToRemove, eligibleWorkers[i].ID)
		}
	}
	
	return workersToRemove
}


// ProcessLearning updates experience for all AI agents and triggers level-ups
func (ep *EventProcessor) ProcessLearning(agents []*types.AIAgent, timeDelta int) {
	// Data exposure is typically 1.0 (full exposure)
	dataExposure := 1.0
	
	for _, agent := range agents {
		// Accumulate experience based on time and data exposure
		agent.AccumulateExperience(timeDelta, dataExposure)
		
		// Check and trigger level-ups
		// An agent might level up multiple times if enough experience is accumulated
		for agent.CheckLevelUp(ep.aiLearningSpeed) {
			// Level up occurred, continue checking in case of multiple level-ups
		}
	}
}
