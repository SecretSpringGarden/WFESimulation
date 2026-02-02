package analytics

import (
	"bytes"
	"strings"
	"testing"
	"workforce-ai-transition-simulator/internal/types"
)

func TestNewAnalyticsEngine(t *testing.T) {
	engine := NewAnalyticsEngine()
	
	if engine == nil {
		t.Fatal("NewAnalyticsEngine returned nil")
	}
	
	if engine.timeSeries == nil {
		t.Error("timeSeries should be initialized")
	}
	
	if engine.metrics == nil {
		t.Error("metrics should be initialized")
	}
}

func TestRecordTimeStep(t *testing.T) {
	engine := NewAnalyticsEngine()
	
	// Create a sample simulation state
	state := types.SimulationState{
		TimeStep:          1,
		TotalCost:         100000,
		AvailableBudget:   50000,
		TotalProductivity: 10.5,
		RevenueOutput:     25000,
		Workforce: types.WorkforceComposition{
			Humans: struct {
				Total          int
				ByExperience   map[types.ExperienceLevel]int
				ByCostCategory map[types.CostCategory]int
			}{
				Total: 5,
			},
			AIAgents: struct {
				Total        int
				ByExperience map[types.ExperienceLevel]int
			}{
				Total: 3,
			},
			OrchestrationUtilization: 75.0,
		},
		CatastrophicFailures: 1,
	}
	
	// Record the time step
	engine.RecordTimeStep(state)
	
	// Verify the state was recorded
	timeSeries := engine.GetTimeSeries()
	if len(timeSeries) != 1 {
		t.Errorf("Expected 1 time step recorded, got %d", len(timeSeries))
	}
	
	if timeSeries[0].TimeStep != 1 {
		t.Errorf("Expected time step 1, got %d", timeSeries[0].TimeStep)
	}
	
	// Verify metrics were recorded
	metrics := engine.GetMetrics()
	if len(metrics) == 0 {
		t.Error("Expected metrics to be recorded")
	}
	
	// Check specific metrics
	if totalCost, exists := metrics["total_cost"]; !exists || len(totalCost) != 1 || totalCost[0] != 100000 {
		t.Errorf("Expected total_cost metric to be 100000, got %v", totalCost)
	}
}

func TestGenerateReport(t *testing.T) {
	engine := NewAnalyticsEngine()
	
	// Create a sample simulation result
	config := types.SimulationConfig{
		InitialHumans: 10,
		FixedBudget:   200000,
	}
	
	timeSeries := []types.SimulationState{
		{
			TimeStep:          0,
			TotalCost:         150000,
			TotalProductivity: 15.0,
			RevenueOutput:     30000,
			Workforce: types.WorkforceComposition{
				Humans: struct {
					Total          int
					ByExperience   map[types.ExperienceLevel]int
					ByCostCategory map[types.CostCategory]int
				}{
					Total: 10,
				},
				AIAgents: struct {
					Total        int
					ByExperience map[types.ExperienceLevel]int
				}{
					Total: 0,
				},
			},
		},
		{
			TimeStep:          5,
			TotalCost:         180000,
			TotalProductivity: 20.0,
			RevenueOutput:     40000,
			Workforce: types.WorkforceComposition{
				Humans: struct {
					Total          int
					ByExperience   map[types.ExperienceLevel]int
					ByCostCategory map[types.CostCategory]int
				}{
					Total: 8,
				},
				AIAgents: struct {
					Total        int
					ByExperience map[types.ExperienceLevel]int
				}{
					Total: 5,
				},
			},
		},
	}
	
	result := types.SimulationResult{
		Config:            config,
		TimeSeries:        timeSeries,
		EquilibriumState:  timeSeries[1],
		TimeToEquilibrium: 5,
	}
	
	// Generate the report
	report := engine.GenerateReport(result)
	
	// Verify report contents
	if report.InitialParameters.InitialHumans != 10 {
		t.Errorf("Expected initial humans 10, got %d", report.InitialParameters.InitialHumans)
	}
	
	if report.TotalSimulationDuration != 5 {
		t.Errorf("Expected simulation duration 5, got %d", report.TotalSimulationDuration)
	}
	
	if len(report.TimeSeriesData) != 2 {
		t.Errorf("Expected 2 time series entries, got %d", len(report.TimeSeriesData))
	}
	
	if len(report.RevenueTimeSeries) != 2 {
		t.Errorf("Expected 2 revenue entries, got %d", len(report.RevenueTimeSeries))
	}
	
	// Verify summary calculations
	if report.Summary.InitialHumanCount != 10 {
		t.Errorf("Expected initial human count 10, got %d", report.Summary.InitialHumanCount)
	}
	
	if report.Summary.FinalHumanCount != 8 {
		t.Errorf("Expected final human count 8, got %d", report.Summary.FinalHumanCount)
	}
	
	if report.Summary.FinalAIAgentCount != 5 {
		t.Errorf("Expected final AI agent count 5, got %d", report.Summary.FinalAIAgentCount)
	}
}

func TestGenerateReportCSV(t *testing.T) {
	engine := NewAnalyticsEngine()
	
	// Create a simple simulation result
	result := types.SimulationResult{
		TimeSeries: []types.SimulationState{
			{
				TimeStep:          0,
				TotalCost:         100000,
				AvailableBudget:   50000,
				TotalProductivity: 10.0,
				RevenueOutput:     20000,
				Workforce: types.WorkforceComposition{
					Humans: struct {
						Total          int
						ByExperience   map[types.ExperienceLevel]int
						ByCostCategory map[types.CostCategory]int
					}{
						Total: 5,
					},
					AIAgents: struct {
						Total        int
						ByExperience map[types.ExperienceLevel]int
					}{
						Total: 2,
					},
					OrchestrationUtilization: 50.0,
				},
				CatastrophicFailures: 0,
				IsEquilibrium:        false,
			},
		},
	}
	
	// Generate CSV data
	csvData, err := engine.GenerateReportCSV(result)
	if err != nil {
		t.Fatalf("Failed to generate CSV: %v", err)
	}
	
	// Verify CSV structure
	if len(csvData) != 2 { // header + 1 data row
		t.Errorf("Expected 2 CSV rows, got %d", len(csvData))
	}
	
	// Verify header
	expectedHeaders := []string{
		"TimeStep", "HumanCount", "AIAgentCount", "TotalWorkforce",
		"TotalCost", "AvailableBudget", "TotalProductivity", "RevenueOutput",
		"OrchestrationUtilization", "CatastrophicFailures", "IsEquilibrium",
	}
	
	if len(csvData[0]) != len(expectedHeaders) {
		t.Errorf("Expected %d headers, got %d", len(expectedHeaders), len(csvData[0]))
	}
	
	// Verify data row
	dataRow := csvData[1]
	if dataRow[0] != "0" { // TimeStep
		t.Errorf("Expected TimeStep 0, got %s", dataRow[0])
	}
	if dataRow[1] != "5" { // HumanCount
		t.Errorf("Expected HumanCount 5, got %s", dataRow[1])
	}
	if dataRow[2] != "2" { // AIAgentCount
		t.Errorf("Expected AIAgentCount 2, got %s", dataRow[2])
	}
}

func TestWriteReportCSV(t *testing.T) {
	engine := NewAnalyticsEngine()
	
	// Create a simple simulation result
	result := types.SimulationResult{
		TimeSeries: []types.SimulationState{
			{
				TimeStep:          0,
				TotalCost:         100000,
				TotalProductivity: 10.0,
				Workforce: types.WorkforceComposition{
					Humans: struct {
						Total          int
						ByExperience   map[types.ExperienceLevel]int
						ByCostCategory map[types.CostCategory]int
					}{
						Total: 5,
					},
					AIAgents: struct {
						Total        int
						ByExperience map[types.ExperienceLevel]int
					}{
						Total: 2,
					},
				},
			},
		},
	}
	
	// Write to buffer
	var buf bytes.Buffer
	err := engine.WriteReportCSV(result, &buf)
	if err != nil {
		t.Fatalf("Failed to write CSV: %v", err)
	}
	
	// Verify output contains expected data
	output := buf.String()
	if !strings.Contains(output, "TimeStep") {
		t.Error("CSV output should contain TimeStep header")
	}
	if !strings.Contains(output, "0,5,2") { // TimeStep=0, HumanCount=5, AIAgentCount=2
		t.Error("CSV output should contain expected data row")
	}
}

func TestRankParameterImpacts(t *testing.T) {
	engine := NewAnalyticsEngine()
	
	// Create mock sensitivity results
	sensitivityResults := map[string]SensitivityResults{
		"FixedBudget": {
			ParameterName:   "FixedBudget",
			ParameterValues: []float64{100000, 200000, 300000},
			Results: []types.SimulationResult{
				{TimeToEquilibrium: 10},
				{TimeToEquilibrium: 5},
				{TimeToEquilibrium: 3},
			},
		},
		"InitialHumans": {
			ParameterName:   "InitialHumans",
			ParameterValues: []float64{5, 10, 15},
			Results: []types.SimulationResult{
				{TimeToEquilibrium: 8},
				{TimeToEquilibrium: 8},
				{TimeToEquilibrium: 9},
			},
		},
	}
	
	// Rank parameter impacts
	impacts := engine.RankParameterImpacts(sensitivityResults)
	
	// Verify we got results for both parameters
	if len(impacts) != 2 {
		t.Errorf("Expected 2 parameter impacts, got %d", len(impacts))
	}
	
	// Verify impacts are sorted (highest impact first)
	// FixedBudget should have higher variance (10,5,3) vs InitialHumans (8,8,9)
	if impacts[0].ParameterName != "FixedBudget" {
		t.Errorf("Expected FixedBudget to have highest impact, got %s", impacts[0].ParameterName)
	}
	
	// Verify impact values are calculated
	if impacts[0].TimeToEquilibriumImpact <= 0 {
		t.Error("Expected positive time to equilibrium impact")
	}
}

func TestCalculateVariance(t *testing.T) {
	engine := NewAnalyticsEngine()
	
	// Test with known values
	values := []float64{1, 2, 3, 4, 5}
	variance := engine.calculateVariance(values)
	
	// Expected variance for [1,2,3,4,5] is 2.5
	expectedVariance := 2.5
	if variance != expectedVariance {
		t.Errorf("Expected variance %.2f, got %.2f", expectedVariance, variance)
	}
	
	// Test with single value
	singleValue := []float64{5}
	variance = engine.calculateVariance(singleValue)
	if variance != 0 {
		t.Errorf("Expected variance 0 for single value, got %.2f", variance)
	}
	
	// Test with empty slice
	emptyValues := []float64{}
	variance = engine.calculateVariance(emptyValues)
	if variance != 0 {
		t.Errorf("Expected variance 0 for empty slice, got %.2f", variance)
	}
}