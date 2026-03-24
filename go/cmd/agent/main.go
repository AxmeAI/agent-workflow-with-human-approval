// Compliance checker agent — Go example.
//
// Validates deployment compliance. After this agent completes,
// the workflow pauses for human CAB approval.
//
// Usage:
//
//	export AXME_API_KEY="<agent-key>"
//	go run agent.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AxmeAI/axme-sdk-go/axme"
)

const agentAddress = "compliance-checker-demo"

func handleIntent(ctx context.Context, client *axme.Client, intentID string) error {
	intentData, err := client.GetIntent(ctx, intentID, axme.RequestOptions{})
	if err != nil {
		return fmt.Errorf("get intent: %w", err)
	}

	intent, _ := intentData["intent"].(map[string]any)
	if intent == nil {
		intent = intentData
	}
	payload, _ := intent["payload"].(map[string]any)
	if payload == nil {
		payload = map[string]any{}
	}
	if pp, ok := payload["parent_payload"].(map[string]any); ok {
		payload = pp
	}

	service, _ := payload["service"].(string)
	if service == "" {
		service = "unknown"
	}
	version, _ := payload["version"].(string)
	if version == "" {
		version = "unknown"
	}
	env, _ := payload["environment"].(string)
	if env == "" {
		env = "unknown"
	}
	risk, _ := payload["risk_level"].(string)
	if risk == "" {
		risk = "unknown"
	}

	fmt.Printf("  Checking compliance for %s v%s -> %s...\n", service, version, env)
	time.Sleep(1 * time.Second)

	fmt.Printf("  Risk level: %s\n", risk)
	fmt.Println("  Validating security policies...")
	time.Sleep(1 * time.Second)

	fmt.Println("  Checking rollback plan...")
	time.Sleep(1 * time.Second)

	passed := risk != "critical"
	riskAssessment := "passed"
	if !passed {
		riskAssessment = "failed"
	}
	recommendation := "approve"
	if !passed {
		recommendation = "reject"
	}

	result := map[string]any{
		"action": "complete",
		"passed": passed,
		"checks": map[string]any{
			"security_policy": "passed",
			"rollback_plan":   "passed",
			"dependency_scan": "passed",
			"risk_assessment": riskAssessment,
		},
		"recommendation": recommendation,
	}

	_, err = client.ResumeIntent(ctx, intentID, result, axme.RequestOptions{})
	if err != nil {
		return fmt.Errorf("resume intent: %w", err)
	}

	statusLabel := "PASSED"
	if !passed {
		statusLabel = "FAILED"
	}
	fmt.Printf("  Compliance check %s. Recommendation: %s\n", statusLabel, recommendation)
	fmt.Println("  Workflow now waits for human CAB approval.")
	fmt.Println("  To approve: axme tasks approve <intent_id>")
	return nil
}

func main() {
	apiKey := os.Getenv("AXME_API_KEY")
	if apiKey == "" {
		log.Fatal("Error: AXME_API_KEY not set.")
	}

	client, err := axme.NewClient(axme.ClientConfig{APIKey: apiKey})
	if err != nil {
		log.Fatalf("create client: %v", err)
	}

	ctx := context.Background()

	fmt.Printf("Agent listening on %s...\n", agentAddress)
	fmt.Println("Waiting for intents (Ctrl+C to stop)")

	intents, errCh := client.Listen(ctx, agentAddress, axme.ListenOptions{})

	go func() {
		for err := range errCh {
			log.Printf("Listen error: %v", err)
		}
	}()

	for delivery := range intents {
		intentID, _ := delivery["intent_id"].(string)
		status, _ := delivery["status"].(string)
		if intentID == "" {
			continue
		}
		if status == "DELIVERED" || status == "CREATED" || status == "IN_PROGRESS" {
			fmt.Printf("[%s] Intent received: %s\n", status, intentID)
			if err := handleIntent(ctx, client, intentID); err != nil {
				fmt.Printf("  Error: %v\n", err)
			}
		}
	}
}
