package controller

import (
	"testing"
	"workforce-ai-transition-simulator/internal/types"
)

func TestNewSimulationController(t *testing.T) {
	config := types.SimulationConfig{
		InitialHumans: 10,
		ExperienceDistribution: types.ExperienceDistribution{
			UniversityHire: 40.0,
			MidLevel:       30.0,
			Senior:         20.0,
			Executive:      10.0,
		},
		CostCategoryDistribution: types.CostCategoryDistribution{
			HighCostUS:   60.0,
			LowCostNonUS: 40.0,
		},
		FixedBudget:     1000000.0,
		RevenueScenario: types.FlatRevenue,
		AILearningSpeeds: types.AILearningSpeed{
			UniversityToMid:   10,
			MidToSenior:       15,
			SeniorToExecutive: 20,
		},
		AttritionConfig: types.AttritionConfig{
			Type:               types.NaturalAttrition,
			NaturalRate:        10.0,
			ForcedAcceleration: 1.0,
		},
		CatastrophicFailureRate: 0.01,
		TimeZoneInefficiency:    0.1,
	}

	controller := NewSimulationController(config, 12345)

	if controller == nil {
		t.Fatal("NewSimulationController returned nil")
	}

	if controller.GetCurrentTimeStep() != 0 {
		t.Errorf("Expected initial time step to be 0, got %d", controller.GetCurrentTimeStep())
	}

	if controller.IsEquilibriumReached() {
		t.Error("Expected equilibrium not to be reached initially")
	}
}

func TestInitialize(t *testing.T) {
	config := types.SimulationConfig{
		InitialHumans: 5,
		ExperienceDistribution: types.ExperienceDistribution{
			UniversityHire: 40.0,
			MidLevel:       30.0,
			Senior:         20.0,
			Executive:      10.0,
		},
		CostCategoryDistribution: types.CostCategoryDistribution{
			HighCostUS:   60.0,
			LowCostNonUS: 40.0,
		},
		FixedBudget:     1000000.0,
		RevenueScenario: types.FlatRevenue,
		AILearningSpeeds: types.AILearningSpeed{
			UniversityToMid:   10,
			MidToSenior:       15,
			SeniorToExecutive: 20,
		},
		AttritionConfig: types.AttritionConfig{
			Type:               types.NaturalAttrition,
			NaturalRate:        10.0,
			ForcedAcceleration: 1.0,
		},
		CatastrophicFailureRate: 0.01,
		TimeZoneInefficiency:    0.1,
	}

	controller := NewSimulationController(config, 12345)

	err := controller.Initialize()
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Check that initial workforce was created
	timeSeries := controller.GetTimeSeries()
	if len(timeSeries) != 1 {
		t.Errorf("Expected 1 initial state, got %d", len(timeSeries))
	}

	initialState := timeSeries[0]
	if initialState.Workforce.Humans.Total != 5 {
		t.Errorf("Expected 5 initial humans, got %d", initialState.Workforce.Humans.Total)
	}

	if initialState.TimeStep != 0 {
		t.Errorf("Expected initial time step to be 0, got %d", initialState.TimeStep)
	}
}

func TestValidateConfiguration(t *testing.T) {
	tests := []struct {
		name        string
		config      types.SimulationConfig
		expectError bool
	}{
		{
			name: "valid configuration",
			config: types.SimulationConfig{
				InitialHumans: 10,
				ExperienceDistribution: types.ExperienceDistribution{
					UniversityHire: 40.0,
					MidLevel:       30.0,
					Senior:         20.0,
					Executive:      10.0,
				},
				CostCategoryDistribution: types.CostCategoryDistribution{
					HighCostUS:   60.0,
					LowCostNonUS: 40.0,
				},
				FixedBudget:     1000000.0,
				RevenueScenario: types.FlatRevenue,
				AILearningSpeeds: types.AILearningSpeed{
					UniversityToMid:   10,
					MidToSenior:       15,
					SeniorToExecutive: 20,
				},
				AttritionConfig: types.AttritionConfig{
					Type:               types.NaturalAttrition,
					NaturalRate:        10.0,
					ForcedAcceleration: 1.0,
				},
				CatastrophicFailureRate: 0.01,
				TimeZoneInefficiency:    0.1,
			},
			expectError: false,
		},
		{
			name: "zero initial humans",
			config: types.SimulationConfig{
				InitialHumans: 0,
			},
			expectError: true,
		},
		{
			name: "invalid experience distribution sum",
			config: types.SimulationConfig{
				InitialHumans: 10,
				ExperienceDistribution: types.ExperienceDistribution{
					UniversityHire: 40.0,
					MidLevel:       30.0,
					Senior:         20.0,
					Executive:      20.0, // Sum = 110%
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := NewSimulationController(tt.config, 12345)
			err := controller.validateConfiguration()

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestStep(t *testing.T) {
	config := types.SimulationConfig{
		InitialHumans: 3,
		ExperienceDistribution: types.ExperienceDistribution{
			UniversityHire: 50.0,
			MidLevel:       30.0,
			Senior:         20.0,
			Executive:      0.0,
		},
		CostCategoryDistribution: types.CostCategoryDistribution{
			HighCostUS:   100.0,
			LowCostNonUS: 0.0,
		},
		FixedBudget:     1000000.0,
		RevenueScenario: types.FlatRevenue,
		AILearningSpeeds: types.AILearningSpeed{
			UniversityToMid:   10,
			MidToSenior:       15,
			SeniorToExecutive: 20,
		},
		AttritionConfig: types.AttritionConfig{
			Type:               types.NaturalAttrition,
			NaturalRate:        0.0, // No attrition for predictable test
			ForcedAcceleration: 1.0,
		},
		CatastrophicFailureRate: 0.0, // No failures for predictable test
		TimeZoneInefficiency:    0.0,
	}

	controller := NewSimulationController(config, 12345)
	err := controller.Initialize()
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	initialTimeStep := controller.GetCurrentTimeStep()
	state := controller.Step()

	// Check that time step advanced
	if controller.GetCurrentTimeStep() != initialTimeStep+1 {
		t.Errorf("Expected time step to advance from %d to %d, got %d",
			initialTimeStep, initialTimeStep+1, controller.GetCurrentTimeStep())
	}

	// Check that state was recorded
	if state.TimeStep != controller.GetCurrentTimeStep() {
		t.Errorf("Expected state time step %d, got %d",
			controller.GetCurrentTimeStep(), state.TimeStep)
	}

	// Check that time series was updated
	timeSeries := controller.GetTimeSeries()
	if len(timeSeries) != 2 { // Initial state + 1 step
		t.Errorf("Expected 2 states in time series, got %d", len(timeSeries))
	}
}
func TestRunUntilEquilibrium(t *testing.T) {
	config := types.SimulationConfig{
		InitialHumans: 3,
		ExperienceDistribution: types.ExperienceDistribution{
			UniversityHire: 50.0,
			MidLevel:       30.0,
			Senior:         20.0,
			Executive:      0.0,
		},
		CostCategoryDistribution: types.CostCategoryDistribution{
			HighCostUS:   100.0,
			LowCostNonUS: 0.0,
		},
		FixedBudget:     500000.0, // Smaller budget to reach equilibrium faster
		RevenueScenario: types.FlatRevenue,
		AILearningSpeeds: types.AILearningSpeed{
			UniversityToMid:   5,
			MidToSenior:       10,
			SeniorToExecutive: 15,
		},
		AttritionConfig: types.AttritionConfig{
			Type:               types.NaturalAttrition,
			NaturalRate:        0.0, // No attrition for predictable test
			ForcedAcceleration: 1.0,
		},
		CatastrophicFailureRate: 0.0, // No failures for predictable test
		TimeZoneInefficiency:    0.0,
	}

	controller := NewSimulationController(config, 12345)
	
	result, err := controller.RunUntilEquilibrium(100) // Max 100 steps
	if err != nil {
		t.Fatalf("RunUntilEquilibrium failed: %v", err)
	}

	// Check that simulation ran
	if len(result.TimeSeries) == 0 {
		t.Error("Expected time series data, got empty")
	}

	// Check that we have initial humans
	if result.TimeSeries[0].Workforce.Humans.Total != 3 {
		t.Errorf("Expected 3 initial humans, got %d", result.TimeSeries[0].Workforce.Humans.Total)
	}

	// Check that time progressed
	if result.TimeToEquilibrium <= 0 {
		t.Errorf("Expected positive time to equilibrium, got %d", result.TimeToEquilibrium)
	}

	// Check that equilibrium state is marked correctly
	if !result.EquilibriumState.IsEquilibrium && result.TimeToEquilibrium < 100 {
		t.Error("Expected equilibrium state to be marked as equilibrium")
	}

	t.Logf("Simulation completed in %d time steps", result.TimeToEquilibrium)
	t.Logf("Final workforce: %d humans, %d AI agents", 
		result.EquilibriumState.Workforce.Humans.Total,
		result.EquilibriumState.Workforce.AIAgents.Total)
}

func TestIsEquilibriumDetailed(t *testing.T) {
	config := types.SimulationConfig{
		InitialHumans: 2,
		ExperienceDistribution: types.ExperienceDistribution{
			UniversityHire: 100.0,
			MidLevel:       0.0,
			Senior:         0.0,
			Executive:      0.0,
		},
		CostCategoryDistribution: types.CostCategoryDistribution{
			HighCostUS:   100.0,
			LowCostNonUS: 0.0,
		},
		FixedBudget:     300000.0, // Adjusted budget
		RevenueScenario: types.FlatRevenue,
		AILearningSpeeds: types.AILearningSpeed{
			UniversityToMid:   10,
			MidToSenior:       15,
			SeniorToExecutive: 20,
		},
		AttritionConfig: types.AttritionConfig{
			Type:               types.NaturalAttrition,
			NaturalRate:        0.0,
			ForcedAcceleration: 1.0,
		},
		CatastrophicFailureRate: 0.0,
		TimeZoneInefficiency:    0.0,
	}

	controller := NewSimulationController(config, 12345)
	err := controller.Initialize()
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Initially should not be in equilibrium
	isEq, reason := controller.IsEquilibriumDetailed()
	if isEq {
		t.Errorf("Expected not to be in equilibrium initially, but got: %s", reason)
	}

	// Run a few steps
	for i := 0; i < 3; i++ {
		controller.Step()
	}

	// Check equilibrium status
	isEq, reason = controller.IsEquilibriumDetailed()
	t.Logf("Equilibrium status after 3 steps: %v, reason: %s", isEq, reason)
}