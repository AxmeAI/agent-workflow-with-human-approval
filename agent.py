"""
Compliance checker agent - validates deployment compliance.

After this agent completes, the workflow pauses for human CAB approval.
Human approves via: axme tasks approve <intent_id>

Usage:
    export AXME_API_KEY="<agent-key>"
    python agent.py
"""

import os
import sys
import time

sys.stdout.reconfigure(line_buffering=True)

from axme import AxmeClient, AxmeClientConfig


AGENT_ADDRESS = "compliance-checker-demo"


def handle_intent(client, intent_id):
    """Run compliance checks and resume with results."""
    intent_data = client.get_intent(intent_id)
    intent = intent_data.get("intent", intent_data)
    payload = intent.get("payload", {})
    if "parent_payload" in payload:
        payload = payload["parent_payload"]

    service = payload.get("service", "unknown")
    version = payload.get("version", "unknown")
    env = payload.get("environment", "unknown")
    risk = payload.get("risk_level", "unknown")

    print(f"  Checking compliance for {service} v{version} -> {env}...")
    time.sleep(1)

    print(f"  Risk level: {risk}")
    print(f"  Validating security policies...")
    time.sleep(1)

    print(f"  Checking rollback plan...")
    time.sleep(1)

    passed = risk != "critical"
    result = {
        "action": "complete",
        "passed": passed,
        "checks": {
            "security_policy": "passed",
            "rollback_plan": "passed",
            "dependency_scan": "passed",
            "risk_assessment": "passed" if passed else "failed",
        },
        "recommendation": "approve" if passed else "reject",
    }

    client.resume_intent(intent_id, result)
    status = "PASSED" if passed else "FAILED"
    print(f"  Compliance check {status}. Recommendation: {result['recommendation']}")
    print(f"  Workflow now waits for human CAB approval.")
    print(f"  To approve: axme tasks approve <intent_id>")


def main():
    api_key = os.environ.get("AXME_API_KEY", "")
    if not api_key:
        print("Error: AXME_API_KEY not set.")
        sys.exit(1)

    client = AxmeClient(AxmeClientConfig(api_key=api_key))
    print(f"Agent listening on {AGENT_ADDRESS}...")
    print("Waiting for intents (Ctrl+C to stop)\n")

    for delivery in client.listen(AGENT_ADDRESS):
        intent_id = delivery.get("intent_id", "")
        status = delivery.get("status", "")
        if not intent_id:
            continue
        if status in ("DELIVERED", "CREATED", "IN_PROGRESS"):
            print(f"[{status}] Intent received: {intent_id}")
            try:
                handle_intent(client, intent_id)
            except Exception as e:
                print(f"  Error: {e}")


if __name__ == "__main__":
    main()
