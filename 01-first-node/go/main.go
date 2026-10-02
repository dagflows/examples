package main

import df "github.com/dagflows/sdk-go"

// Node handlers follow the signature func(ctx *df.Ctx, in In) (Out, error).
// Root nodes take df.None since they have no parent inputs.
func ingestOrders(ctx *df.Ctx, _ df.None) (map[string]any, error) {
	orders := []map[string]any{
		{"id": 1001, "customer_id": "cus_alpha", "amount_cents": 1250, "currency": "usd"},
		{"id": 1002, "customer_id": "cus_beta", "amount_cents": 34000, "currency": "USD"},
		{"id": 1003, "customer_id": "cus_gamma", "amount_cents": 9999, "currency": "usd"},
		{"id": 1004, "customer_id": "cus_delta", "amount_cents": 500, "currency": "eur"},
		{"id": 1005, "customer_id": "cus_epsilon", "amount_cents": 0, "currency": "usd"},
	}

	// ctx.Log() is a *slog.Logger attached to the workflow run's log stream.
	ctx.Log().Info("ingested orders", "run", ctx.Run().WorkflowRunID, "count", len(orders))

	return map[string]any{"orders": orders}, nil
}

func main() {
	wf := df.NewWorkflow("order-pipeline", df.WorkflowOptions{})

	// Register ingestOrders as a root node (df.Root) with an explicit node key.
	wf.Node(ingestOrders, df.Root, df.NodeOptions{Key: "ingest_orders"})

	// Dispatches SDK commands: dumps the workflow graph during build,
	// and invokes the target node handler during execution.
	df.Main()
}
