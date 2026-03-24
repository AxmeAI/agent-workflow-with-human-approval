/**
 * Agent workflow with human approval — TypeScript example.
 *
 * AI compliance checker processes a deployment request, then pauses
 * for human CAB (Change Advisory Board) approval. Human approves via
 * CLI or email magic link. Agent continues with the decision.
 *
 * Usage:
 *   npm install @axme/axme
 *   export AXME_API_KEY="your-key"
 *   npx tsx main.ts
 */

import { AxmeClient } from "@axme/axme";

async function main() {
  const client = new AxmeClient({ apiKey: process.env.AXME_API_KEY! });

  // Step 1: Agent performs automated compliance checks
  console.log("Running automated compliance checks...");
  const complianceResult = {
    piiScan: "pass",
    dependencyAudit: "pass",
    securityReview: "2 low-severity findings",
    testCoverage: "94%",
  };
  console.log("Checks complete:", complianceResult);

  // Step 2: Request human approval from CAB reviewer
  const intentId = await client.sendIntent({
    intentType: "human_approval.v1",
    toAgent: "agent://myorg/production/cab-reviewer",
    payload: {
      action: "deploy",
      service: "payments-api",
      environment: "production",
      riskLevel: "high",
      changeSummary: "Upgrade payment processor SDK to v3.2",
      complianceResult,
      requestor: "deploy-agent@myorg",
    },
  });
  console.log(`\nApproval requested: ${intentId}`);
  console.log("Waiting for CAB reviewer decision...");

  // Step 3: Wait for human decision — agent suspends durably
  const result = await client.waitFor(intentId);
  console.log(`\nCAB decision: ${result.status}`);

  if (result.status === "COMPLETED") {
    console.log("Proceeding with deployment...");
  } else {
    console.log("Deployment blocked. Review the decision details.");
  }
}

main().catch(console.error);
