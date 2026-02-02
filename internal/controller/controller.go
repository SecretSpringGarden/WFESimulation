package controller

import (
	"errors"
	"fmt"
	"math/rand"
	"workforce-ai-transition-simulator/internal/economic"
	"workforce-ai-transition-simulator/internal/events"
	"workforce-ai-transition-simulator/internal/types"
	"workforce-ai-transition-simulator/internal/workforce"
)

// SimulationController coordinates WorkforceManager, EconomicModel, and EventProcessor
// and tracks simulation state throughout the execution
type SimulationController struct {
	config           types.SimulationConfig
	workforceManager *workforce.WorkforceManager
	economicModel    *economic.EconomicModel
	eventProcessor   *events.EventProcessor
	
	// Simulation state tracking
	currentTimeStep           int
	timeSeries               []types.SimulationState
	totalCatastrophicFailures int
	equilibriumReached        bool
	
	// Random number generator for reproducible results
	rng *rand.Rand
}

// NewSimulationController creates a new SimulationController instance
func NewSimulationController(config types.SimulationConfig, seed int64) *SimulationController {
	// Create random number generator with seed for reproducibility
	rng := rand.New(rand.NewSource(seed))
	
	// Create component instances
	workforceManager := workforce.NewWorkforceManager()
	economicModel := economic.NewEconomicModel(config.FixedBudget, config.RevenueScenario)
	eventProcessor := events.NewEventProcessor(
		config.AttritionConfig,
		config.CatastrophicFailureRate,
		config.AILearningSpeeds,
		config.TimeZoneInefficiency,
		rng,
	)
	
	return &SimulationController{
		config:                    config,
		workforceManager:         workforceManager,
		economicModel:            economicModel,
		eventProcessor:           eventProcessor,
		currentTimeStep:          0,
		timeSeries:               make([]types.SimulationState, 0),
		totalCatastrophicFailures: 0,
		equilibriumReached:       false,
		rng:                      rng,
	}
}

// GetConfig returns the simulation configuration
func (sc *SimulationController) GetConfig() types.SimulationConfig {
	return sc.config
}

// GetCurrentTimeStep returns the current simulation time step
func (sc *SimulationController) GetCurrentTimeStep() int {
	return sc.currentTimeStep
}

// GetTimeSeries returns the complete time series data
func (sc *SimulationController) GetTimeSeries() []types.SimulationState {
	return sc.timeSeries
}

// GetTotalCatastrophicFailures returns the total number of catastrophic failures encountered
func (sc *SimulationController) GetTotalCatastrophicFailures() int {
	return sc.totalCatastrophicFailures
}

// IsEquilibriumReached returns whether equilibrium has been reached
func (sc *SimulationController) IsEquilibriumReached() bool {
	return sc.equilibriumReached
}

// Initialize sets up the initial workforce based on configuration and validates parameters
// Returns an error if the configuration is invalid or initialization fails
func (sc *SimulationController) Initialize() error {
	// Validate configuration parameters
	if err := sc.validateConfiguration(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}
	
	// Reset simulation state
	sc.currentTimeStep = 0
	sc.timeSeries = make([]types.SimulationState, 0)
	sc.totalCatastrophicFailures = 0
	sc.equilibriumReached = false
	
	// Create initial workforce based on configuration
	if err := sc.createInitialWorkforce(); err != nil {
		return fmt.Errorf("failed to create initial workforce: %w", err)
	}
	
	// Record initial state
	initialState := sc.captureCurrentState()
	sc.timeSeries = append(sc.timeSeries, initialState)
	
	return nil
}

// validateConfiguration checks if all configuration parameters are valid
func (sc *SimulationController) validateConfiguration() error {
	config := sc.config
	
	// Check initial humans count
	if config.InitialHumans <= 0 {
		return errors.New("initial humans count must be greater than 0")
	}
	
	// Check experience distribution sums to 100%
	expSum := config.ExperienceDistribution.UniversityHire +
		config.ExperienceDistribution.MidLevel +
		config.ExperienceDistribution.Senior +
		config.ExperienceDistribution.Executive
	if expSum < 99.9 || expSum > 100.1 { // Allow small floating point errors
		return fmt.Errorf("experience distribution must sum to 100%%, got %.2f%%", expSum)
	}
	
	// Check cost category distribution sums to 100%
	costSum := config.CostCategoryDistribution.HighCostUS +
		config.CostCategoryDistribution.LowCostNonUS
	if costSum < 99.9 || costSum > 100.1 { // Allow small floating point errors
		return fmt.Errorf("cost category distribution must sum to 100%%, got %.2f%%", costSum)
	}
	
	// Check fixed budget is positive
	if config.FixedBudget <= 0 {
		return errors.New("fixed budget must be greater than 0")
	}
	
	// Check AI learning speeds are positive
	if config.AILearningSpeeds.UniversityToMid <= 0 ||
		config.AILearningSpeeds.MidToSenior <= 0 ||
		config.AILearningSpeeds.SeniorToExecutive <= 0 {
		return errors.New("AI learning speeds must be greater than 0")
	}
	
	// Check attrition rate is valid (0-100%)
	if config.AttritionConfig.NaturalRate < 0 || config.AttritionConfig.NaturalRate > 100 {
		return fmt.Errorf("natural attrition rate must be between 0-100%%, got %.2f%%", config.AttritionConfig.NaturalRate)
	}
	
	// Check forced acceleration is positive
	if config.AttritionConfig.ForcedAcceleration < 0 {
		return errors.New("forced acceleration must be non-negative")
	}
	
	// Check catastrophic failure rate is valid (0-1)
	if config.CatastrophicFailureRate < 0 || config.CatastrophicFailureRate > 1 {
		return fmt.Errorf("catastrophic failure rate must be between 0-1, got %.4f", config.CatastrophicFailureRate)
	}
	
	// Check time zone inefficiency is valid (0-1)
	if config.TimeZoneInefficiency < 0 || config.TimeZoneInefficiency > 1 {
		return fmt.Errorf("time zone inefficiency must be between 0-1, got %.4f", config.TimeZoneInefficiency)
	}
	
	return nil
}

// createInitialWorkforce creates the initial human workforce based on configuration
func (sc *SimulationController) createInitialWorkforce() error {
	config := sc.config
	
	// Calculate number of workers for each experience level
	expDist := config.ExperienceDistribution
	universityHireCount := int(float64(config.InitialHumans) * expDist.UniversityHire / 100.0)
	midLevelCount := int(float64(config.InitialHumans) * expDist.MidLevel / 100.0)
	seniorCount := int(float64(config.InitialHumans) * expDist.Senior / 100.0)
	executiveCount := int(float64(config.InitialHumans) * expDist.Executive / 100.0)
	
	// Handle rounding errors by adjusting the largest group
	totalAssigned := universityHireCount + midLevelCount + seniorCount + executiveCount
	if totalAssigned < config.InitialHumans {
		// Add remaining workers to the largest group
		remaining := config.InitialHumans - totalAssigned
		if expDist.UniversityHire >= expDist.MidLevel && expDist.UniversityHire >= expDist.Senior && expDist.UniversityHire >= expDist.Executive {
			universityHireCount += remaining
		} else if expDist.MidLevel >= expDist.Senior && expDist.MidLevel >= expDist.Executive {
			midLevelCount += remaining
		} else if expDist.Senior >= expDist.Executive {
			seniorCount += remaining
		} else {
			executiveCount += remaining
		}
	}
	
	// Calculate cost category distribution
	costDist := config.CostCategoryDistribution
	highCostCount := int(float64(config.InitialHumans) * costDist.HighCostUS / 100.0)
	
	// Create workers for each experience level
	experienceLevels := []struct {
		level types.ExperienceLevel
		count int
	}{
		{types.UniversityHire, universityHireCount},
		{types.MidLevel, midLevelCount},
		{types.Senior, seniorCount},
		{types.Executive, executiveCount},
	}
	
	businessOwnerAssigned := false
	
	for _, expLevel := range experienceLevels {
		for i := 0; i < expLevel.count; i++ {
			// Determine cost category for this worker
			var costCategory types.CostCategory
			if highCostCount > 0 {
				costCategory = types.HighCostUS
				highCostCount--
			} else {
				costCategory = types.LowCostNonUS
			}
			
			// Assign business owner to the first worker if not yet assigned
			isBusinessOwner := !businessOwnerAssigned
			if isBusinessOwner {
				businessOwnerAssigned = true
			}
			
			// Create the human worker
			_, err := sc.workforceManager.AddHuman(expLevel.level, costCategory, isBusinessOwner)
			if err != nil {
				return fmt.Errorf("failed to add human worker: %w", err)
			}
		}
	}
	
	// Ensure at least one business owner exists (requirement 1.9)
	if !businessOwnerAssigned {
		return errors.New("no business owner was assigned during workforce creation")
	}
	
	// Validate that initial workforce fits within budget
	humans := sc.workforceManager.GetAllHumans()
	agents := sc.workforceManager.GetAllAIAgents()
	totalCost := sc.economicModel.CalculateWorkforceCost(humans, agents)
	
	if totalCost > config.FixedBudget {
		return fmt.Errorf("initial workforce cost (%.2f) exceeds fixed budget (%.2f)", totalCost, config.FixedBudget)
	}
	
	return nil
}

// captureCurrentState captures the current simulation state for recording
func (sc *SimulationController) captureCurrentState() types.SimulationState {
	humans := sc.workforceManager.GetAllHumans()
	agents := sc.workforceManager.GetAllAIAgents()
	
	// Calculate metrics
	totalCost := sc.economicModel.CalculateWorkforceCost(humans, agents)
	availableBudget := sc.economicModel.GetAvailableBudget(humans, agents)
	totalProductivity := sc.workforceManager.CalculateTotalProductivity(sc.config.TimeZoneInefficiency)
	revenueOutput := sc.economicModel.CalculateRevenue(totalProductivity, sc.currentTimeStep)
	
	// Get workforce composition
	workforce := sc.workforceManager.GetWorkforceComposition()
	
	return types.SimulationState{
		TimeStep:             sc.currentTimeStep,
		Workforce:            workforce,
		TotalCost:            totalCost,
		AvailableBudget:      availableBudget,
		TotalProductivity:    totalProductivity,
		RevenueOutput:        revenueOutput,
		IsEquilibrium:        sc.equilibriumReached,
		CatastrophicFailures: sc.totalCatastrophicFailures,
	}
}

// Step executes one simulation time step
// Processes attrition, learning, optimization, and metrics according to requirements 10.2-10.7
func (sc *SimulationController) Step() types.SimulationState {
	// Increment time step
	sc.currentTimeStep++
	
	// Step 1: Process human attrition events (Requirement 10.2)
	sc.processAttrition()
	
	// Step 2: Update AI agent experience and learning progression (Requirement 10.3)
	sc.processLearning()
	
	// Step 3: Handle catastrophic failures
	sc.processCatastrophicFailures()
	
	// Step 4: Evaluate and execute workforce composition changes (Requirement 10.4)
	sc.processWorkforceOptimization()
	
	// Step 5: Calculate current revenue output (Requirement 10.5)
	// This is done in captureCurrentState()
	
	// Step 6: Check for equilibrium conditions (Requirement 10.6)
	sc.checkEquilibrium()
	
	// Step 7: Record workforce state and metrics at each time step (Requirement 10.7)
	currentState := sc.captureCurrentState()
	sc.timeSeries = append(sc.timeSeries, currentState)
	
	return currentState
}

// processAttrition handles human worker attrition based on configured attrition type
func (sc *SimulationController) processAttrition() {
	humans := sc.workforceManager.GetAllHumans()
	workersToRemove := sc.eventProcessor.ProcessAttrition(humans, sc.currentTimeStep)
	
	// Remove the selected workers
	for _, workerID := range workersToRemove {
		err := sc.workforceManager.RemoveHuman(workerID)
		if err != nil {
			// Log error but continue simulation
			// In a production system, this would use proper logging
			fmt.Printf("Warning: Failed to remove human worker %s: %v\n", workerID, err)
		}
	}
}

// processLearning updates AI agent experience and triggers level-ups
func (sc *SimulationController) processLearning() {
	agents := sc.workforceManager.GetAllAIAgents()
	// Process learning with time delta of 1 (one time step)
	sc.eventProcessor.ProcessLearning(agents, 1)
}

// processCatastrophicFailures generates and handles catastrophic failure events
func (sc *SimulationController) processCatastrophicFailures() {
	// Generate potential catastrophic failure
	failure := sc.eventProcessor.GenerateCatastrophicFailure(sc.currentTimeStep)
	if failure != nil {
		sc.totalCatastrophicFailures++
		
		// Evaluate workforce response to the failure
		humans := sc.workforceManager.GetAllHumans()
		agents := sc.workforceManager.GetAllAIAgents()
		outcome := sc.eventProcessor.EvaluateFailureResponse(failure, humans, agents)
		
		// Apply productivity penalties if workforce cannot handle the failure
		if !outcome.CanHandle && outcome.ProductivityPenalty > 0 {
			// In a more sophisticated implementation, we would apply temporary
			// productivity penalties. For now, we track the failure count.
			// The penalty could be applied by modifying productivity calculations
			// in subsequent steps, but this would require additional state tracking.
		}
	}
}

// processWorkforceOptimization evaluates and executes workforce composition changes
func (sc *SimulationController) processWorkforceOptimization() {
	humans := sc.workforceManager.GetAllHumans()
	agents := sc.workforceManager.GetAllAIAgents()
	
	// Calculate available budget and orchestration capacity
	availableBudget := sc.economicModel.GetAvailableBudget(humans, agents)
	availableCapacity := sc.workforceManager.GetAvailableOrchestrationCapacity()
	
	// Get optimization recommendations
	changes := sc.eventProcessor.OptimizeWorkforce(humans, agents, availableBudget, availableCapacity)
	
	// Execute agent releases first (to free up budget)
	for _, agentID := range changes.ReleaseAIAgents {
		err := sc.workforceManager.ReleaseAIAgent(agentID)
		if err != nil {
			fmt.Printf("Warning: Failed to release AI agent %s: %v\n", agentID, err)
		}
	}
	
	// Execute agent hires
	if changes.HireAIAgents > 0 && changes.OrchestratorID != "" {
		for i := 0; i < changes.HireAIAgents; i++ {
			_, err := sc.workforceManager.AddAIAgent(changes.OrchestratorID, sc.currentTimeStep)
			if err != nil {
				// If we can't hire more agents, stop trying
				fmt.Printf("Warning: Failed to hire AI agent: %v\n", err)
				break
			}
		}
	}
}

// checkEquilibrium determines if equilibrium conditions have been met
func (sc *SimulationController) checkEquilibrium() {
	// Simple equilibrium detection: check if workforce composition has been stable
	// for the last few time steps
	
	const stabilityWindow = 5 // Number of time steps to check for stability
	
	if len(sc.timeSeries) < stabilityWindow {
		// Not enough history to determine stability
		return
	}
	
	// Get the last few states
	recentStates := sc.timeSeries[len(sc.timeSeries)-stabilityWindow:]
	
	// Check if workforce composition has remained stable
	firstState := recentStates[0]
	isStable := true
	
	for i := 1; i < len(recentStates); i++ {
		state := recentStates[i]
		
		// Compare workforce composition
		if state.Workforce.Humans.Total != firstState.Workforce.Humans.Total ||
			state.Workforce.AIAgents.Total != firstState.Workforce.AIAgents.Total {
			isStable = false
			break
		}
		
		// Check if available budget is too low to hire more agents
		// (indicating cost-effectiveness equilibrium)
		if state.AvailableBudget > 0 {
			// Still have budget, check if we have orchestration capacity
			if state.Workforce.OrchestrationUtilization < 100.0 {
				// Have both budget and capacity, but no hiring occurred
				// This suggests equilibrium has been reached
				continue
			}
		}
	}
	
	// Additional check: if we have no available orchestration capacity
	// and no budget for more humans, we've reached equilibrium
	currentState := sc.captureCurrentState()
	if currentState.Workforce.OrchestrationUtilization >= 100.0 || currentState.AvailableBudget <= 0 {
		isStable = true
	}
	
	sc.equilibriumReached = isStable
}
// IsEquilibrium detects when equilibrium conditions are met
// Checks workforce composition stability according to requirements 8.1, 8.2, 8.3
func (sc *SimulationController) IsEquilibrium() bool {
	// Return the cached equilibrium state (updated in checkEquilibrium)
	return sc.equilibriumReached
}

// IsEquilibriumDetailed provides detailed equilibrium analysis
// This method provides more granular equilibrium detection logic
func (sc *SimulationController) IsEquilibriumDetailed() (bool, string) {
	if len(sc.timeSeries) == 0 {
		return false, "no simulation data available"
	}
	
	currentState := sc.timeSeries[len(sc.timeSeries)-1]
	
	// Check if we have reached maximum orchestration capacity
	if currentState.Workforce.OrchestrationUtilization >= 100.0 {
		return true, "maximum orchestration capacity reached"
	}
	
	// Check if we have no available budget for more agents
	if currentState.AvailableBudget <= 0 {
		return true, "no available budget for workforce expansion"
	}
	
	// Check if the cost of adding additional AI agents exceeds productivity benefit
	// This requires checking if we have budget and capacity but no hiring occurred
	const stabilityWindow = 5
	if len(sc.timeSeries) >= stabilityWindow {
		recentStates := sc.timeSeries[len(sc.timeSeries)-stabilityWindow:]
		
		// Check if workforce composition has been stable
		firstState := recentStates[0]
		isStable := true
		
		for i := 1; i < len(recentStates); i++ {
			state := recentStates[i]
			if state.Workforce.Humans.Total != firstState.Workforce.Humans.Total ||
				state.Workforce.AIAgents.Total != firstState.Workforce.AIAgents.Total {
				isStable = false
				break
			}
		}
		
		if isStable {
			// Check if we had opportunities to hire but didn't
			hasOpportunity := false
			for _, state := range recentStates {
				if state.AvailableBudget > types.AIAgentCosts[types.UniversityHire] &&
					state.Workforce.OrchestrationUtilization < 100.0 {
					hasOpportunity = true
					break
				}
			}
			
			if hasOpportunity {
				return true, "workforce composition stable despite hiring opportunities (cost-effectiveness equilibrium)"
			} else {
				return true, "workforce composition stable with no hiring opportunities"
			}
		}
	}
	
	return false, "equilibrium conditions not yet met"
}
// RunUntilEquilibrium executes the simulation loop until equilibrium is reached
// Returns complete simulation result according to requirements 8.3, 8.4
func (sc *SimulationController) RunUntilEquilibrium(maxTimeSteps int) (types.SimulationResult, error) {
	// Initialize the simulation if not already done
	if len(sc.timeSeries) == 0 {
		if err := sc.Initialize(); err != nil {
			return types.SimulationResult{}, fmt.Errorf("initialization failed: %w", err)
		}
	}
	
	// Execute simulation steps until equilibrium or max steps reached
	for sc.currentTimeStep < maxTimeSteps && !sc.equilibriumReached {
		sc.Step()
		
		// Safety check to prevent infinite loops
		if sc.currentTimeStep >= maxTimeSteps {
			break
		}
	}
	
	// Determine final equilibrium state
	var equilibriumState types.SimulationState
	if len(sc.timeSeries) > 0 {
		equilibriumState = sc.timeSeries[len(sc.timeSeries)-1]
		equilibriumState.IsEquilibrium = sc.equilibriumReached
	}
	
	// Create and return simulation result
	result := types.SimulationResult{
		Config:                    sc.config,
		TimeSeries:               sc.timeSeries,
		EquilibriumState:         equilibriumState,
		TimeToEquilibrium:        sc.currentTimeStep,
		TotalCatastrophicFailures: sc.totalCatastrophicFailures,
	}
	
	return result, nil
}

// Reset resets the simulation controller to initial state
// Useful for running multiple simulations with the same configuration
func (sc *SimulationController) Reset() {
	sc.currentTimeStep = 0
	sc.timeSeries = make([]types.SimulationState, 0)
	sc.totalCatastrophicFailures = 0
	sc.equilibriumReached = false
	
	// Reset component states
	sc.workforceManager = workforce.NewWorkforceManager()
	sc.economicModel = economic.NewEconomicModel(sc.config.FixedBudget, sc.config.RevenueScenario)
}