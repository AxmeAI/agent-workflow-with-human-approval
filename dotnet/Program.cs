// Agent workflow with human approval — .NET example.
//
// AI compliance checker processes a deployment request, then pauses
// for human CAB (Change Advisory Board) approval. Human approves via
// CLI or email magic link. Agent continues with the decision.
//
// Usage:
//   export AXME_API_KEY="your-key"
//   dotnet run

using Axme.Sdk;

var client = new AxmeClient(new AxmeClientConfig
{
    ApiKey = Environment.GetEnvironmentVariable("AXME_API_KEY")!
});

// Step 1: Agent performs automated compliance checks
Console.WriteLine("Running automated compliance checks...");
var complianceResult = new
{
    pii_scan = "pass",
    dependency_audit = "pass",
    security_review = "2 low-severity findings",
    test_coverage = "94%"
};
Console.WriteLine($"Checks complete: {complianceResult}");

// Step 2: Request human approval from CAB reviewer
var intentId = await client.SendIntentAsync(new
{
    intent_type = "human_approval.v1",
    to_agent = "agent://myorg/production/cab-reviewer",
    payload = new
    {
        action = "deploy",
        service = "payments-api",
        environment = "production",
        risk_level = "high",
        change_summary = "Upgrade payment processor SDK to v3.2",
        compliance_result = complianceResult,
        requestor = "deploy-agent@myorg"
    }
});
Console.WriteLine($"\nApproval requested: {intentId}");
Console.WriteLine("Waiting for CAB reviewer decision...");

// Step 3: Wait for human decision — agent suspends durably
var result = await client.WaitForAsync(intentId);
Console.WriteLine($"\nCAB decision: {result.Status}");

if (result.Status == "COMPLETED")
{
    Console.WriteLine("Proceeding with deployment...");
}
else
{
    Console.WriteLine("Deployment blocked. Review the decision details.");
}
