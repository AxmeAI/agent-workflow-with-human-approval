"""
Agent workflow with human approval — Python example.

AI compliance checker processes a deployment request, then pauses
for human CAB (Change Advisory Board) approval. Human approves via
CLI or email magic link. Agent continues with the decision.

Usage:
    pip install axme
    export AXME_API_KEY="your-key"
    python main.py
"""

import os
from axme import AxmeClient, AxmeClientConfig


def main():
    client = AxmeClient(
        AxmeClientConfig(api_key=os.environ["AXME_API_KEY"])
    )

    # Step 1: Agent performs automated compliance checks
    print("Running automated compliance checks...")
    compliance_result = {
        "pii_scan": "pass",
        "dependency_audit": "pass",
        "security_review": "2 low-severity findings",
        "test_coverage": "94%",
    }
    print(f"Checks complete: {compliance_result}")

    # Step 2: Request human approval from CAB reviewer
    intent_id = client.send_intent(
        {
            "intent_type": "human_approval.v1",
            "to_agent": "agent://myorg/production/cab-reviewer",
            "payload": {
                "action": "deploy",
                "service": "payments-api",
                "environment": "production",
                "risk_level": "high",
                "change_summary": "Upgrade payment processor SDK to v3.2",
                "compliance_result": compliance_result,
                "requestor": "deploy-agent@myorg",
            },
        }
    )
    print(f"\nApproval requested: {intent_id}")
    print("Waiting for CAB reviewer decision...")

    # Step 3: Observe lifecycle events in real time
    for event in client.observe(intent_id):
        status = event.get("status", "")
        print(f"  [{status}] {event.get('event_type', '')}")
        if status in ("COMPLETED", "FAILED", "TIMED_OUT", "CANCELLED"):
            break

    # Step 4: Fetch final decision
    intent = client.get_intent(intent_id)
    decision = intent["intent"]["lifecycle_status"]
    print(f"\nCAB decision: {decision}")

    if decision == "COMPLETED":
        print("Proceeding with deployment...")
    else:
        print("Deployment blocked. Review the decision details.")


if __name__ == "__main__":
    main()
