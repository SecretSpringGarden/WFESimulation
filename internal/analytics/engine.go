package analytics

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"sort"
	"sync"
	"workforce-ai-transition-simulator/internal/controller"
	"workforce-ai-transition-simulator/internal/types"
)

// AnalyticsEngine stores time-series data and metrics for simulation analysis
// Requirement 10.7: Store time-series data and metrics
type AnalyticsEngine struct {
	// Time-series data storage
	timeSeries []types.SimulationState
	
	// Metrics storage
	metrics map[string][]float64
	
	// Mutex for thread-safe operations during parallel sensitivity analysis
	mu sync.RWMutex
}

// NewAnalyticsEngine creates a new AnalyticsEngine instance
func NewAnalyticsEngine() *AnalyticsEngine {
	return &AnalyticsEngine{
		timeSeries: make([]types.SimulationState, 0),
		metrics:    make(map[string][]float64),
	}
}

// Reset clears all stored data and metrics
func (ae *AnalyticsEngine) Reset() {
	ae.mu.Lock()
	defer ae.mu.Unlock()
	
	ae.timeSeries = make([]types.SimulationState, 0)
	ae.metrics = make(map[string][]float64)
}

// GetTimeSeries returns a copy of the stored time series data
func (ae *AnalyticsEngine) GetTimeSeries() []types.SimulationState {
	ae.mu.RLock()
	defer ae.mu.RUnlock()
	
	// Return a copy to prevent external modification
	result := make([]types.SimulationState, len(ae.timeSeries))
	copy(result, ae.timeSeries)
	return result
}

// GetMetrics returns a copy of the stored metrics
func (ae *AnalyticsEngine) GetMetrics() map[string][]float64 {
	ae.mu.RLock()
	defer ae.mu.RUnlock()
	
	// Return a deep copy to prevent external modification
	result := make(map[string][]float64)
	for key, values := range ae.metrics {
		result[key] = make([]float64, len(values))
		copy(result[key], values)
	}
	return result
}

// SensitivityResults represents the results of a sensitivity analysis
type SensitivityResults struct {
	ParameterName                    string
	ParameterValues                  []float64
	Results                         []types.SimulationResult
	TimeToEquilibriumByValue        map[float64]int
	EquilibriumCompositionByValue   map[float64]types.WorkforceComposition
}

// ParameterRanges defines the ranges for sensitivity analysis parameters
type ParameterRanges struct {
	FixedBudget             []float64
	InitialHumans           []int
	CatastrophicFailureRate []float64
	TimeZoneInefficiency    []float64
	NaturalAttritionRate    []float64
	ForcedAcceleration      []float64
	UniversityToMid         []int
	MidToSenior            []int
	SeniorToExecutive      []int
}

// ParameterImpact represents the impact of a parameter on simulation outcomes
type ParameterImpact struct {
	ParameterName           string
	TimeToEquilibriumImpact float64 // variance in time to equilibrium
	CompositionImpact       float64 // variance in final composition
}

// Report represents a comprehensive simulation report
type Report struct {
	InitialParameters       types.SimulationConfig
	TimeSeriesData         []types.SimulationState
	RevenueTimeSeries      []float64
	EquilibriumDetails     types.SimulationState
	TotalSimulationDuration int
	Summary                ReportSummary
}

// ReportSummary provides key metrics and insights from the simulation
type ReportSummary struct {
	InitialWorkforceSize    int
	FinalWorkforceSize      int
	InitialHumanCount       int
	FinalHumanCount         int
	InitialAIAgentCount     int
	FinalAIAgentCount       int
	TotalRevenueGenerated   float64
	AverageProductivity     float64
	CostEfficiencyRatio     float64 // final productivity / final cost
}

// SensitivityReport represents a sensitivity analysis report
type SensitivityReport struct {
	ParameterRankings       []ParameterImpact
	DetailedResults         map[string]SensitivityResults
	Summary                 SensitivitySummary
}

// SensitivitySummary provides key insights from sensitivity analysis
type SensitivitySummary struct {
	MostImpactfulParameter     string
	LeastImpactfulParameter    string
	AverageTimeToEquilibrium   float64
	TimeToEquilibriumVariance  float64
	OptimalParameterValues     map[string]float64
}

// RecordTimeStep captures and stores simulation state at each time step
// Requirement 10.7: Capture and store simulation state at each time step
func (ae *AnalyticsEngine) RecordTimeStep(state types.SimulationState) {
	ae.mu.Lock()
	defer ae.mu.Unlock()
	
	// Store the simulation state
	ae.timeSeries = append(ae.timeSeries, state)
	
	// Extract and store key metrics for analysis
	ae.recordMetric("total_cost", state.TotalCost)
	ae.recordMetric("available_budget", state.AvailableBudget)
	ae.recordMetric("total_productivity", state.TotalProductivity)
	ae.recordMetric("revenue_output", state.RevenueOutput)
	ae.recordMetric("human_count", float64(state.Workforce.Humans.Total))
	ae.recordMetric("ai_agent_count", float64(state.Workforce.AIAgents.Total))
	ae.recordMetric("orchestration_utilization", state.Workforce.OrchestrationUtilization)
	ae.recordMetric("catastrophic_failures", float64(state.CatastrophicFailures))
	
	// Calculate and store derived metrics
	totalWorkforce := float64(state.Workforce.Humans.Total + state.Workforce.AIAgents.Total)
	ae.recordMetric("total_workforce", totalWorkforce)
	
	// Cost efficiency ratio (productivity per unit cost)
	if state.TotalCost > 0 {
		costEfficiency := state.TotalProductivity / state.TotalCost
		ae.recordMetric("cost_efficiency", costEfficiency)
	}
	
	// AI agent ratio (percentage of workforce that is AI)
	if totalWorkforce > 0 {
		aiRatio := float64(state.Workforce.AIAgents.Total) / totalWorkforce * 100.0
		ae.recordMetric("ai_ratio", aiRatio)
	}
}

// recordMetric is a helper method to store individual metrics
func (ae *AnalyticsEngine) recordMetric(name string, value float64) {
	if ae.metrics[name] == nil {
		ae.metrics[name] = make([]float64, 0)
	}
	ae.metrics[name] = append(ae.metrics[name], value)
}

// RecordSimulationResult records the complete result of a simulation run
func (ae *AnalyticsEngine) RecordSimulationResult(result types.SimulationResult) {
	ae.mu.Lock()
	defer ae.mu.Unlock()
	
	// Store all time series data from the simulation
	ae.timeSeries = make([]types.SimulationState, len(result.TimeSeries))
	copy(ae.timeSeries, result.TimeSeries)
	
	// Clear and rebuild metrics from the time series
	ae.metrics = make(map[string][]float64)
	
	for _, state := range result.TimeSeries {
		ae.recordMetric("total_cost", state.TotalCost)
		ae.recordMetric("available_budget", state.AvailableBudget)
		ae.recordMetric("total_productivity", state.TotalProductivity)
		ae.recordMetric("revenue_output", state.RevenueOutput)
		ae.recordMetric("human_count", float64(state.Workforce.Humans.Total))
		ae.recordMetric("ai_agent_count", float64(state.Workforce.AIAgents.Total))
		ae.recordMetric("orchestration_utilization", state.Workforce.OrchestrationUtilization)
		ae.recordMetric("catastrophic_failures", float64(state.CatastrophicFailures))
		
		// Derived metrics
		totalWorkforce := float64(state.Workforce.Humans.Total + state.Workforce.AIAgents.Total)
		ae.recordMetric("total_workforce", totalWorkforce)
		
		if state.TotalCost > 0 {
			costEfficiency := state.TotalProductivity / state.TotalCost
			ae.recordMetric("cost_efficiency", costEfficiency)
		}
		
		if totalWorkforce > 0 {
			aiRatio := float64(state.Workforce.AIAgents.Total) / totalWorkforce * 100.0
			ae.recordMetric("ai_ratio", aiRatio)
		}
	}
}
// RunSensitivityAnalysis executes multiple simulations with parameter variations
// Requirements 11.1, 11.2: Execute multiple simulations varying one parameter at a time
// Uses Go goroutines for parallel execution
func (ae *AnalyticsEngine) RunSensitivityAnalysis(baseConfig types.SimulationConfig, paramRanges ParameterRanges, maxTimeSteps int, seed int64) (map[string]SensitivityResults, error) {
	results := make(map[string]SensitivityResults)
	
	// Channel for collecting results from goroutines
	type paramResult struct {
		paramName string
		result    SensitivityResults
		err       error
	}
	resultChan := make(chan paramResult, 10) // Buffer for up to 10 parameters
	
	// WaitGroup to wait for all goroutines to complete
	var wg sync.WaitGroup
	
	// Run sensitivity analysis for FixedBudget parameter
	if len(paramRanges.FixedBudget) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := ae.runParameterSensitivity("FixedBudget", baseConfig, paramRanges.FixedBudget, maxTimeSteps, seed, func(config *types.SimulationConfig, value float64) {
				config.FixedBudget = value
			})
			resultChan <- paramResult{"FixedBudget", result, err}
		}()
	}
	
	// Run sensitivity analysis for InitialHumans parameter
	if len(paramRanges.InitialHumans) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			intValues := make([]float64, len(paramRanges.InitialHumans))
			for i, v := range paramRanges.InitialHumans {
				intValues[i] = float64(v)
			}
			result, err := ae.runParameterSensitivity("InitialHumans", baseConfig, intValues, maxTimeSteps, seed+1, func(config *types.SimulationConfig, value float64) {
				config.InitialHumans = int(value)
			})
			resultChan <- paramResult{"InitialHumans", result, err}
		}()
	}
	
	// Run sensitivity analysis for CatastrophicFailureRate parameter
	if len(paramRanges.CatastrophicFailureRate) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := ae.runParameterSensitivity("CatastrophicFailureRate", baseConfig, paramRanges.CatastrophicFailureRate, maxTimeSteps, seed+2, func(config *types.SimulationConfig, value float64) {
				config.CatastrophicFailureRate = value
			})
			resultChan <- paramResult{"CatastrophicFailureRate", result, err}
		}()
	}
	
	// Run sensitivity analysis for TimeZoneInefficiency parameter
	if len(paramRanges.TimeZoneInefficiency) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := ae.runParameterSensitivity("TimeZoneInefficiency", baseConfig, paramRanges.TimeZoneInefficiency, maxTimeSteps, seed+3, func(config *types.SimulationConfig, value float64) {
				config.TimeZoneInefficiency = value
			})
			resultChan <- paramResult{"TimeZoneInefficiency", result, err}
		}()
	}
	
	// Run sensitivity analysis for NaturalAttritionRate parameter
	if len(paramRanges.NaturalAttritionRate) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := ae.runParameterSensitivity("NaturalAttritionRate", baseConfig, paramRanges.NaturalAttritionRate, maxTimeSteps, seed+4, func(config *types.SimulationConfig, value float64) {
				config.AttritionConfig.NaturalRate = value
			})
			resultChan <- paramResult{"NaturalAttritionRate", result, err}
		}()
	}
	
	// Run sensitivity analysis for ForcedAcceleration parameter
	if len(paramRanges.ForcedAcceleration) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := ae.runParameterSensitivity("ForcedAcceleration", baseConfig, paramRanges.ForcedAcceleration, maxTimeSteps, seed+5, func(config *types.SimulationConfig, value float64) {
				config.AttritionConfig.ForcedAcceleration = value
			})
			resultChan <- paramResult{"ForcedAcceleration", result, err}
		}()
	}
	
	// Run sensitivity analysis for AI learning speed parameters
	if len(paramRanges.UniversityToMid) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			intValues := make([]float64, len(paramRanges.UniversityToMid))
			for i, v := range paramRanges.UniversityToMid {
				intValues[i] = float64(v)
			}
			result, err := ae.runParameterSensitivity("UniversityToMid", baseConfig, intValues, maxTimeSteps, seed+6, func(config *types.SimulationConfig, value float64) {
				config.AILearningSpeeds.UniversityToMid = int(value)
			})
			resultChan <- paramResult{"UniversityToMid", result, err}
		}()
	}
	
	if len(paramRanges.MidToSenior) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			intValues := make([]float64, len(paramRanges.MidToSenior))
			for i, v := range paramRanges.MidToSenior {
				intValues[i] = float64(v)
			}
			result, err := ae.runParameterSensitivity("MidToSenior", baseConfig, intValues, maxTimeSteps, seed+7, func(config *types.SimulationConfig, value float64) {
				config.AILearningSpeeds.MidToSenior = int(value)
			})
			resultChan <- paramResult{"MidToSenior", result, err}
		}()
	}
	
	if len(paramRanges.SeniorToExecutive) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			intValues := make([]float64, len(paramRanges.SeniorToExecutive))
			for i, v := range paramRanges.SeniorToExecutive {
				intValues[i] = float64(v)
			}
			result, err := ae.runParameterSensitivity("SeniorToExecutive", baseConfig, intValues, maxTimeSteps, seed+8, func(config *types.SimulationConfig, value float64) {
				config.AILearningSpeeds.SeniorToExecutive = int(value)
			})
			resultChan <- paramResult{"SeniorToExecutive", result, err}
		}()
	}
	
	// Close the result channel when all goroutines are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	
	// Collect results from all goroutines
	for result := range resultChan {
		if result.err != nil {
			return nil, fmt.Errorf("sensitivity analysis failed for parameter %s: %w", result.paramName, result.err)
		}
		results[result.paramName] = result.result
	}
	
	return results, nil
}

// runParameterSensitivity runs sensitivity analysis for a single parameter
func (ae *AnalyticsEngine) runParameterSensitivity(paramName string, baseConfig types.SimulationConfig, values []float64, maxTimeSteps int, seed int64, setter func(*types.SimulationConfig, float64)) (SensitivityResults, error) {
	results := make([]types.SimulationResult, len(values))
	timeToEquilibrium := make(map[float64]int)
	equilibriumComposition := make(map[float64]types.WorkforceComposition)
	
	// Run simulation for each parameter value
	for i, value := range values {
		// Create a copy of the base configuration
		config := baseConfig
		
		// Apply the parameter value using the setter function
		setter(&config, value)
		
		// Create a new simulation controller with unique seed
		simController := controller.NewSimulationController(config, seed+int64(i))
		
		// Run the simulation
		result, err := simController.RunUntilEquilibrium(maxTimeSteps)
		if err != nil {
			return SensitivityResults{}, fmt.Errorf("simulation failed for %s=%f: %w", paramName, value, err)
		}
		
		// Store the results
		results[i] = result
		timeToEquilibrium[value] = result.TimeToEquilibrium
		equilibriumComposition[value] = result.EquilibriumState.Workforce
	}
	
	return SensitivityResults{
		ParameterName:                   paramName,
		ParameterValues:                 values,
		Results:                        results,
		TimeToEquilibriumByValue:       timeToEquilibrium,
		EquilibriumCompositionByValue:  equilibriumComposition,
	}, nil
}
// RankParameterImpacts calculates and ranks parameter impacts on equilibrium time and composition
// Requirements 11.5, 11.6: Rank parameters by their impact on time to equilibrium and final workforce composition
func (ae *AnalyticsEngine) RankParameterImpacts(sensitivityResults map[string]SensitivityResults) []ParameterImpact {
	impacts := make([]ParameterImpact, 0, len(sensitivityResults))
	
	for paramName, results := range sensitivityResults {
		// Calculate impact on time to equilibrium
		timeToEquilibriumImpact := ae.calculateVariance(ae.extractTimeToEquilibrium(results))
		
		// Calculate impact on workforce composition
		compositionImpact := ae.calculateCompositionVariance(results)
		
		impacts = append(impacts, ParameterImpact{
			ParameterName:           paramName,
			TimeToEquilibriumImpact: timeToEquilibriumImpact,
			CompositionImpact:       compositionImpact,
		})
	}
	
	// Sort by combined impact (time to equilibrium impact + composition impact)
	sort.Slice(impacts, func(i, j int) bool {
		impactI := impacts[i].TimeToEquilibriumImpact + impacts[i].CompositionImpact
		impactJ := impacts[j].TimeToEquilibriumImpact + impacts[j].CompositionImpact
		return impactI > impactJ // Sort in descending order (highest impact first)
	})
	
	return impacts
}

// extractTimeToEquilibrium extracts time to equilibrium values from sensitivity results
func (ae *AnalyticsEngine) extractTimeToEquilibrium(results SensitivityResults) []float64 {
	values := make([]float64, len(results.Results))
	for i, result := range results.Results {
		values[i] = float64(result.TimeToEquilibrium)
	}
	return values
}

// calculateVariance calculates the variance of a slice of float64 values
func (ae *AnalyticsEngine) calculateVariance(values []float64) float64 {
	if len(values) <= 1 {
		return 0.0
	}
	
	// Calculate mean
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(len(values))
	
	// Calculate variance
	sumSquaredDiffs := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquaredDiffs += diff * diff
	}
	
	return sumSquaredDiffs / float64(len(values)-1)
}

// calculateCompositionVariance calculates the variance in workforce composition across parameter values
func (ae *AnalyticsEngine) calculateCompositionVariance(results SensitivityResults) float64 {
	if len(results.Results) <= 1 {
		return 0.0
	}
	
	// Extract composition metrics for variance calculation
	humanCounts := make([]float64, len(results.Results))
	aiCounts := make([]float64, len(results.Results))
	orchestrationUtils := make([]float64, len(results.Results))
	
	for i, result := range results.Results {
		composition := result.EquilibriumState.Workforce
		humanCounts[i] = float64(composition.Humans.Total)
		aiCounts[i] = float64(composition.AIAgents.Total)
		orchestrationUtils[i] = composition.OrchestrationUtilization
	}
	
	// Calculate variance for each composition metric
	humanVariance := ae.calculateVariance(humanCounts)
	aiVariance := ae.calculateVariance(aiCounts)
	orchestrationVariance := ae.calculateVariance(orchestrationUtils)
	
	// Return combined variance (weighted average)
	return (humanVariance + aiVariance + orchestrationVariance/100.0) / 3.0
}

// RankParametersByTimeImpact ranks parameters specifically by their impact on time to equilibrium
func (ae *AnalyticsEngine) RankParametersByTimeImpact(sensitivityResults map[string]SensitivityResults) []ParameterImpact {
	impacts := make([]ParameterImpact, 0, len(sensitivityResults))
	
	for paramName, results := range sensitivityResults {
		timeToEquilibriumImpact := ae.calculateVariance(ae.extractTimeToEquilibrium(results))
		
		impacts = append(impacts, ParameterImpact{
			ParameterName:           paramName,
			TimeToEquilibriumImpact: timeToEquilibriumImpact,
			CompositionImpact:       0, // Not used for this ranking
		})
	}
	
	// Sort by time to equilibrium impact only
	sort.Slice(impacts, func(i, j int) bool {
		return impacts[i].TimeToEquilibriumImpact > impacts[j].TimeToEquilibriumImpact
	})
	
	return impacts
}

// RankParametersByCompositionImpact ranks parameters specifically by their impact on final workforce composition
func (ae *AnalyticsEngine) RankParametersByCompositionImpact(sensitivityResults map[string]SensitivityResults) []ParameterImpact {
	impacts := make([]ParameterImpact, 0, len(sensitivityResults))
	
	for paramName, results := range sensitivityResults {
		compositionImpact := ae.calculateCompositionVariance(results)
		
		impacts = append(impacts, ParameterImpact{
			ParameterName:           paramName,
			TimeToEquilibriumImpact: 0, // Not used for this ranking
			CompositionImpact:       compositionImpact,
		})
	}
	
	// Sort by composition impact only
	sort.Slice(impacts, func(i, j int) bool {
		return impacts[i].CompositionImpact > impacts[j].CompositionImpact
	})
	
	return impacts
}

// CalculateSensitivitySummary calculates summary statistics from sensitivity analysis results
func (ae *AnalyticsEngine) CalculateSensitivitySummary(sensitivityResults map[string]SensitivityResults) SensitivitySummary {
	if len(sensitivityResults) == 0 {
		return SensitivitySummary{}
	}
	
	// Rank parameters by combined impact
	impacts := ae.RankParameterImpacts(sensitivityResults)
	
	var mostImpactful, leastImpactful string
	if len(impacts) > 0 {
		mostImpactful = impacts[0].ParameterName
		leastImpactful = impacts[len(impacts)-1].ParameterName
	}
	
	// Calculate average time to equilibrium across all parameter variations
	totalTime := 0.0
	totalCount := 0
	timeValues := make([]float64, 0)
	
	for _, results := range sensitivityResults {
		for _, result := range results.Results {
			totalTime += float64(result.TimeToEquilibrium)
			totalCount++
			timeValues = append(timeValues, float64(result.TimeToEquilibrium))
		}
	}
	
	averageTime := 0.0
	if totalCount > 0 {
		averageTime = totalTime / float64(totalCount)
	}
	
	// Calculate variance in time to equilibrium
	timeVariance := ae.calculateVariance(timeValues)
	
	// Find optimal parameter values (those that minimize time to equilibrium)
	optimalValues := ae.findOptimalParameterValues(sensitivityResults)
	
	return SensitivitySummary{
		MostImpactfulParameter:     mostImpactful,
		LeastImpactfulParameter:    leastImpactful,
		AverageTimeToEquilibrium:   averageTime,
		TimeToEquilibriumVariance:  timeVariance,
		OptimalParameterValues:     optimalValues,
	}
}

// findOptimalParameterValues finds parameter values that minimize time to equilibrium
func (ae *AnalyticsEngine) findOptimalParameterValues(sensitivityResults map[string]SensitivityResults) map[string]float64 {
	optimal := make(map[string]float64)
	
	for paramName, results := range sensitivityResults {
		minTime := math.Inf(1)
		optimalValue := 0.0
		
		for i, result := range results.Results {
			if float64(result.TimeToEquilibrium) < minTime {
				minTime = float64(result.TimeToEquilibrium)
				optimalValue = results.ParameterValues[i]
			}
		}
		
		optimal[paramName] = optimalValue
	}
	
	return optimal
}
// GenerateReport creates a comprehensive simulation report with all required data
// Requirements 12.1, 12.2, 12.3, 12.4, 12.5: Generate report containing initial parameters,
// time-series data, revenue output, equilibrium state details, and total simulation duration
func (ae *AnalyticsEngine) GenerateReport(result types.SimulationResult) Report {
	// Extract revenue time series from simulation states
	revenueTimeSeries := make([]float64, len(result.TimeSeries))
	for i, state := range result.TimeSeries {
		revenueTimeSeries[i] = state.RevenueOutput
	}
	
	// Calculate summary statistics
	summary := ae.calculateReportSummary(result)
	
	return Report{
		InitialParameters:       result.Config,
		TimeSeriesData:         result.TimeSeries,
		RevenueTimeSeries:      revenueTimeSeries,
		EquilibriumDetails:     result.EquilibriumState,
		TotalSimulationDuration: result.TimeToEquilibrium,
		Summary:                summary,
	}
}

// calculateReportSummary calculates key metrics and insights from the simulation result
func (ae *AnalyticsEngine) calculateReportSummary(result types.SimulationResult) ReportSummary {
	if len(result.TimeSeries) == 0 {
		return ReportSummary{}
	}
	
	initialState := result.TimeSeries[0]
	finalState := result.EquilibriumState
	
	// Calculate total revenue generated throughout the simulation
	totalRevenue := 0.0
	for _, state := range result.TimeSeries {
		totalRevenue += state.RevenueOutput
	}
	
	// Calculate average productivity across the simulation
	totalProductivity := 0.0
	for _, state := range result.TimeSeries {
		totalProductivity += state.TotalProductivity
	}
	averageProductivity := totalProductivity / float64(len(result.TimeSeries))
	
	// Calculate cost efficiency ratio (final productivity / final cost)
	costEfficiencyRatio := 0.0
	if finalState.TotalCost > 0 {
		costEfficiencyRatio = finalState.TotalProductivity / finalState.TotalCost
	}
	
	return ReportSummary{
		InitialWorkforceSize:    initialState.Workforce.Humans.Total + initialState.Workforce.AIAgents.Total,
		FinalWorkforceSize:      finalState.Workforce.Humans.Total + finalState.Workforce.AIAgents.Total,
		InitialHumanCount:       initialState.Workforce.Humans.Total,
		FinalHumanCount:         finalState.Workforce.Humans.Total,
		InitialAIAgentCount:     initialState.Workforce.AIAgents.Total,
		FinalAIAgentCount:       finalState.Workforce.AIAgents.Total,
		TotalRevenueGenerated:   totalRevenue,
		AverageProductivity:     averageProductivity,
		CostEfficiencyRatio:     costEfficiencyRatio,
	}
}

// GenerateReportJSON generates a JSON representation of the simulation report
func (ae *AnalyticsEngine) GenerateReportJSON(result types.SimulationResult) ([]byte, error) {
	report := ae.GenerateReport(result)
	return json.MarshalIndent(report, "", "  ")
}

// WriteReportJSON writes the simulation report to a JSON file
func (ae *AnalyticsEngine) WriteReportJSON(result types.SimulationResult, writer io.Writer) error {
	jsonData, err := ae.GenerateReportJSON(result)
	if err != nil {
		return fmt.Errorf("failed to generate JSON report: %w", err)
	}
	
	_, err = writer.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to write JSON report: %w", err)
	}
	
	return nil
}

// GenerateReportCSV generates a CSV representation of the time series data
func (ae *AnalyticsEngine) GenerateReportCSV(result types.SimulationResult) ([][]string, error) {
	if len(result.TimeSeries) == 0 {
		return nil, fmt.Errorf("no time series data available")
	}
	
	// Create CSV header
	header := []string{
		"TimeStep",
		"HumanCount",
		"AIAgentCount",
		"TotalWorkforce",
		"TotalCost",
		"AvailableBudget",
		"TotalProductivity",
		"RevenueOutput",
		"OrchestrationUtilization",
		"CatastrophicFailures",
		"IsEquilibrium",
	}
	
	// Create CSV data
	data := make([][]string, len(result.TimeSeries)+1)
	data[0] = header
	
	for i, state := range result.TimeSeries {
		row := []string{
			fmt.Sprintf("%d", state.TimeStep),
			fmt.Sprintf("%d", state.Workforce.Humans.Total),
			fmt.Sprintf("%d", state.Workforce.AIAgents.Total),
			fmt.Sprintf("%d", state.Workforce.Humans.Total+state.Workforce.AIAgents.Total),
			fmt.Sprintf("%.2f", state.TotalCost),
			fmt.Sprintf("%.2f", state.AvailableBudget),
			fmt.Sprintf("%.2f", state.TotalProductivity),
			fmt.Sprintf("%.2f", state.RevenueOutput),
			fmt.Sprintf("%.2f", state.Workforce.OrchestrationUtilization),
			fmt.Sprintf("%d", state.CatastrophicFailures),
			fmt.Sprintf("%t", state.IsEquilibrium),
		}
		data[i+1] = row
	}
	
	return data, nil
}

// WriteReportCSV writes the simulation report to a CSV file
func (ae *AnalyticsEngine) WriteReportCSV(result types.SimulationResult, writer io.Writer) error {
	csvData, err := ae.GenerateReportCSV(result)
	if err != nil {
		return fmt.Errorf("failed to generate CSV report: %w", err)
	}
	
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()
	
	for _, row := range csvData {
		if err := csvWriter.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}
	
	return nil
}
// GenerateSensitivityReport creates a sensitivity analysis report with parameter rankings
// Requirements 12.6, 12.7: Create sensitivity analysis report with parameter rankings in CSV/JSON format
func (ae *AnalyticsEngine) GenerateSensitivityReport(sensitivityResults map[string]SensitivityResults) SensitivityReport {
	// Calculate parameter rankings
	parameterRankings := ae.RankParameterImpacts(sensitivityResults)
	
	// Calculate summary statistics
	summary := ae.CalculateSensitivitySummary(sensitivityResults)
	
	return SensitivityReport{
		ParameterRankings: parameterRankings,
		DetailedResults:   sensitivityResults,
		Summary:          summary,
	}
}

// GenerateSensitivityReportJSON generates a JSON representation of the sensitivity analysis report
func (ae *AnalyticsEngine) GenerateSensitivityReportJSON(sensitivityResults map[string]SensitivityResults) ([]byte, error) {
	report := ae.GenerateSensitivityReport(sensitivityResults)
	return json.MarshalIndent(report, "", "  ")
}

// WriteSensitivityReportJSON writes the sensitivity analysis report to a JSON file
func (ae *AnalyticsEngine) WriteSensitivityReportJSON(sensitivityResults map[string]SensitivityResults, writer io.Writer) error {
	jsonData, err := ae.GenerateSensitivityReportJSON(sensitivityResults)
	if err != nil {
		return fmt.Errorf("failed to generate JSON sensitivity report: %w", err)
	}
	
	_, err = writer.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to write JSON sensitivity report: %w", err)
	}
	
	return nil
}

// GenerateSensitivityReportCSV generates a CSV representation of the parameter rankings
func (ae *AnalyticsEngine) GenerateSensitivityReportCSV(sensitivityResults map[string]SensitivityResults) ([][]string, error) {
	if len(sensitivityResults) == 0 {
		return nil, fmt.Errorf("no sensitivity results available")
	}
	
	// Generate parameter rankings
	rankings := ae.RankParameterImpacts(sensitivityResults)
	
	// Create CSV header for parameter rankings
	header := []string{
		"Rank",
		"ParameterName",
		"TimeToEquilibriumImpact",
		"CompositionImpact",
		"CombinedImpact",
	}
	
	// Create CSV data
	data := make([][]string, len(rankings)+1)
	data[0] = header
	
	for i, impact := range rankings {
		combinedImpact := impact.TimeToEquilibriumImpact + impact.CompositionImpact
		row := []string{
			fmt.Sprintf("%d", i+1),
			impact.ParameterName,
			fmt.Sprintf("%.4f", impact.TimeToEquilibriumImpact),
			fmt.Sprintf("%.4f", impact.CompositionImpact),
			fmt.Sprintf("%.4f", combinedImpact),
		}
		data[i+1] = row
	}
	
	return data, nil
}

// WriteSensitivityReportCSV writes the sensitivity analysis report to a CSV file
func (ae *AnalyticsEngine) WriteSensitivityReportCSV(sensitivityResults map[string]SensitivityResults, writer io.Writer) error {
	csvData, err := ae.GenerateSensitivityReportCSV(sensitivityResults)
	if err != nil {
		return fmt.Errorf("failed to generate CSV sensitivity report: %w", err)
	}
	
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()
	
	for _, row := range csvData {
		if err := csvWriter.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}
	
	return nil
}

// GenerateDetailedSensitivityCSV generates a detailed CSV with all parameter variations and their results
func (ae *AnalyticsEngine) GenerateDetailedSensitivityCSV(sensitivityResults map[string]SensitivityResults) ([][]string, error) {
	if len(sensitivityResults) == 0 {
		return nil, fmt.Errorf("no sensitivity results available")
	}
	
	// Create CSV header
	header := []string{
		"ParameterName",
		"ParameterValue",
		"TimeToEquilibrium",
		"FinalHumanCount",
		"FinalAIAgentCount",
		"FinalTotalCost",
		"FinalProductivity",
		"FinalRevenue",
		"OrchestrationUtilization",
		"CatastrophicFailures",
	}
	
	// Calculate total rows needed
	totalRows := 1 // header
	for _, results := range sensitivityResults {
		totalRows += len(results.Results)
	}
	
	// Create CSV data
	data := make([][]string, totalRows)
	data[0] = header
	
	rowIndex := 1
	for paramName, results := range sensitivityResults {
		for i, result := range results.Results {
			paramValue := results.ParameterValues[i]
			equilibrium := result.EquilibriumState
			
			row := []string{
				paramName,
				fmt.Sprintf("%.4f", paramValue),
				fmt.Sprintf("%d", result.TimeToEquilibrium),
				fmt.Sprintf("%d", equilibrium.Workforce.Humans.Total),
				fmt.Sprintf("%d", equilibrium.Workforce.AIAgents.Total),
				fmt.Sprintf("%.2f", equilibrium.TotalCost),
				fmt.Sprintf("%.2f", equilibrium.TotalProductivity),
				fmt.Sprintf("%.2f", equilibrium.RevenueOutput),
				fmt.Sprintf("%.2f", equilibrium.Workforce.OrchestrationUtilization),
				fmt.Sprintf("%d", result.TotalCatastrophicFailures),
			}
			data[rowIndex] = row
			rowIndex++
		}
	}
	
	return data, nil
}

// WriteDetailedSensitivityCSV writes the detailed sensitivity analysis results to a CSV file
func (ae *AnalyticsEngine) WriteDetailedSensitivityCSV(sensitivityResults map[string]SensitivityResults, writer io.Writer) error {
	csvData, err := ae.GenerateDetailedSensitivityCSV(sensitivityResults)
	if err != nil {
		return fmt.Errorf("failed to generate detailed CSV sensitivity report: %w", err)
	}
	
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()
	
	for _, row := range csvData {
		if err := csvWriter.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}
	
	return nil
}