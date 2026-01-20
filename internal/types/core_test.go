package types

import (
	"testing"
)

func TestExperienceLevelString(t *testing.T) {
	tests := []struct {
		level    ExperienceLevel
		expected string
	}{
		{UniversityHire, "University_Hire"},
		{MidLevel, "Mid_Level"},
		{Senior, "Senior"},
		{Executive, "Executive"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.level.String(); got != tt.expected {
				t.Errorf("ExperienceLevel.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCostCategoryString(t *testing.T) {
	tests := []struct {
		category CostCategory
		expected string
	}{
		{HighCostUS, "High_Cost_US"},
		{LowCostNonUS, "Low_Cost_Non_US"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.category.String(); got != tt.expected {
				t.Errorf("CostCategory.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRevenueScenarioString(t *testing.T) {
	tests := []struct {
		scenario RevenueScenario
		expected string
	}{
		{FlatRevenue, "Flat_Revenue"},
		{ExplosiveGrowth, "Explosive_Growth"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.scenario.String(); got != tt.expected {
				t.Errorf("RevenueScenario.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAttritionTypeString(t *testing.T) {
	tests := []struct {
		attrition AttritionType
		expected  string
	}{
		{NaturalAttrition, "Natural_Attrition"},
		{HiringFreeze, "Hiring_Freeze"},
		{ReductionInForce, "Reduction_In_Force"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.attrition.String(); got != tt.expected {
				t.Errorf("AttritionType.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestOrchestrationLimit(t *testing.T) {
	if OrchestrationLimit != 6 {
		t.Errorf("OrchestrationLimit = %v, want 6", OrchestrationLimit)
	}
}

func TestNewHumanWorker(t *testing.T) {
	tests := []struct {
		name             string
		experienceLevel  ExperienceLevel
		costCategory     CostCategory
		isBusinessOwner  bool
		expectedCost     float64
		expectedProd     float64
	}{
		{
			name:             "University Hire High Cost US",
			experienceLevel:  UniversityHire,
			costCategory:     HighCostUS,
			isBusinessOwner:  false,
			expectedCost:     100000,
			expectedProd:     1.0,
		},
		{
			name:             "Mid Level Low Cost Non-US",
			experienceLevel:  MidLevel,
			costCategory:     LowCostNonUS,
			isBusinessOwner:  false,
			expectedCost:     60000,
			expectedProd:     2.0,
		},
		{
			name:             "Senior High Cost US Business Owner",
			experienceLevel:  Senior,
			costCategory:     HighCostUS,
			isBusinessOwner:  true,
			expectedCost:     200000,
			expectedProd:     3.5,
		},
		{
			name:             "Executive Low Cost Non-US",
			experienceLevel:  Executive,
			costCategory:     LowCostNonUS,
			isBusinessOwner:  false,
			expectedCost:     120000,
			expectedProd:     5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			worker := NewHumanWorker("test-id", tt.experienceLevel, tt.costCategory, tt.isBusinessOwner)
			
			if worker.ID != "test-id" {
				t.Errorf("ID = %v, want test-id", worker.ID)
			}
			if worker.ExperienceLevel != tt.experienceLevel {
				t.Errorf("ExperienceLevel = %v, want %v", worker.ExperienceLevel, tt.experienceLevel)
			}
			if worker.CostCategory != tt.costCategory {
				t.Errorf("CostCategory = %v, want %v", worker.CostCategory, tt.costCategory)
			}
			if worker.BaseCost != tt.expectedCost {
				t.Errorf("BaseCost = %v, want %v", worker.BaseCost, tt.expectedCost)
			}
			if worker.BaseProductivity != tt.expectedProd {
				t.Errorf("BaseProductivity = %v, want %v", worker.BaseProductivity, tt.expectedProd)
			}
			if worker.IsBusinessOwner != tt.isBusinessOwner {
				t.Errorf("IsBusinessOwner = %v, want %v", worker.IsBusinessOwner, tt.isBusinessOwner)
			}
			if len(worker.AssignedAgents) != 0 {
				t.Errorf("AssignedAgents length = %v, want 0", len(worker.AssignedAgents))
			}
		})
	}
}

func TestGetEffectiveProductivity(t *testing.T) {
	tests := []struct {
		name                 string
		costCategory         CostCategory
		baseProductivity     float64
		timeZoneInefficiency float64
		expected             float64
	}{
		{
			name:                 "High Cost US - no penalty",
			costCategory:         HighCostUS,
			baseProductivity:     3.5,
			timeZoneInefficiency: 0.2,
			expected:             3.5,
		},
		{
			name:                 "Low Cost Non-US - with penalty",
			costCategory:         LowCostNonUS,
			baseProductivity:     3.5,
			timeZoneInefficiency: 0.2,
			expected:             2.8, // 3.5 * (1 - 0.2) = 2.8
		},
		{
			name:                 "Low Cost Non-US - zero penalty",
			costCategory:         LowCostNonUS,
			baseProductivity:     2.0,
			timeZoneInefficiency: 0.0,
			expected:             2.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			worker := &HumanWorker{
				CostCategory:     tt.costCategory,
				BaseProductivity: tt.baseProductivity,
			}
			
			got := worker.GetEffectiveProductivity(tt.timeZoneInefficiency)
			// Use a small tolerance for floating point comparison
			const tolerance = 1e-9
			if diff := got - tt.expected; diff < -tolerance || diff > tolerance {
				t.Errorf("GetEffectiveProductivity() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestOrchestrationCapacity(t *testing.T) {
	worker := NewHumanWorker("test-id", MidLevel, HighCostUS, false)
	
	// Initially should have full capacity
	if !worker.CanOrchestrateMoreAgents() {
		t.Error("CanOrchestrateMoreAgents() = false, want true")
	}
	if got := worker.GetOrchestrationCapacity(); got != 6 {
		t.Errorf("GetOrchestrationCapacity() = %v, want 6", got)
	}
	
	// Add 3 agents
	worker.AssignedAgents = []string{"agent1", "agent2", "agent3"}
	if !worker.CanOrchestrateMoreAgents() {
		t.Error("CanOrchestrateMoreAgents() = false, want true")
	}
	if got := worker.GetOrchestrationCapacity(); got != 3 {
		t.Errorf("GetOrchestrationCapacity() = %v, want 3", got)
	}
	
	// Add 3 more agents (total 6)
	worker.AssignedAgents = append(worker.AssignedAgents, "agent4", "agent5", "agent6")
	if worker.CanOrchestrateMoreAgents() {
		t.Error("CanOrchestrateMoreAgents() = true, want false")
	}
	if got := worker.GetOrchestrationCapacity(); got != 0 {
		t.Errorf("GetOrchestrationCapacity() = %v, want 0", got)
	}
}

func TestNewAIAgent(t *testing.T) {
	agent := NewAIAgent("agent-1", "orchestrator-1", 10)
	
	if agent.ID != "agent-1" {
		t.Errorf("ID = %v, want agent-1", agent.ID)
	}
	if agent.ExperienceLevel != UniversityHire {
		t.Errorf("ExperienceLevel = %v, want UniversityHire", agent.ExperienceLevel)
	}
	if agent.ExperiencePoints != 0.0 {
		t.Errorf("ExperiencePoints = %v, want 0.0", agent.ExperiencePoints)
	}
	if agent.Cost != AIAgentCosts[UniversityHire] {
		t.Errorf("Cost = %v, want %v", agent.Cost, AIAgentCosts[UniversityHire])
	}
	if agent.OrchestratorID != "orchestrator-1" {
		t.Errorf("OrchestratorID = %v, want orchestrator-1", agent.OrchestratorID)
	}
	if agent.CreationTime != 10 {
		t.Errorf("CreationTime = %v, want 10", agent.CreationTime)
	}
}

func TestAccumulateExperience(t *testing.T) {
	agent := NewAIAgent("agent-1", "orchestrator-1", 0)
	
	// Initially should have 0 experience
	if agent.ExperiencePoints != 0.0 {
		t.Errorf("Initial ExperiencePoints = %v, want 0.0", agent.ExperiencePoints)
	}
	
	// Accumulate experience with time delta 5 and data exposure 1.0
	agent.AccumulateExperience(5, 1.0)
	if agent.ExperiencePoints != 5.0 {
		t.Errorf("ExperiencePoints after first accumulation = %v, want 5.0", agent.ExperiencePoints)
	}
	
	// Accumulate more experience
	agent.AccumulateExperience(3, 2.0)
	expected := 5.0 + (3.0 * 2.0) // 5.0 + 6.0 = 11.0
	if agent.ExperiencePoints != expected {
		t.Errorf("ExperiencePoints after second accumulation = %v, want %v", agent.ExperiencePoints, expected)
	}
}

func TestCheckLevelUp(t *testing.T) {
	learningSpeed := AILearningSpeed{
		UniversityToMid:   10,
		MidToSenior:       20,
		SeniorToExecutive: 30,
	}
	
	agent := NewAIAgent("agent-1", "orchestrator-1", 0)
	
	// Initially at UniversityHire
	if agent.ExperienceLevel != UniversityHire {
		t.Errorf("Initial ExperienceLevel = %v, want UniversityHire", agent.ExperienceLevel)
	}
	
	// Accumulate experience but not enough to level up
	agent.AccumulateExperience(5, 1.0)
	if leveledUp := agent.CheckLevelUp(learningSpeed); leveledUp {
		t.Error("CheckLevelUp() = true, want false (not enough experience)")
	}
	if agent.ExperienceLevel != UniversityHire {
		t.Errorf("ExperienceLevel = %v, want UniversityHire", agent.ExperienceLevel)
	}
	
	// Accumulate enough experience to level up to MidLevel
	agent.AccumulateExperience(5, 1.0) // Total: 10
	if leveledUp := agent.CheckLevelUp(learningSpeed); !leveledUp {
		t.Error("CheckLevelUp() = false, want true (should level up)")
	}
	if agent.ExperienceLevel != MidLevel {
		t.Errorf("ExperienceLevel = %v, want MidLevel", agent.ExperienceLevel)
	}
	if agent.ExperiencePoints != 0.0 {
		t.Errorf("ExperiencePoints after level up = %v, want 0.0", agent.ExperiencePoints)
	}
	if agent.Cost != AIAgentCosts[MidLevel] {
		t.Errorf("Cost after level up = %v, want %v", agent.Cost, AIAgentCosts[MidLevel])
	}
	
	// Level up to Senior
	agent.AccumulateExperience(20, 1.0)
	if leveledUp := agent.CheckLevelUp(learningSpeed); !leveledUp {
		t.Error("CheckLevelUp() = false, want true (should level up to Senior)")
	}
	if agent.ExperienceLevel != Senior {
		t.Errorf("ExperienceLevel = %v, want Senior", agent.ExperienceLevel)
	}
	
	// Level up to Executive
	agent.AccumulateExperience(30, 1.0)
	if leveledUp := agent.CheckLevelUp(learningSpeed); !leveledUp {
		t.Error("CheckLevelUp() = false, want true (should level up to Executive)")
	}
	if agent.ExperienceLevel != Executive {
		t.Errorf("ExperienceLevel = %v, want Executive", agent.ExperienceLevel)
	}
	
	// Try to level up beyond Executive (should not level up)
	agent.AccumulateExperience(100, 1.0)
	if leveledUp := agent.CheckLevelUp(learningSpeed); leveledUp {
		t.Error("CheckLevelUp() = true, want false (already at max level)")
	}
	if agent.ExperienceLevel != Executive {
		t.Errorf("ExperienceLevel = %v, want Executive", agent.ExperienceLevel)
	}
}

func TestGetProductivity(t *testing.T) {
	tests := []struct {
		name            string
		experienceLevel ExperienceLevel
		expected        float64
	}{
		{"University Hire", UniversityHire, 0.8},
		{"Mid Level", MidLevel, 1.8},
		{"Senior", Senior, 3.2},
		{"Executive", Executive, 4.8},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := NewAIAgent("agent-1", "orchestrator-1", 0)
			agent.ExperienceLevel = tt.experienceLevel
			
			got := agent.GetProductivity()
			if got != tt.expected {
				t.Errorf("GetProductivity() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetCost(t *testing.T) {
	tests := []struct {
		name            string
		experienceLevel ExperienceLevel
		expected        float64
	}{
		{"University Hire", UniversityHire, 20000},
		{"Mid Level", MidLevel, 40000},
		{"Senior", Senior, 70000},
		{"Executive", Executive, 100000},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := NewAIAgent("agent-1", "orchestrator-1", 0)
			agent.ExperienceLevel = tt.experienceLevel
			agent.Cost = AIAgentCosts[tt.experienceLevel]
			
			got := agent.GetCost()
			if got != tt.expected {
				t.Errorf("GetCost() = %v, want %v", got, tt.expected)
			}
		})
	}
}
