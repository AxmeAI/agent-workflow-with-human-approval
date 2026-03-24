# Agent Workflow with Human Approval

Your AI agent processes a deployment request. It needs manager approval before proceeding. So you build a webhook endpoint, an email notification service, a polling loop, a timeout handler, and a database to track approval state. That's 200+ lines before your agent does anything useful.

**There is a better way.** Send a human-approval intent, wait for resolution, continue — with built-in reminders, timeouts, and audit trail.

> **Alpha** · Built with [AXME](https://github.com/AxmeAI/axme) (AXP Intent Protocol).
> [cloud.axme.ai](https://cloud.axme.ai) · [hello@axme.ai](mailto:hello@axme.ai)

---

## The Problem

Every team adding human approval to AI agent workflows reinvents the same broken stack:

```
Agent  → needs approval → send email → store state in DB
Agent  → poll DB every 5s... 30 times... still waiting
Human  → clicks link → webhook fires → update DB → notify agent
Agent  → poll DB again → finally sees approval → continue
```

What breaks:
- **HITL is DIY** — 200+ lines of webhook + email + DB + polling for each approval gate
- **No durability** — agent crashes mid-wait, approval state is lost
- **Agents hang** — human doesn't respond for 3 hours, agent blocks forever
- **No timeout semantics** — no escalation, no reminder, no fallback
- **No audit trail** — who approved what, when, with what context?

---

## The Solution: Human Approval Intent

```
Agent  → send_intent("human_approval") → intent_id
Agent  → wait_for(intent_id)

Human  → approves via CLI / email link / form
Agent  ← resumes with approval result
```

One call to request approval, one call to wait. The platform handles delivery, reminders, timeouts, and audit trail.

---

## Quick Start

### Python

```bash
pip install axme
export AXME_API_KEY="your-key"   # Get one: axme login
```

```python
from axme import AxmeClient, AxmeClientConfig
import os

client = AxmeClient(AxmeClientConfig(api_key=os.environ["AXME_API_KEY"]))

# Agent requests human approval for a deployment
intent_id = client.send_intent({
    "intent_type": "human_approval.v1",
    "to_agent": "agent://myorg/production/cab-reviewer",
    "payload": {
        "action": "deploy",
        "service": "payments-api",
        "environment": "production",
        "risk_level": "high",
        "change_summary": "Upgrade payment processor SDK to v3.2",
    },
})

print(f"Approval requested: {intent_id}")

# Wait for human decision — no polling, no webhooks
# Agent suspends durably. Human responds in 5 min or 5 hours.
result = client.wait_for(intent_id)
print(f"Decision: {result['status']}")  # COMPLETED (approved) or REJECTED
```

### TypeScript

```bash
npm install @axme/axme
```

```typescript
import { AxmeClient } from "@axme/axme";

const client = new AxmeClient({ apiKey: process.env.AXME_API_KEY! });

const intentId = await client.sendIntent({
  intentType: "human_approval.v1",
  toAgent: "agent://myorg/production/cab-reviewer",
  payload: {
    action: "deploy",
    service: "payments-api",
    environment: "production",
    riskLevel: "high",
    changeSummary: "Upgrade payment processor SDK to v3.2",
  },
});

console.log(`Approval requested: ${intentId}`);

// Wait for human decision — agent suspends durably
const result = await client.waitFor(intentId);
console.log(`Decision: ${result.status}`);
```

---

## More Languages

Full implementations in all 5 languages:

| Language | Directory | Install |
|----------|-----------|---------|
| [Python](python/) | `python/` | `pip install axme` |
| [TypeScript](typescript/) | `typescript/` | `npm install @axme/axme` |
| [Go](go/) | `go/` | `go get github.com/AxmeAI/axme-sdk-go` |
| [Java](java/) | `java/` | Maven Central: `ai.axme:axme-sdk` |
| [.NET](dotnet/) | `dotnet/` | `dotnet add package Axme.Sdk` |

---

## Before / After

### Before: DIY Approval Gate (200+ lines)

```python
# Webhook endpoint for approval callbacks
@app.post("/webhooks/approval-callback")
async def approval_callback(req):
    data = req.json()
    verify_signature(req.headers["x-signature"], req.body)  # fails silently
    db.update("approvals", data["request_id"], status=data["decision"])
    notify_agent(data["request_id"])  # another webhook...

# Email notification service
def send_approval_email(reviewer_email, request_id, context):
    approval_link = f"https://myapp.com/approve/{request_id}?token={generate_token()}"
    send_email(reviewer_email, "Approval needed", render_template(context, approval_link))
    db.insert("approvals", request_id=request_id, status="pending", created=now())

# Polling loop in the agent (runs every 5 seconds)
async def wait_for_approval(request_id, timeout=3600):
    start = time.time()
    while time.time() - start < timeout:
        row = db.get("approvals", request_id)
        if row["status"] != "pending":
            return row
        await asyncio.sleep(5)  # 720 polls per hour
    raise TimeoutError("Approval timed out")

# Plus: token expiry, reminder emails, escalation logic, audit table,
# DB cleanup cron, orphan request detector, retry on email failure...
```

### After: AXME Human Approval Intent (3 lines)

```python
from axme import AxmeClient, AxmeClientConfig

client = AxmeClient(AxmeClientConfig(api_key=os.environ["AXME_API_KEY"]))

intent_id = client.send_intent({
    "intent_type": "human_approval.v1",
    "to_agent": "agent://myorg/production/cab-reviewer",
    "payload": {"action": "deploy", "service": "payments-api", "environment": "production"},
})

result = client.wait_for(intent_id)
print(result["status"])  # COMPLETED, REJECTED, or TIMED_OUT
```

No webhook endpoint. No email service. No polling loop. No token management. No DB tables. No cleanup cron.

---

## Three Approval Paths

Humans approve through whichever channel fits their workflow:

### CLI

```bash
# Reviewer sees pending approvals
axme tasks list --status pending

# Approve with comment
axme tasks approve intent_abc123 --comment "LGTM, deploy window confirmed"
```

### Email (Magic Link)

```
Reviewer receives email:
  "Approval needed: Deploy payments-api to production"
  [Approve] [Reject] [View Details]

One click. No login required. Cryptographically signed link.
```

### Form (Structured Response)

```
Reviewer opens form link:
  - Decision: [Approve / Reject / Escalate]
  - Conditions: "Only during maintenance window (2-4am UTC)"
  - Risk acknowledgment: [x] I understand this is a high-risk change

Structured response flows back to the agent as intent payload.
```

---

## Use Cases

| Use Case | Agent Does | Human Approves |
|----------|-----------|----------------|
| **Deployment approval** | Runs CI, prepares release | CAB reviewer approves go-live |
| **Budget sign-off** | Calculates cost, generates PO | Finance manager signs off |
| **Content review** | Generates marketing copy | Brand manager reviews tone |
| **Compliance gate** | Scans for PII, flags risks | Compliance officer clears |
| **Data access request** | Validates schema, checks policy | Data owner grants access |
| **Incident escalation** | Detects anomaly, triages | On-call engineer confirms action |

---

## How It Works

```
┌──────────┐  send_intent()     ┌──────────────┐   notify     ┌──────────┐
│   Agent   │ ────────────────► │  AXME Cloud   │ ──────────► │  Human   │
│           │                   │  (platform)   │             │(reviewer)│
│  suspends │                   │               │  approve /  │          │
│  durably  │ ◄── wait_for() ── │  reminders,   │ ◄── reject  │ via CLI, │
│           │   resumes with    │  timeouts,    │             │ email,   │
│ continues │   decision        │  audit trail  │             │ or form  │
└──────────┘                   └──────────────┘             └──────────┘
```

1. Agent sends a **human-approval intent** with context (what, why, risk level)
2. Platform **notifies** the reviewer via configured channel (email, Slack, CLI)
3. Agent **suspends durably** — no resources consumed while waiting
4. Human **approves, rejects, or escalates** through any channel
5. Agent **resumes** with the decision and any structured response
6. Platform records the full **audit trail** (who, when, decision, context)

---

## Works With Any Agent Framework

AXME complements agent frameworks — it handles coordination and durability, not reasoning or planning.

| Framework | How AXME Fits |
|-----------|--------------|
| **LangGraph** | Add approval gates between graph nodes |
| **CrewAI** | Pause crew tasks for human sign-off |
| **AutoGen** | Insert approval checkpoints in multi-agent chat |
| **OpenAI Agents SDK** | Gate tool calls behind human approval |
| **Any Python agent** | `send_intent()` + `wait_for()` in any code |

---

## Run the Full Example

```bash
# Install CLI (one-time)
curl -fsSL https://raw.githubusercontent.com/AxmeAI/axme-cli/main/install.sh | sh
source ~/.zshrc

# Log in
axme login

# Run the built-in example
axme examples run approval/human-gate
```

---

## Related

- [AXME](https://github.com/AxmeAI/axme) — project overview
- [AXP Spec](https://github.com/AxmeAI/axme-spec) — open Intent Protocol specification
- [AXME Examples](https://github.com/AxmeAI/axme-examples) — 20+ runnable examples across 5 languages
- [AXME CLI](https://github.com/AxmeAI/axme-cli) — manage intents, agents, scenarios from the terminal
- [Long-Running API Without Polling](https://github.com/AxmeAI/long-running-api-without-polling) — async operations without webhook glue

---

Built with [AXME](https://github.com/AxmeAI/axme) (AXP Intent Protocol).
