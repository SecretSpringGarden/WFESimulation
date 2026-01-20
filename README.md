# Workforce AI Transition Simulator

A discrete-event simulation system that models the evolution of an engineering workforce from purely human to AI-augmented composition under budget constraints.

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

## Project Structure

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
├── pkg/                    # Public packages (to be added)
├── go.mod                  # Go module definition
└── README.md               # This file
```

## Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...
```

## Building

```bash
# Build the simulator
go build -o simulator ./cmd/simulator

# Run the simulator
./simulator
```

## Configuration

Configuration parameters are defined in `internal/types/config.go`. The simulator accepts:

- Initial workforce size and distribution
- Fixed budget constraints
- Revenue growth scenarios
- AI learning speed parameters
- Attrition configuration
- Catastrophic failure rates
- Time zone inefficiency penalties

## Development Status

This project is currently under development. Core types and project structure have been established.

## License

TBD
