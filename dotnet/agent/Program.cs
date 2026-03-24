// Compliance checker agent — .NET example.
//
// Fetches an intent by ID, runs compliance checks, and resumes with result.
// After this agent completes, the workflow pauses for human CAB approval.
//
// Usage:
//   export AXME_API_KEY="<agent-key>"
//   dotnet run -- <intent_id>

using Axme.Sdk;
using System.Text.Json.Nodes;

if (args.Length < 1)
{
    Console.Error.WriteLine("Usage: dotnet run -- <intent_id>");
    return 1;
}

var apiKey = Environment.GetEnvironmentVariable("AXME_API_KEY");
if (string.IsNullOrEmpty(apiKey))
{
    Console.Error.WriteLine("Error: AXME_API_KEY not set.");
    return 1;
}

var intentId = args[0];
var client = new AxmeClient(new AxmeClientConfig { ApiKey = apiKey });

Console.WriteLine($"Processing intent: {intentId}");

var intentData = await client.GetIntentAsync(intentId);
var intent = intentData["intent"]?.AsObject() ?? intentData;
var payload = intent["payload"]?.AsObject() ?? new JsonObject();
if (payload["parent_payload"] is JsonObject parentPayload)
{
    payload = parentPayload;
}

var service = payload["service"]?.ToString() ?? "unknown";
var version = payload["version"]?.ToString() ?? "unknown";
var env = payload["environment"]?.ToString() ?? "unknown";
var risk = payload["risk_level"]?.ToString() ?? "unknown";

Console.WriteLine($"  Checking compliance for {service} v{version} -> {env}...");
await Task.Delay(1000);

Console.WriteLine($"  Risk level: {risk}");
Console.WriteLine("  Validating security policies...");
await Task.Delay(1000);

Console.WriteLine("  Checking rollback plan...");
await Task.Delay(1000);

var passed = risk != "critical";
var riskAssessment = passed ? "passed" : "failed";
var recommendation = passed ? "approve" : "reject";

var result = new JsonObject
{
    ["action"] = "complete",
    ["passed"] = passed,
    ["checks"] = new JsonObject
    {
        ["security_policy"] = "passed",
        ["rollback_plan"] = "passed",
        ["dependency_scan"] = "passed",
        ["risk_assessment"] = riskAssessment
    },
    ["recommendation"] = recommendation
};

await client.ResumeIntentAsync(intentId, result);
var status = passed ? "PASSED" : "FAILED";
Console.WriteLine($"  Compliance check {status}. Recommendation: {recommendation}");
Console.WriteLine("  Workflow now waits for human CAB approval.");
Console.WriteLine("  To approve: axme tasks approve <intent_id>");
return 0;
