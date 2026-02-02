# Workforce AI Transition Simulator

A discrete-event simulation system that models the evolution of an engineering workforce from purely human to AI-augmented composition under budget constraints. The simulator tracks the transition from a purely human workforce to an AI-augmented equilibrium state, measuring productivity through revenue output.

## Features

- **Agent-Based Modeling**: Individual workers (human and AI) are modeled as entities with distinct characteristics
- **Budget-Driven Optimization**: All workforce decisions are constrained by a fixed budget
- **AI Learning Progression**: AI agents learn and improve over time, progressing through experience levels
- **Multiple Attrition Scenarios**: Support for natural attrition, hiring freezes, and reduction in force
- **Catastrophic Failure Modeling**: Probabilistic failure events that test workforce resilience
- **Sensitivity Analysis**: Run multiple simulations with parameter variations to identify key factors
- **Comprehensive Reporting**: Generate detailed reports in CSV and JSON formats

## Prerequisites

- Go 1.21 or higher
- Git

## Installation

1. Install Go from https://golang.org/dl/

2. Clone this repository and navigate to the project directory

3. Install dependencies:
```bash
go mod download
```

4. Build the simulator:
```bash
go build -o simulator ./cmd/simulator
```

## Quick Start

### Running a Single Simulation

```bash
# Run with default configuration
./simulator

# Run with a specific configuration file
./simulator -config examples/small_team_natural_attrition.json

# Run with YAML configuration
./simulator -config examples/medium_team_fast_learning.yaml
```

### Running Sensitivity Analysis

```bash
# Run sensitivity analysis on all parameters
./simulator -sensitivity -config examples/small_team_natural_attrition.json

# Run sensitivity analysis with custom output directory
./simulator -sensitivity -config examples/large_team_global_distributed.yaml -output results/
```

## Configuration

The simulator accepts configuration in JSON or YAML format. Configuration files define all simulation parameters:

### Configuration Parameters

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `InitialHumans` | int | Starting number of human workers | `10` |
| `ExperienceDistribution` | object | Percentage distribution across experience levels | See examples |
| `CostCategoryDistribution` | object | Percentage distribution across cost categories | See examples |
| `FixedBudget` | float | Total fixed monetary allocation for workforce | `1800000.0` |
| `RevenueScenario` | int | Revenue growth pattern (0=Flat, 1=Explosive) | `0` |
| `AILearningSpeeds` | object | Time steps required for AI level progression | See examples |
| `AttritionConfig` | object | Human attrition behavior configuration | See examples |
| `CatastrophicFailureRate` | float | Probability of failure events per time step | `0.015` |
| `TimeZoneInefficiency` | float | Productivity penalty for distributed workers | `0.15` |

### Experience Levels

- **University_Hire**: Entry-level workers
- **Mid_Level**: Experienced workers
- **Senior**: Senior-level workers
- **Executive**: Leadership-level workers

### Cost Categories

- **High_Cost_US**: US-based workers with higher costs
- **Low_Cost_Non_US**: Non-US workers with lower costs and potential time zone inefficiency

### Revenue Scenarios

- **Flat_Revenue** (0): Constant revenue targets over time
- **Explosive_Growth** (1): Exponentially increasing revenue targets

### Attrition Types

- **Natural_Attrition** (0): Probabilistic worker departure at natural rate
- **Hiring_Freeze** (1): No new human hiring, natural attrition continues
- **Reduction_In_Force** (2): Active worker removal with acceleration

## Example Configurations

The `examples/` directory contains pre-configured scenarios:

### Small Team Scenarios (10 humans)
- `small_team_natural_attrition.json`: Basic small team with natural attrition
- `small_team_explosive_growth.yaml`: Small team with explosive revenue growth
- `startup_scenario.json`: Early-stage startup with high attrition and fast AI learning

### Medium Team Scenarios (50 humans)
- `medium_team_hiring_freeze.json`: Hiring freeze with accelerated attrition
- `medium_team_fast_learning.yaml`: Fast AI learning with explosive growth

### Large Team Scenarios (200 humans)
- `large_team_reduction_in_force.json`: Large team undergoing workforce reduction
- `large_team_global_distributed.yaml`: Global team with high time zone inefficiency
- `enterprise_conservative.yaml`: Conservative enterprise with slow AI adoption

## Command Line Options

```bash
Usage: ./simulator [options]

Options:
  -config string
        Path to configuration file (JSON or YAML) (default "example_config.json")
  -sensitivity
        Run sensitivity analysis instead of single simulation
  -output string
        Output directory for reports (default ".")
  -help
        Show this help message
```

## Output Formats

### Single Simulation Output

The simulator generates several output files:

1. **Simulation Report** (`simulation_report_YYYYMMDD_HHMMSS.json`):
   - Complete simulation configuration
   - Time-series data of workforce composition
   - Revenue output over time
   - Equilibrium state details
   - Total simulation duration

2. **CSV Export** (`simulation_report_YYYYMMDD_HHMMSS.csv`):
   - Time-series data in spreadsheet format
   - Suitable for visualization tools

### Sensitivity Analysis Output

1. **Sensitivity Report** (`sensitivity_report_YYYYMMDD_HHMMSS.json`):
   - Parameter impact rankings
   - Time to equilibrium for each parameter variation
   - Equilibrium workforce composition variations

2. **Detailed Results** (`sensitivity_detailed_YYYYMMDD_HHMMSS.csv`):
   - Complete results for each parameter variation
   - Suitable for statistical analysis

3. **Parameter Rankings** (`sensitivity_rankings_YYYYMMDD_HHMMSS.csv`):
   - Ranked list of parameters by impact
   - Impact scores for time to equilibrium and workforce composition

## Understanding Results

### Key Metrics

- **Time to Equilibrium**: Number of simulation steps to reach stable workforce composition
- **Equilibrium Composition**: Final ratio of humans to AI agents
- **Total Productivity**: Combined output of all workers
- **Revenue Output**: Monetary value of workforce productivity
- **Orchestration Utilization**: Percentage of human orchestration capacity used

### Interpreting Equilibrium

The simulation reaches equilibrium when:
- Adding more AI agents becomes cost-ineffective
- Workforce composition remains stable over time
- Budget constraints prevent further optimization

## Project Structure

```
workforce-ai-transition-simulator/
├── cmd/
│   └── simulator/          # Main application entry point
├── internal/
│   ├── analytics/          # Analytics engine and reporting
│   ├── controller/         # Simulation controller
│   ├── economic/           # Economic model and budget management
│   ├── events/             # Event processor (attrition, learning, failures)
│   ├── types/              # Core types and configuration
│   └── workforce/          # Workforce manager and worker models
├── examples/               # Example configuration files
├── go.mod                  # Go module definition
└── README.md               # This file
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run property-based tests (requires gopter)
go test -v ./... -tags=property
```

### Building from Source

```bash
# Build for current platform
go build -o simulator ./cmd/simulator

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o simulator.exe ./cmd/simulator

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o simulator ./cmd/simulator

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o simulator ./cmd/simulator
```

## Troubleshooting

### Common Issues

1. **Configuration Validation Errors**
   - Ensure experience distribution percentages sum to 100%
   - Ensure cost category distribution percentages sum to 100%
   - Verify all required fields are present

2. **Memory Issues with Large Simulations**
   - Reduce the number of initial humans for very large scenarios
   - Consider running sensitivity analysis in smaller batches

3. **Long Running Simulations**
   - Some configurations may take a long time to reach equilibrium
   - Monitor progress through generated reports
   - Consider adjusting learning speeds or attrition rates

### Getting Help

- Check the example configurations for reference
- Review the requirements and design documents in `.kiro/specs/`
- Examine the test files for usage examples

## License

This project is licensed under the MIT License - see the LICENSE file for details.
