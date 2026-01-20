package workforce

import (
	"errors"
	"fmt"
	"workforce-ai-transition-simulator/internal/types"
)

// WorkforceManager manages the collection of human workers and AI agents
type WorkforceManager struct {
	humans         map[string]*types.HumanWorker
	aiAgents       map[string]*types.AIAgent
	businessOwnerID string
	nextHumanID    int
	nextAgentID    int
}

// NewWorkforceManager creates a new WorkforceManager instance
func NewWorkforceManager() *WorkforceManager {
	return &WorkforceManager{
		humans:      make(map[string]*types.HumanWorker),
		aiAgents:    make(map[string]*types.AIAgent),
		nextHumanID: 1,
		nextAgentID: 1,
	}
}

// GetHuman returns a human worker by ID
func (wm *WorkforceManager) GetHuman(id string) (*types.HumanWorker, bool) {
	human, exists := wm.humans[id]
	return human, exists
}

// GetAIAgent returns an AI agent by ID
func (wm *WorkforceManager) GetAIAgent(id string) (*types.AIAgent, bool) {
	agent, exists := wm.aiAgents[id]
	return agent, exists
}

// GetAllHumans returns all human workers
func (wm *WorkforceManager) GetAllHumans() []*types.HumanWorker {
	humans := make([]*types.HumanWorker, 0, len(wm.humans))
	for _, human := range wm.humans {
		humans = append(humans, human)
	}
	return humans
}

// GetAllAIAgents returns all AI agents
func (wm *WorkforceManager) GetAllAIAgents() []*types.AIAgent {
	agents := make([]*types.AIAgent, 0, len(wm.aiAgents))
	for _, agent := range wm.aiAgents {
		agents = append(agents, agent)
	}
	return agents
}

// GetBusinessOwner returns the business owner human worker
func (wm *WorkforceManager) GetBusinessOwner() (*types.HumanWorker, error) {
	if wm.businessOwnerID == "" {
		return nil, errors.New("no business owner exists")
	}
	human, exists := wm.humans[wm.businessOwnerID]
	if !exists {
		return nil, errors.New("business owner not found")
	}
	return human, nil
}

// AddHuman creates and adds a human worker with specified attributes
// If isBusinessOwner is true and no business owner exists, this worker becomes the business owner
// Returns the created human worker or an error
func (wm *WorkforceManager) AddHuman(experienceLevel types.ExperienceLevel, costCategory types.CostCategory, isBusinessOwner bool) (*types.HumanWorker, error) {
	// Generate unique ID
	id := fmt.Sprintf("human-%d", wm.nextHumanID)
	wm.nextHumanID++
	
	// If this is marked as business owner, ensure we don't already have one
	if isBusinessOwner && wm.businessOwnerID != "" {
		return nil, errors.New("business owner already exists")
	}
	
	// If no business owner exists yet, make this the business owner
	if wm.businessOwnerID == "" {
		isBusinessOwner = true
	}
	
	// Create the human worker
	human := types.NewHumanWorker(id, experienceLevel, costCategory, isBusinessOwner)
	
	// Add to collection
	wm.humans[id] = human
	
	// Track business owner
	if isBusinessOwner {
		wm.businessOwnerID = id
	}
	
	return human, nil
}

// RemoveHuman removes a human worker and releases all their assigned AI agents
// Prevents removal of the business owner
// Returns an error if the worker is the business owner or doesn't exist
func (wm *WorkforceManager) RemoveHuman(workerID string) error {
	// Check if worker exists
	human, exists := wm.humans[workerID]
	if !exists {
		return fmt.Errorf("human worker %s not found", workerID)
	}
	
	// Prevent removal of business owner
	if human.IsBusinessOwner {
		return errors.New("cannot remove business owner")
	}
	
	// Release all assigned AI agents
	for _, agentID := range human.AssignedAgents {
		// Remove the agent from the collection
		delete(wm.aiAgents, agentID)
	}
	
	// Remove the human worker
	delete(wm.humans, workerID)
	
	return nil
}

// AddAIAgent creates and assigns an AI agent to a human with available capacity
// Returns the created AI agent or an error if no capacity is available
func (wm *WorkforceManager) AddAIAgent(orchestratorID string, creationTime int) (*types.AIAgent, error) {
	// Check if orchestrator exists
	human, exists := wm.humans[orchestratorID]
	if !exists {
		return nil, fmt.Errorf("orchestrator %s not found", orchestratorID)
	}
	
	// Check if orchestrator has capacity
	if !human.CanOrchestrateMoreAgents() {
		return nil, fmt.Errorf("orchestrator %s has reached orchestration limit", orchestratorID)
	}
	
	// Generate unique ID
	id := fmt.Sprintf("agent-%d", wm.nextAgentID)
	wm.nextAgentID++
	
	// Create the AI agent
	agent := types.NewAIAgent(id, orchestratorID, creationTime)
	
	// Add to collection
	wm.aiAgents[id] = agent
	
	// Assign to orchestrator
	human.AssignedAgents = append(human.AssignedAgents, id)
	
	return agent, nil
}

// ReleaseAIAgent removes an AI agent instantaneously
// Returns an error if the agent doesn't exist
func (wm *WorkforceManager) ReleaseAIAgent(agentID string) error {
	// Check if agent exists
	agent, exists := wm.aiAgents[agentID]
	if !exists {
		return fmt.Errorf("AI agent %s not found", agentID)
	}
	
	// Remove agent from orchestrator's assigned list
	orchestrator, exists := wm.humans[agent.OrchestratorID]
	if exists {
		// Find and remove the agent ID from the orchestrator's list
		for i, id := range orchestrator.AssignedAgents {
			if id == agentID {
				// Remove by swapping with last element and truncating
				orchestrator.AssignedAgents[i] = orchestrator.AssignedAgents[len(orchestrator.AssignedAgents)-1]
				orchestrator.AssignedAgents = orchestrator.AssignedAgents[:len(orchestrator.AssignedAgents)-1]
				break
			}
		}
	}
	
	// Remove the agent from the collection
	delete(wm.aiAgents, agentID)
	
	return nil
}

// GetAvailableOrchestrationCapacity calculates the total available capacity across all humans
// Returns the sum of available capacity from all human workers
func (wm *WorkforceManager) GetAvailableOrchestrationCapacity() int {
	totalCapacity := 0
	for _, human := range wm.humans {
		totalCapacity += human.GetOrchestrationCapacity()
	}
	return totalCapacity
}

// CalculateTotalProductivity sums productivity from all humans and AI agents
// timeZoneInefficiency is the productivity penalty for Low_Cost_Non_US workers (0-1)
func (wm *WorkforceManager) CalculateTotalProductivity(timeZoneInefficiency float64) float64 {
	totalProductivity := 0.0
	
	// Sum human productivity
	for _, human := range wm.humans {
		totalProductivity += human.GetEffectiveProductivity(timeZoneInefficiency)
	}
	
	// Sum AI agent productivity
	for _, agent := range wm.aiAgents {
		totalProductivity += agent.GetProductivity()
	}
	
	return totalProductivity
}

// GetWorkforceComposition returns detailed workforce statistics
func (wm *WorkforceManager) GetWorkforceComposition() types.WorkforceComposition {
	composition := types.WorkforceComposition{}
	
	// Initialize maps
	composition.Humans.ByExperience = make(map[types.ExperienceLevel]int)
	composition.Humans.ByCostCategory = make(map[types.CostCategory]int)
	composition.AIAgents.ByExperience = make(map[types.ExperienceLevel]int)
	
	// Count humans
	composition.Humans.Total = len(wm.humans)
	for _, human := range wm.humans {
		composition.Humans.ByExperience[human.ExperienceLevel]++
		composition.Humans.ByCostCategory[human.CostCategory]++
	}
	
	// Count AI agents
	composition.AIAgents.Total = len(wm.aiAgents)
	for _, agent := range wm.aiAgents {
		composition.AIAgents.ByExperience[agent.ExperienceLevel]++
	}
	
	// Calculate orchestration utilization
	totalCapacity := len(wm.humans) * types.OrchestrationLimit
	if totalCapacity > 0 {
		usedCapacity := len(wm.aiAgents)
		composition.OrchestrationUtilization = (float64(usedCapacity) / float64(totalCapacity)) * 100.0
	} else {
		composition.OrchestrationUtilization = 0.0
	}
	
	return composition
}
