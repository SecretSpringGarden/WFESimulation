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


// CatastrophicFailure represents a critical system failure event
type CatastrophicFailure struct {
	TimeStep int
	Severity float64 // 0-1, where 1 is most severe
}

// GenerateCatastrophicFailure probabilistically generates failure events
// Returns a failure event or nil if no failure occurs
func (ep *EventProcessor) GenerateCatastrophicFailure(timeStep int) *CatastrophicFailure {
	// Check if a failure occurs based on the configured rate
	if ep.rng.Float64() < ep.catastrophicFailureRate {
		// Generate a failure with random severity
		severity := ep.rng.Float64()
		return &CatastrophicFailure{
			TimeStep: timeStep,
			Severity: severity,
		}
	}
	
	return nil
}


// FailureOutcome represents the result of evaluating a catastrophic failure
type FailureOutcome struct {
	CanHandle            bool
	ProductivityPenalty  float64 // 0-1, percentage reduction in productivity
	RequiresHumanIntervention bool
}

// EvaluateFailureResponse assesses workforce capability to handle failures
// Determines if productivity penalties should be applied
func (ep *EventProcessor) EvaluateFailureResponse(
	failure *CatastrophicFailure,
	humans []*types.HumanWorker,
	agents []*types.AIAgent,
) FailureOutcome {
	// Count senior+ humans (Senior and Executive)
	seniorHumanCount := 0
	for _, human := range humans {
		if human.ExperienceLevel >= types.Senior {
			seniorHumanCount++
		}
	}
	
	// Count senior+ AI agents
	seniorAgentCount := 0
	for _, agent := range agents {
		if agent.ExperienceLevel >= types.Senior {
			seniorAgentCount++
		}
	}
	
	// Calculate workforce capability score
	// Senior humans are more valuable for handling failures
	humanCapability := float64(seniorHumanCount) * 1.0
	agentCapability := float64(seniorAgentCount) * 0.5 // AI agents are less capable
	
	totalCapability := humanCapability + agentCapability
	
	// Determine if workforce can handle the failure
	// Require at least one senior human for any failure
	if seniorHumanCount == 0 {
		// No senior humans - cannot handle failure
		return FailureOutcome{
			CanHandle:                 false,
			ProductivityPenalty:       failure.Severity * 0.5, // 50% of severity as penalty
			RequiresHumanIntervention: true,
		}
	}
	
	// Check if capability is sufficient for the failure severity
	requiredCapability := failure.Severity * 3.0 // Scale severity to required capability
	
	if totalCapability >= requiredCapability {
		// Workforce can handle the failure
		return FailureOutcome{
			CanHandle:                 true,
			ProductivityPenalty:       0.0,
			RequiresHumanIntervention: false,
		}
	}
	
	// Workforce cannot fully handle the failure
	// Apply productivity penalty proportional to the capability gap
	capabilityGap := (requiredCapability - totalCapability) / requiredCapability
	penalty := failure.Severity * capabilityGap * 0.3 // Up to 30% penalty
	
	return FailureOutcome{
		CanHandle:                 false,
		ProductivityPenalty:       penalty,
		RequiresHumanIntervention: true,
	}
}


// WorkforceChange represents a proposed change to the workforce
type WorkforceChange struct {
	HireAIAgents     int      // Number of AI agents to hire
	ReleaseAIAgents  []string // IDs of AI agents to release
	OrchestratorID   string   // ID of human to assign new agents to
}

// OptimizeWorkforce evaluates hiring/release opportunities
// Prioritizes cost-effective decisions while respecting budget and orchestration constraints
func (ep *EventProcessor) OptimizeWorkforce(
	humans []*types.HumanWorker,
	agents []*types.AIAgent,
	availableBudget float64,
	availableOrchestrationCapacity int,
) WorkforceChange {
	change := WorkforceChange{
		HireAIAgents:    0,
		ReleaseAIAgents: make([]string, 0),
	}
	
	// If no orchestration capacity, we can't hire agents
	if availableOrchestrationCapacity <= 0 {
		return change
	}
	
	// Calculate cost-effectiveness of hiring a new AI agent
	// Start with University_Hire level agent
	newAgentCost := types.AIAgentCosts[types.UniversityHire]
	newAgentProductivity := types.AIAgentProductivity[types.UniversityHire]
	
	// Check if we can afford at least one agent
	if availableBudget < newAgentCost {
		return change
	}
	
	// Calculate cost per productivity unit for new agent
	newAgentCostPerProductivity := newAgentCost / newAgentProductivity
	
	// Find the most cost-effective human to compare against
	// (This helps decide if we should hire AI instead of humans)
	bestHumanCostPerProductivity := 0.0
	for _, human := range humans {
		effectiveProductivity := human.GetEffectiveProductivity(ep.timeZoneInefficiency)
		if effectiveProductivity > 0 {
			costPerProductivity := human.BaseCost / effectiveProductivity
			if bestHumanCostPerProductivity == 0 || costPerProductivity < bestHumanCostPerProductivity {
				bestHumanCostPerProductivity = costPerProductivity
			}
		}
	}
	
	// Hire AI agents if they are more cost-effective than humans
	// or if we have budget and capacity available
	if newAgentCostPerProductivity < bestHumanCostPerProductivity || bestHumanCostPerProductivity == 0 {
		// Calculate how many agents we can hire
		maxAgentsByBudget := int(availableBudget / newAgentCost)
		maxAgentsToHire := maxAgentsByBudget
		if maxAgentsToHire > availableOrchestrationCapacity {
			maxAgentsToHire = availableOrchestrationCapacity
		}
		
		// Find the best orchestrator (human with most available capacity)
		var bestOrchestrator *types.HumanWorker
		maxCapacity := 0
		for _, human := range humans {
			capacity := human.GetOrchestrationCapacity()
			if capacity > maxCapacity {
				maxCapacity = capacity
				bestOrchestrator = human
			}
		}
		
		if bestOrchestrator != nil && maxAgentsToHire > 0 {
			// Hire agents up to the orchestrator's capacity
			agentsToHire := maxAgentsToHire
			if agentsToHire > bestOrchestrator.GetOrchestrationCapacity() {
				agentsToHire = bestOrchestrator.GetOrchestrationCapacity()
			}
			
			change.HireAIAgents = agentsToHire
			change.OrchestratorID = bestOrchestrator.ID
		}
	}
	
	// Check if we should release any agents due to budget constraints
	// This would happen if we're over budget (shouldn't normally occur)
	// or if agents are not cost-effective
	if availableBudget < 0 {
		// We're over budget - need to release agents
		// Release the least productive agents first
		type agentScore struct {
			id           string
			productivity float64
		}
		
		agentScores := make([]agentScore, 0, len(agents))
		for _, agent := range agents {
			agentScores = append(agentScores, agentScore{
				id:           agent.ID,
				productivity: agent.GetProductivity(),
			})
		}
		
		// Sort by productivity (ascending - least productive first)
		for i := 0; i < len(agentScores)-1; i++ {
			for j := i + 1; j < len(agentScores); j++ {
				if agentScores[i].productivity > agentScores[j].productivity {
					agentScores[i], agentScores[j] = agentScores[j], agentScores[i]
				}
			}
		}
		
		// Release agents until we're back under budget
		budgetDeficit := -availableBudget
		for _, score := range agentScores {
			if budgetDeficit <= 0 {
				break
			}
			
			// Find the agent and get its cost
			for _, agent := range agents {
				if agent.ID == score.id {
					change.ReleaseAIAgents = append(change.ReleaseAIAgents, agent.ID)
					budgetDeficit -= agent.GetCost()
					break
				}
			}
		}
	}
	
	return change
}
