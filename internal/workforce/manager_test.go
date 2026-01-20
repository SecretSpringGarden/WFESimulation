package workforce

import (
	"testing"
	"workforce-ai-transition-simulator/internal/types"
)

func TestNewWorkforceManager(t *testing.T) {
	wm := NewWorkforceManager()
	
	if wm == nil {
		t.Fatal("NewWorkforceManager() returned nil")
	}
	
	if len(wm.humans) != 0 {
		t.Errorf("Expected 0 humans, got %d", len(wm.humans))
	}
	
	if len(wm.aiAgents) != 0 {
		t.Errorf("Expected 0 AI agents, got %d", len(wm.aiAgents))
	}
}

func TestAddHuman(t *testing.T) {
	wm := NewWorkforceManager()
	
	// Add first human (should become business owner)
	human1, err := wm.AddHuman(types.MidLevel, types.HighCostUS, false)
	if err != nil {
		t.Fatalf("AddHuman() error = %v", err)
	}
	
	if !human1.IsBusinessOwner {
		t.Error("First human should be business owner")
	}
	
	if wm.businessOwnerID != human1.ID {
		t.Error("Business owner ID not set correctly")
	}
	
	// Add second human (should not be business owner)
	human2, err := wm.AddHuman(types.Senior, types.LowCostNonUS, false)
	if err != nil {
		t.Fatalf("AddHuman() error = %v", err)
	}
	
	if human2.IsBusinessOwner {
		t.Error("Second human should not be business owner")
	}
	
	// Try to add another business owner (should fail)
	_, err = wm.AddHuman(types.Executive, types.HighCostUS, true)
	if err == nil {
		t.Error("Expected error when adding second business owner")
	}
}

func TestRemoveHuman(t *testing.T) {
	wm := NewWorkforceManager()
	
	// Add humans
	human1, _ := wm.AddHuman(types.MidLevel, types.HighCostUS, false)
	human2, _ := wm.AddHuman(types.Senior, types.LowCostNonUS, false)
	
	// Add AI agents to human2
	agent1, _ := wm.AddAIAgent(human2.ID, 0)
	agent2, _ := wm.AddAIAgent(human2.ID, 0)
	
	// Try to remove business owner (should fail)
	err := wm.RemoveHuman(human1.ID)
	if err == nil {
		t.Error("Expected error when removing business owner")
	}
	
	// Remove human2 (should also remove their agents)
	err = wm.RemoveHuman(human2.ID)
	if err != nil {
		t.Fatalf("RemoveHuman() error = %v", err)
	}
	
	// Verify human2 is removed
	_, exists := wm.humans[human2.ID]
	if exists {
		t.Error("Human should be removed")
	}
	
	// Verify agents are removed
	_, exists = wm.aiAgents[agent1.ID]
	if exists {
		t.Error("Agent1 should be removed")
	}
	
	_, exists = wm.aiAgents[agent2.ID]
	if exists {
		t.Error("Agent2 should be removed")
	}
}

func TestAddAIAgent(t *testing.T) {
	wm := NewWorkforceManager()
	
	// Add human
	human, _ := wm.AddHuman(types.MidLevel, types.HighCostUS, false)
	
	// Add AI agent
	agent, err := wm.AddAIAgent(human.ID, 10)
	if err != nil {
		t.Fatalf("AddAIAgent() error = %v", err)
	}
	
	if agent.OrchestratorID != human.ID {
		t.Error("Agent orchestrator ID not set correctly")
	}
	
	if agent.CreationTime != 10 {
		t.Error("Agent creation time not set correctly")
	}
	
	// Verify agent is in collection
	_, exists := wm.aiAgents[agent.ID]
	if !exists {
		t.Error("Agent should be in collection")
	}
	
	// Verify agent is assigned to human
	if len(human.AssignedAgents) != 1 {
		t.Errorf("Expected 1 assigned agent, got %d", len(human.AssignedAgents))
	}
	
	// Try to add agent to non-existent orchestrator
	_, err = wm.AddAIAgent("non-existent", 0)
	if err == nil {
		t.Error("Expected error when adding agent to non-existent orchestrator")
	}
}

func TestAddAIAgentCapacityLimit(t *testing.T) {
	wm := NewWorkforceManager()
	
	// Add human
	human, _ := wm.AddHuman(types.MidLevel, types.HighCostUS, false)
	
	// Add 6 agents (max capacity)
	for i := 0; i < 6; i++ {
		_, err := wm.AddAIAgent(human.ID, i)
		if err != nil {
			t.Fatalf("AddAIAgent() error = %v at iteration %d", err, i)
		}
	}
	
	// Try to add 7th agent (should fail)
	_, err := wm.AddAIAgent(human.ID, 6)
	if err == nil {
		t.Error("Expected error when exceeding orchestration limit")
	}
}

func TestReleaseAIAgent(t *testing.T) {
	wm := NewWorkforceManager()
	
	// Add human and agents
	human, _ := wm.AddHuman(types.MidLevel, types.HighCostUS, false)
	agent1, _ := wm.AddAIAgent(human.ID, 0)
	agent2, _ := wm.AddAIAgent(human.ID, 0)
	
	// Release agent1
	err := wm.ReleaseAIAgent(agent1.ID)
	if err != nil {
		t.Fatalf("ReleaseAIAgent() error = %v", err)
	}
	
	// Verify agent1 is removed
	_, exists := wm.aiAgents[agent1.ID]
	if exists {
		t.Error("Agent should be removed from collection")
	}
	
	// Verify agent1 is removed from human's assigned list
	if len(human.AssignedAgents) != 1 {
		t.Errorf("Expected 1 assigned agent, got %d", len(human.AssignedAgents))
	}
	
	// Verify agent2 is still there
	if human.AssignedAgents[0] != agent2.ID {
		t.Error("Wrong agent in assigned list")
	}
	
	// Try to release non-existent agent
	err = wm.ReleaseAIAgent("non-existent")
	if err == nil {
		t.Error("Expected error when releasing non-existent agent")
	}
}

func TestGetAvailableOrchestrationCapacity(t *testing.T) {
	wm := NewWorkforceManager()
	
	// Initially should be 0
	capacity := wm.GetAvailableOrchestrationCapacity()
	if capacity != 0 {
		t.Errorf("Expected 0 capacity, got %d", capacity)
	}
	
	// Add 2 humans (12 total capacity)
	human1, _ := wm.AddHuman(types.MidLevel, types.HighCostUS, false)
	_, _ = wm.AddHuman(types.Senior, types.LowCostNonUS, false)
	
	capacity = wm.GetAvailableOrchestrationCapacity()
	if capacity != 12 {
		t.Errorf("Expected 12 capacity, got %d", capacity)
	}
	
	// Add 3 agents to human1
	wm.AddAIAgent(human1.ID, 0)
	wm.AddAIAgent(human1.ID, 0)
	wm.AddAIAgent(human1.ID, 0)
	
	capacity = wm.GetAvailableOrchestrationCapacity()
	if capacity != 9 {
		t.Errorf("Expected 9 capacity, got %d", capacity)
	}
}

func TestCalculateTotalProductivity(t *testing.T) {
	wm := NewWorkforceManager()
	
	// Add humans
	human1, _ := wm.AddHuman(types.MidLevel, types.HighCostUS, false)      // productivity: 2.0
	_, _ = wm.AddHuman(types.Senior, types.LowCostNonUS, false)            // productivity: 3.5 * 0.8 = 2.8 (with 20% penalty)
	
	// Add AI agents
	wm.AddAIAgent(human1.ID, 0) // University hire: 0.8
	wm.AddAIAgent(human1.ID, 0) // University hire: 0.8
	
	// Calculate with 20% time zone inefficiency
	productivity := wm.CalculateTotalProductivity(0.2)
	expected := 2.0 + 2.8 + 0.8 + 0.8 // 6.4
	
	const tolerance = 1e-9
	if diff := productivity - expected; diff < -tolerance || diff > tolerance {
		t.Errorf("Expected productivity %v, got %v", expected, productivity)
	}
}

func TestGetWorkforceComposition(t *testing.T) {
	wm := NewWorkforceManager()
	
	// Add humans
	wm.AddHuman(types.MidLevel, types.HighCostUS, false)
	wm.AddHuman(types.MidLevel, types.HighCostUS, false)
	human3, _ := wm.AddHuman(types.Senior, types.LowCostNonUS, false)
	
	// Add AI agents
	agent1, _ := wm.AddAIAgent(human3.ID, 0)
	wm.AddAIAgent(human3.ID, 0)
	
	// Level up agent1 to MidLevel
	agent1.ExperienceLevel = types.MidLevel
	
	composition := wm.GetWorkforceComposition()
	
	// Check human counts
	if composition.Humans.Total != 3 {
		t.Errorf("Expected 3 humans, got %d", composition.Humans.Total)
	}
	
	if composition.Humans.ByExperience[types.MidLevel] != 2 {
		t.Errorf("Expected 2 MidLevel humans, got %d", composition.Humans.ByExperience[types.MidLevel])
	}
	
	if composition.Humans.ByExperience[types.Senior] != 1 {
		t.Errorf("Expected 1 Senior human, got %d", composition.Humans.ByExperience[types.Senior])
	}
	
	if composition.Humans.ByCostCategory[types.HighCostUS] != 2 {
		t.Errorf("Expected 2 HighCostUS humans, got %d", composition.Humans.ByCostCategory[types.HighCostUS])
	}
	
	if composition.Humans.ByCostCategory[types.LowCostNonUS] != 1 {
		t.Errorf("Expected 1 LowCostNonUS human, got %d", composition.Humans.ByCostCategory[types.LowCostNonUS])
	}
	
	// Check AI agent counts
	if composition.AIAgents.Total != 2 {
		t.Errorf("Expected 2 AI agents, got %d", composition.AIAgents.Total)
	}
	
	if composition.AIAgents.ByExperience[types.UniversityHire] != 1 {
		t.Errorf("Expected 1 UniversityHire agent, got %d", composition.AIAgents.ByExperience[types.UniversityHire])
	}
	
	if composition.AIAgents.ByExperience[types.MidLevel] != 1 {
		t.Errorf("Expected 1 MidLevel agent, got %d", composition.AIAgents.ByExperience[types.MidLevel])
	}
	
	// Check orchestration utilization (2 agents / 18 total capacity = 11.11%)
	expectedUtilization := (2.0 / 18.0) * 100.0
	const tolerance = 0.01
	if diff := composition.OrchestrationUtilization - expectedUtilization; diff < -tolerance || diff > tolerance {
		t.Errorf("Expected utilization %v%%, got %v%%", expectedUtilization, composition.OrchestrationUtilization)
	}
}
