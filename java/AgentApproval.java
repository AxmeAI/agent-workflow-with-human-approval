/*
 * Agent workflow with human approval — Java example.
 *
 * AI compliance checker processes a deployment request, then pauses
 * for human CAB (Change Advisory Board) approval. Human approves via
 * CLI or email magic link. Agent continues with the decision.
 *
 * Usage:
 *   export AXME_API_KEY="your-key"
 *   mvn compile exec:java -Dexec.mainClass="AgentApproval"
 */

import dev.axme.sdk.AxmeClient;
import dev.axme.sdk.AxmeClientConfig;
import dev.axme.sdk.RequestOptions;
import dev.axme.sdk.ObserveOptions;
import java.util.Map;

public class AgentApproval {
    public static void main(String[] args) throws Exception {
        var client = new AxmeClient(
            AxmeClientConfig.forCloud(System.getenv("AXME_API_KEY"))
        );

        // Step 1: Agent performs automated compliance checks
        System.out.println("Running automated compliance checks...");
        var complianceResult = Map.of(
            "pii_scan", "pass",
            "dependency_audit", "pass",
            "security_review", "2 low-severity findings",
            "test_coverage", "94%"
        );
        System.out.println("Checks complete: " + complianceResult);

        // Step 2: Request human approval from CAB reviewer
        String intentId = client.sendIntent(Map.of(
            "intent_type", "human_approval.v1",
            "to_agent", "agent://myorg/production/cab-reviewer",
            "payload", Map.of(
                "action", "deploy",
                "service", "payments-api",
                "environment", "production",
                "risk_level", "high",
                "change_summary", "Upgrade payment processor SDK to v3.2",
                "compliance_result", complianceResult,
                "requestor", "deploy-agent@myorg"
            )
        ), new RequestOptions());
        System.out.println("\nApproval requested: " + intentId);
        System.out.println("Waiting for CAB reviewer decision...");

        // Step 3: Wait for human decision — agent suspends durably
        var result = client.waitFor(intentId, new ObserveOptions());
        System.out.println("\nCAB decision: " + result.get("status"));

        if ("COMPLETED".equals(result.get("status"))) {
            System.out.println("Proceeding with deployment...");
        } else {
            System.out.println("Deployment blocked. Review the decision details.");
        }
    }
}
