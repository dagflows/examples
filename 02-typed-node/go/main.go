package main

import df "github.com/dagflows/sdk-go"

// Structs and their json tags are reflected into the manifest as JSON Schema.
type Order struct {
	ID          int64  `json:"id"`
	CustomerID  string `json:"customer_id"`
	AmountCents int64  `json:"amount_cents"`
	Currency    string `json:"currency"`
}

type Orders struct {
	Orders []Order `json:"orders"`
}

// Node handlers follow the signature func(ctx *df.Ctx, in In) (Out, error).
// Root nodes take df.None since they have no parent inputs. Out becomes the
// node's output schema.
func ingestOrders(ctx *df.Ctx, _ df.None) (Orders, error) {
	orders := []Order{
		{ID: 1001, CustomerID: "cus_alpha", AmountCents: 1250, Currency: "usd"},
		{ID: 1002, CustomerID: "cus_beta", AmountCents: 34000, Currency: "USD"},
		{ID: 1003, CustomerID: "cus_gamma", AmountCents: 9999, Currency: "usd"},
		{ID: 1004, CustomerID: "cus_delta", AmountCents: 500, Currency: "eur"},
		{ID: 1005, CustomerID: "cus_epsilon", AmountCents: 0, Currency: "usd"},
	}

	// ctx.Log() is a *slog.Logger attached to the workflow run's log stream.
	ctx.Log().Info("ingested orders", "run", ctx.Run().WorkflowRunID, "count", len(orders))

	return Orders{Orders: orders}, nil
}

func main() {
	wf := df.NewWorkflow("order-pipeline", df.WorkflowOptions{})

	// Register ingestOrders as a root node (df.Root) with an explicit node key.
	wf.Node(ingestOrders, df.Root, df.NodeOptions{Key: "ingest_orders"})

	// Dispatches SDK commands: dumps the workflow graph during build,
	// and invokes the target node handler during execution.
	df.Main()
}
