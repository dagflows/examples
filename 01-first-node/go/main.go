// A workflow with one node, run by hand.
package main

import df "github.com/dagflows/sdk-go"

// ingestOrders returns a batch of orders as plain data, and logs how many there are.
// ctx is the node's context: ctx.Log() writes to the run's logs, and ctx.Run() says
// which run this is. df.None says the node takes no input.
func ingestOrders(ctx *df.Ctx, _ df.None) (map[string]any, error) {
	orders := []map[string]any{
		{"id": 1001, "customer_id": "cus_alpha", "amount_cents": 1250, "currency": "usd"},
		{"id": 1002, "customer_id": "cus_beta", "amount_cents": 34000, "currency": "USD"},
		{"id": 1003, "customer_id": "cus_gamma", "amount_cents": 9999, "currency": "usd"},
		{"id": 1004, "customer_id": "cus_delta", "amount_cents": 500, "currency": "eur"},
		{"id": 1005, "customer_id": "cus_epsilon", "amount_cents": 0, "currency": "usd"},
	}

	ctx.Log().Info("ingested orders", "run", ctx.Run().WorkflowRunID, "count", len(orders))
	return map[string]any{"orders": orders}, nil
}

func main() {
	// The workflow every step of these examples builds on.
	wf := df.NewWorkflow("order-pipeline", df.WorkflowOptions{})

	// df.Root makes it a node with no parents.
	wf.Node(ingestOrders, df.Root, df.NodeOptions{Key: "ingest_orders"})

	// Answers the platform: emits the manifest when building, runs the node when invoked.
	df.Main()
}
