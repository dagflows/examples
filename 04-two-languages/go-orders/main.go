package main

import (
	"strings"

	df "github.com/dagflows/sdk-go"
)

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

type ValidatedOrders struct {
	Orders   []Order `json:"orders"`
	Rejected int     `json:"rejected"`
}

// A node with one parent receives that parent's output, decoded into In.
func validateOrders(ctx *df.Ctx, in Orders) (ValidatedOrders, error) {
	out := ValidatedOrders{Orders: make([]Order, 0, len(in.Orders))}
	for _, order := range in.Orders {
		if !strings.EqualFold(order.Currency, "usd") || order.AmountCents <= 0 {
			out.Rejected++
			continue
		}

		order.Currency = "usd"
		out.Orders = append(out.Orders, order)
	}

	// ctx.Log() is a *slog.Logger attached to the workflow run's log stream.
	ctx.Log().Info("validated orders", "kept", len(out.Orders), "rejected", out.Rejected)

	return out, nil
}

func main() {
	// A project contributing nodes to a workflow another project owns uses an unnamed workflow.
	wf := df.NewWorkflow("", df.WorkflowOptions{})

	// wf.External names a node of another project by its key, typed with what this node
	// expects of it, which Dagflows checks against what that project declares at build time.
	wf.Node(validateOrders, wf.External[Orders]("ingest_orders"), df.NodeOptions{Key: "validate_orders"})

	// Dispatches SDK commands: dumps the workflow graph during build,
	// and invokes the target node handler during execution.
	df.Main()
}
