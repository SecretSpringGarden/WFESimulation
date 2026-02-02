package types

// ExperienceDistribution defines the percentage distribution of workers across experience levels
type ExperienceDistribution struct {
	UniversityHire float64 // percentage (0-100)
	MidLevel       float64 // percentage (0-100)
	Senior         float64 // percentage (0-100)
	Executive      float64 // percentage (0-100)
}

// CostCategoryDistribution defines the percentage distribution of workers across cost categories
type CostCategoryDistribution struct {
	HighCostUS    float64 // percentage (0-100)
	LowCostNonUS  float64 // percentage (0-100)
}

// AILearningSpeed defines the time steps required for AI agents to progress through experience levels
type AILearningSpeed struct {
	UniversityToMid int // time steps required
	MidToSenior     int // time steps required
	SeniorToExecutive int // time steps required
}

// AttritionConfig defines the attrition behavior for human workers
type AttritionConfig struct {
	Type                AttritionType
	NaturalRate         float64 // annual percentage (0-100)
	ForcedAcceleration  float64 // multiplier for attrition rate
}

// SimulationConfig contains all configuration parameters for a simulation run
type SimulationConfig struct {
	// Initial workforce configuration
	InitialHumans            int
	ExperienceDistribution   ExperienceDistribution
	CostCategoryDistribution CostCategoryDistribution
	
	// Economic configuration
	FixedBudget      float64
	RevenueScenario  RevenueScenario
	
	// AI learning configuration
	AILearningSpeeds AILearningSpeed
	
	// Attrition configuration
	AttritionConfig AttritionConfig
	
	// Failure and inefficiency configuration
	CatastrophicFailureRate float64 // probability per time step (0-1)
	TimeZoneInefficiency    float64 // productivity penalty for Low_Cost_Non_US (0-1)
}

// Validate checks if the configuration is valid
func (c *SimulationConfig) Validate() error {
	// Validation will be implemented as part of the simulation controller
	// This method is a placeholder for now
	return nil
}

// WorkforceComposition represents detailed workforce statistics
type WorkforceComposition struct {
	Humans struct {
		Total          int
		ByExperience   map[ExperienceLevel]int
		ByCostCategory map[CostCategory]int
	}
	AIAgents struct {
		Total        int
		ByExperience map[ExperienceLevel]int
	}
	OrchestrationUtilization float64 // percentage of capacity used (0-100)
}

// SimulationState represents the state of the simulation at a specific time step
type SimulationState struct {
	TimeStep                  int
	Workforce                 WorkforceComposition
	TotalCost                 float64
	AvailableBudget          float64
	TotalProductivity        float64
	RevenueOutput            float64
	IsEquilibrium            bool
	CatastrophicFailures     int
}

// SimulationResult represents the complete result of a simulation run
type SimulationResult struct {
	Config                    SimulationConfig
	TimeSeries               []SimulationState
	EquilibriumState         SimulationState
	TimeToEquilibrium        int
	TotalCatastrophicFailures int
}
