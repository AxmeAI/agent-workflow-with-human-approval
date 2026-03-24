/*
 * Compliance checker agent — Java example.
 *
 * Fetches an intent by ID, runs compliance checks, and resumes with result.
 * After this agent completes, the workflow pauses for human CAB approval.
 *
 * Usage:
 *   export AXME_API_KEY="<agent-key>"
 *   javac -cp axme-sdk.jar Agent.java
 *   java -cp .:axme-sdk.jar Agent <intent_id>
 */

import dev.axme.sdk.AxmeClient;
import dev.axme.sdk.AxmeClientConfig;
import dev.axme.sdk.RequestOptions;
import java.util.Map;

public class Agent {
    public static void main(String[] args) throws Exception {
        if (args.length < 1) {
            System.err.println("Usage: java Agent <intent_id>");
            System.exit(1);
        }

        String apiKey = System.getenv("AXME_API_KEY");
        if (apiKey == null || apiKey.isEmpty()) {
            System.err.println("Error: AXME_API_KEY not set.");
            System.exit(1);
        }

        String intentId = args[0];
        var client = new AxmeClient(AxmeClientConfig.forCloud(apiKey));

        System.out.println("Processing intent: " + intentId);

        var intentData = client.getIntent(intentId, new RequestOptions());
        @SuppressWarnings("unchecked")
        var intent = (Map<String, Object>) intentData.getOrDefault("intent", intentData);
        @SuppressWarnings("unchecked")
        var payload = (Map<String, Object>) intent.getOrDefault("payload", Map.of());
        if (payload.containsKey("parent_payload")) {
            @SuppressWarnings("unchecked")
            var pp = (Map<String, Object>) payload.get("parent_payload");
            payload = pp;
        }

        String service = (String) payload.getOrDefault("service", "unknown");
        String version = (String) payload.getOrDefault("version", "unknown");
        String env = (String) payload.getOrDefault("environment", "unknown");
        String risk = (String) payload.getOrDefault("risk_level", "unknown");

        System.out.println("  Checking compliance for " + service + " v" + version + " -> " + env + "...");
        Thread.sleep(1000);

        System.out.println("  Risk level: " + risk);
        System.out.println("  Validating security policies...");
        Thread.sleep(1000);

        System.out.println("  Checking rollback plan...");
        Thread.sleep(1000);

        boolean passed = !risk.equals("critical");
        String riskAssessment = passed ? "passed" : "failed";
        String recommendation = passed ? "approve" : "reject";

        var result = Map.<String, Object>of(
            "action", "complete",
            "passed", passed,
            "checks", Map.of(
                "security_policy", "passed",
                "rollback_plan", "passed",
                "dependency_scan", "passed",
                "risk_assessment", riskAssessment
            ),
            "recommendation", recommendation
        );

        client.resumeIntent(intentId, result, new RequestOptions());
        String status = passed ? "PASSED" : "FAILED";
        System.out.println("  Compliance check " + status + ". Recommendation: " + recommendation);
        System.out.println("  Workflow now waits for human CAB approval.");
        System.out.println("  To approve: axme tasks approve <intent_id>");
    }
}
