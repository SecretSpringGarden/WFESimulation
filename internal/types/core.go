package types

// ExperienceLevel represents the capability classification of workers
type ExperienceLevel int

const (
	UniversityHire ExperienceLevel = iota
	MidLevel
	Senior
	Executive
)

// String returns the string representation of ExperienceLevel
func (e ExperienceLevel) String() string {
	switch e {
	case UniversityHire:
		return "University_Hire"
	case MidLevel:
		return "Mid_Level"
	case Senior:
		return "Senior"
	case Executive:
		return "Executive"
	default:
		return "Unknown"
	}
}

// CostCategory represents the cost classification of human workers
type CostCategory int

const (
	HighCostUS CostCategory = iota
	LowCostNonUS
)

// String returns the string representation of CostCategory
func (c CostCategory) String() string {
	switch c {
	case HighCostUS:
		return "High_Cost_US"
	case LowCostNonUS:
		return "Low_Cost_Non_US"
	default:
		return "Unknown"
	}
}

// RevenueScenario represents the revenue growth pattern
type RevenueScenario int

const (
	FlatRevenue RevenueScenario = iota
	ExplosiveGrowth
)

// String returns the string representation of RevenueScenario
func (r RevenueScenario) String() string {
	switch r {
	case FlatRevenue:
		return "Flat_Revenue"
	case ExplosiveGrowth:
		return "Explosive_Growth"
	default:
		return "Unknown"
	}
}

// AttritionType represents the type of human worker attrition
type AttritionType int

const (
	NaturalAttrition AttritionType = iota
	HiringFreeze
	ReductionInForce
)

// String returns the string representation of AttritionType
func (a AttritionType) String() string {
	switch a {
	case NaturalAttrition:
		return "Natural_Attrition"
	case HiringFreeze:
		return "Hiring_Freeze"
	case ReductionInForce:
		return "Reduction_In_Force"
	default:
		return "Unknown"
	}
}

// OrchestrationLimit is the maximum number of AI agents a single human can manage
const OrchestrationLimit = 6

// Cost and productivity values for human workers based on experience level and cost category
var (
	// BaseCosts maps experience level and cost category to annual cost
	BaseCosts = map[ExperienceLevel]map[CostCategory]float64{
		UniversityHire: {
			HighCostUS:   100000,
			LowCostNonUS: 40000,
		},
		MidLevel: {
			HighCostUS:   150000,
			LowCostNonUS: 60000,
		},
		Senior: {
			HighCostUS:   200000,
			LowCostNonUS: 80000,
		},
		Executive: {
			HighCostUS:   300000,
			LowCostNonUS: 120000,
		},
	}

	// BaseProductivity maps experience level to productivity value
	BaseProductivity = map[ExperienceLevel]float64{
		UniversityHire: 1.0,
		MidLevel:       2.0,
		Senior:         3.5,
		Executive:      5.0,
	}
)

// HumanWorker represents a human employee in the workforce
type HumanWorker struct {
	ID               string
	ExperienceLevel  ExperienceLevel
	CostCategory     CostCategory
	BaseCost         float64
	BaseProductivity float64
	AssignedAgents   []string // IDs of assigned AI agents
	IsBusinessOwner  bool
}

// NewHumanWorker creates a new HumanWorker with attributes assigned based on experience level and cost category
func NewHumanWorker(id string, experienceLevel ExperienceLevel, costCategory CostCategory, isBusinessOwner bool) *HumanWorker {
	baseCost := BaseCosts[experienceLevel][costCategory]
	baseProductivity := BaseProductivity[experienceLevel]

	return &HumanWorker{
		ID:               id,
		ExperienceLevel:  experienceLevel,
		CostCategory:     costCategory,
		BaseCost:         baseCost,
		BaseProductivity: baseProductivity,
		AssignedAgents:   make([]string, 0),
		IsBusinessOwner:  isBusinessOwner,
	}
}

// GetEffectiveProductivity calculates the effective productivity of the human worker
// applying time zone inefficiency penalty for Low_Cost_Non_US workers
func (h *HumanWorker) GetEffectiveProductivity(timeZoneInefficiency float64) float64 {
	if h.CostCategory == LowCostNonUS {
		return h.BaseProductivity * (1.0 - timeZoneInefficiency)
	}
	return h.BaseProductivity
}

// CanOrchestrateMoreAgents checks if the human worker can orchestrate additional AI agents
func (h *HumanWorker) CanOrchestrateMoreAgents() bool {
	return len(h.AssignedAgents) < OrchestrationLimit
}

// GetOrchestrationCapacity returns the number of additional AI agents this human can orchestrate
func (h *HumanWorker) GetOrchestrationCapacity() int {
	return OrchestrationLimit - len(h.AssignedAgents)
}

// AI Agent cost and productivity values based on experience level
var (
	// AIAgentCosts maps experience level to annual cost for AI agents
	AIAgentCosts = map[ExperienceLevel]float64{
		UniversityHire: 20000,
		MidLevel:       40000,
		Senior:         70000,
		Executive:      100000,
	}

	// AIAgentProductivity maps experience level to productivity value for AI agents
	AIAgentProductivity = map[ExperienceLevel]float64{
		UniversityHire: 0.8,
		MidLevel:       1.8,
		Senior:         3.2,
		Executive:      4.8,
	}
)

// AIAgent represents an AI agent in the workforce
type AIAgent struct {
	ID              string
	ExperienceLevel ExperienceLevel
	ExperiencePoints float64
	Cost            float64
	OrchestratorID  string
	CreationTime    int // time step when the agent was created
}

// NewAIAgent creates a new AIAgent initialized at University_Hire level
func NewAIAgent(id string, orchestratorID string, creationTime int) *AIAgent {
	return &AIAgent{
		ID:              id,
		ExperienceLevel: UniversityHire,
		ExperiencePoints: 0.0,
		Cost:            AIAgentCosts[UniversityHire],
		OrchestratorID:  orchestratorID,
		CreationTime:    creationTime,
	}
}

// AccumulateExperience calculates and adds experience points based on time and data exposure
// timeDelta is the number of time steps elapsed
// dataExposure is a multiplier representing the amount of data the agent has been exposed to (typically 1.0)
func (a *AIAgent) AccumulateExperience(timeDelta int, dataExposure float64) {
	// Experience accumulation is proportional to time and data exposure
	experienceGain := float64(timeDelta) * dataExposure
	a.ExperiencePoints += experienceGain
}

// CheckLevelUp checks if the agent has accumulated enough experience to progress to the next level
// Returns true if a level up occurred
// learningSpeed contains the thresholds for each level progression
func (a *AIAgent) CheckLevelUp(learningSpeed AILearningSpeed) bool {
	var threshold float64
	var nextLevel ExperienceLevel
	
	switch a.ExperienceLevel {
	case UniversityHire:
		threshold = float64(learningSpeed.UniversityToMid)
		nextLevel = MidLevel
	case MidLevel:
		threshold = float64(learningSpeed.MidToSenior)
		nextLevel = Senior
	case Senior:
		threshold = float64(learningSpeed.SeniorToExecutive)
		nextLevel = Executive
	case Executive:
		// Already at max level
		return false
	default:
		return false
	}
	
	// Check if experience points exceed the threshold
	if a.ExperiencePoints >= threshold {
		a.ExperienceLevel = nextLevel
		a.ExperiencePoints = 0.0 // Reset experience points for the new level
		// Update cost based on new experience level
		a.Cost = AIAgentCosts[nextLevel]
		return true
	}
	
	return false
}

// GetProductivity returns the productivity value based on the agent's current experience level
func (a *AIAgent) GetProductivity() float64 {
	return AIAgentProductivity[a.ExperienceLevel]
}

// GetCost returns the cost of the agent based on their current experience level
func (a *AIAgent) GetCost() float64 {
	return a.Cost
}
