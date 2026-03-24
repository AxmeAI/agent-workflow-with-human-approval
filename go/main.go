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
	client, err := axme.NewClient(axme.ClientConfig{
		APIKey: os.Getenv("AXME_API_KEY"),
	})
	if err != nil {
		log.Fatalf("create client: %v", err)
	}

	ctx := context.Background()

	// Step 1: Agent performs automated compliance checks
	fmt.Println("Running automated compliance checks...")
	complianceResult := map[string]any{
		"pii_scan":         "pass",
		"dependency_audit": "pass",
		"security_review":  "2 low-severity findings",
		"test_coverage":    "94%",
	}
	fmt.Printf("Checks complete: %v\n", complianceResult)

	// Step 2: Request human approval from CAB reviewer
	intentID, err := client.SendIntent(ctx, map[string]any{
		"intent_type":       "human_approval.v1",
		"to_agent":          "agent://myorg/production/cab-reviewer",
		"action":            "deploy",
		"service":           "payments-api",
		"environment":       "production",
		"risk_level":        "high",
		"change_summary":    "Upgrade payment processor SDK to v3.2",
		"compliance_result": complianceResult,
		"requestor":         "deploy-agent@myorg",
	}, axme.RequestOptions{})
	if err != nil {
		log.Fatalf("send intent: %v", err)
	}
	fmt.Printf("\nApproval requested: %s\n", intentID)
	fmt.Println("Waiting for CAB reviewer decision...")

	// Step 3: Wait for human decision — agent suspends durably
	result, err := client.WaitFor(ctx, intentID, axme.ObserveOptions{})
	if err != nil {
		log.Fatalf("wait: %v", err)
	}
	fmt.Printf("\nCAB decision: %v\n", result["status"])

	if result["status"] == "COMPLETED" {
		fmt.Println("Proceeding with deployment...")
	} else {
		fmt.Println("Deployment blocked. Review the decision details.")
	}
}
