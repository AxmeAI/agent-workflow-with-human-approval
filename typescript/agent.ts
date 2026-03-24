/**
 * Compliance checker agent — TypeScript example.
 *
 * Validates deployment compliance. After this agent completes,
 * the workflow pauses for human CAB approval.
 *
 * Usage:
 *   export AXME_API_KEY="<agent-key>"
 *   npx tsx agent.ts
 */

import { AxmeClient } from "@axme/axme";

const AGENT_ADDRESS = "compliance-checker-demo";

async function handleIntent(client: AxmeClient, intentId: string) {
  const intentData = await client.getIntent(intentId);
  const intent = intentData.intent ?? intentData;
  let payload = intent.payload ?? {};
  if (payload.parent_payload) {
    payload = payload.parent_payload;
  }

  const service = payload.service ?? "unknown";
  const version = payload.version ?? "unknown";
  const env = payload.environment ?? "unknown";
  const risk = payload.risk_level ?? "unknown";

  console.log(`  Checking compliance for ${service} v${version} -> ${env}...`);
  await new Promise((r) => setTimeout(r, 1000));

  console.log(`  Risk level: ${risk}`);
  console.log(`  Validating security policies...`);
  await new Promise((r) => setTimeout(r, 1000));

  console.log(`  Checking rollback plan...`);
  await new Promise((r) => setTimeout(r, 1000));

  const passed = risk !== "critical";
  const result = {
    action: "complete",
    passed,
    checks: {
      security_policy: "passed",
      rollback_plan: "passed",
      dependency_scan: "passed",
      risk_assessment: passed ? "passed" : "failed",
    },
    recommendation: passed ? "approve" : "reject",
  };

  await client.resumeIntent(intentId, result, { ownerAgent: "compliance-checker-demo" });
  const checkStatus = passed ? "PASSED" : "FAILED";
  console.log(`  Compliance check ${checkStatus}. Recommendation: ${result.recommendation}`);
  console.log(`  Workflow now waits for human CAB approval.`);
  console.log(`  To approve: axme tasks approve <intent_id>`);
}

async function main() {
  const apiKey = process.env.AXME_API_KEY;
  if (!apiKey) {
    console.error("Error: AXME_API_KEY not set.");
    process.exit(1);
  }

  const client = new AxmeClient({ apiKey });

  console.log(`Agent listening on ${AGENT_ADDRESS}...`);
  console.log("Waiting for intents (Ctrl+C to stop)\n");

  for await (const delivery of client.listen(AGENT_ADDRESS)) {
    const intentId = delivery.intent_id;
    const status = delivery.status;
    if (!intentId) continue;
    if (["DELIVERED", "CREATED", "IN_PROGRESS"].includes(status)) {
      console.log(`[${status}] Intent received: ${intentId}`);
      try {
        await handleIntent(client, intentId);
      } catch (e) {
        console.error(`  Error: ${e}`);
      }
    }
  }
}

main().catch(console.error);
