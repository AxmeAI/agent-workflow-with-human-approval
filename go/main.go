// Agent workflow with human approval — Go example.
//
// AI compliance checker processes a deployment request, then pauses
// for human CAB (Change Advisory Board) approval. Human approves via
// CLI or email magic link. Agent continues with the decision.
//
// Usage:
//
//	export AXME_API_KEY="your-key"
//	go run main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/AxmeAI/axme-sdk-go/axme"
)

func main() {
	client := axme.NewClient(axme.Config{
		APIKey: os.Getenv("AXME_API_KEY"),
	})

	ctx := context.Background()

	// Step 1: Agent performs automated compliance checks
	fmt.Println("Running automated compliance checks...")
	complianceResult := map[string]interface{}{
		"pii_scan":         "pass",
		"dependency_audit": "pass",
		"security_review":  "2 low-severity findings",
		"test_coverage":    "94%",
	}
	fmt.Printf("Checks complete: %v\n", complianceResult)

	// Step 2: Request human approval from CAB reviewer
	intentID, err := client.SendIntent(ctx, axme.SendIntentRequest{
		IntentType: "human_approval.v1",
		ToAgent:    "agent://myorg/production/cab-reviewer",
		Payload: map[string]interface{}{
			"action":            "deploy",
			"service":           "payments-api",
			"environment":       "production",
			"risk_level":        "high",
			"change_summary":    "Upgrade payment processor SDK to v3.2",
			"compliance_result": complianceResult,
			"requestor":         "deploy-agent@myorg",
		},
	})
	if err != nil {
		log.Fatalf("send intent: %v", err)
	}
	fmt.Printf("\nApproval requested: %s\n", intentID)
	fmt.Println("Waiting for CAB reviewer decision...")

	// Step 3: Wait for human decision — agent suspends durably
	result, err := client.WaitFor(ctx, intentID)
	if err != nil {
		log.Fatalf("wait: %v", err)
	}
	fmt.Printf("\nCAB decision: %s\n", result.Status)

	if result.Status == "COMPLETED" {
		fmt.Println("Proceeding with deployment...")
	} else {
		fmt.Println("Deployment blocked. Review the decision details.")
	}
}
