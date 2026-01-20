# Setup Instructions

## Current Status

The project structure and core types have been created. However, **Go is not currently installed** on this system, so we cannot run tests or build the project yet.

## Next Steps

### 1. Install Go

Download and install Go 1.21 or higher from: https://golang.org/dl/

For Windows:
- Download the Windows installer (.msi file)
- Run the installer and follow the prompts
- Restart your terminal/command prompt after installation

### 2. Verify Go Installation

Open a new terminal/command prompt and run:
```bash
go version
```

You should see output like: `go version go1.21.x windows/amd64`

### 3. Initialize the Project

Once Go is installed, run the setup script:

**Windows:**
```bash
setup.bat
```

**Linux/Mac:**
```bash
chmod +x setup.sh
./setup.sh
```

Or manually:
```bash
# Download dependencies
go mod download

# Run tests
go test ./...

# Build the simulator
go build -o simulator ./cmd/simulator
```

## What Has Been Created

### Directory Structure
```
workforce-ai-transition-simulator/
├── cmd/
│   └── simulator/          # Main application entry point
│       └── main.go
├── internal/
│   └── types/              # Core types and configuration
│       ├── core.go         # Enums and constants
│       ├── config.go       # Configuration structs
│       └── core_test.go    # Unit tests
├── pkg/                    # Public packages (for future use)
├── go.mod                  # Go module definition with gopter dependency
├── .gitignore              # Git ignore rules
├── README.md               # Project documentation
├── setup.sh                # Linux/Mac setup script
└── setup.bat               # Windows setup script
```

### Core Types Defined

1. **ExperienceLevel** enum: UniversityHire, MidLevel, Senior, Executive
2. **CostCategory** enum: HighCostUS, LowCostNonUS
3. **RevenueScenario** enum: FlatRevenue, ExplosiveGrowth
4. **AttritionType** enum: NaturalAttrition, HiringFreeze, ReductionInForce
5. **OrchestrationLimit** constant: 6 (max AI agents per human)

### Configuration Structures

1. **ExperienceDistribution**: Percentage distribution across experience levels
2. **CostCategoryDistribution**: Percentage distribution across cost categories
3. **AILearningSpeed**: Time steps required for AI progression
4. **AttritionConfig**: Attrition behavior configuration
5. **SimulationConfig**: Complete simulation configuration with all parameters

### Testing Framework

- Go's built-in testing package is configured
- gopter (property-based testing library) is included in go.mod
- Basic unit tests created for core type string representations

## Requirements Satisfied

This implementation satisfies the following requirements from task 1:
- ✅ Create Go module with appropriate directory structure (cmd/, internal/, pkg/)
- ✅ Define core enums and constants (ExperienceLevel, CostCategory, RevenueScenario, AttritionType)
- ✅ Define SimulationConfig struct with all configuration parameters
- ✅ Set up testing framework (Go's built-in testing + gopter for property-based testing)

**Requirements covered:** 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7, 1.8
