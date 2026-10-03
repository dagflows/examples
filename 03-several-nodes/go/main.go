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

type Report struct {
	TotalOrders int   `json:"total_orders"`
	Rejected    int   `json:"rejected"`
	GrossCents  int64 `json:"gross_cents"`
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

	ctx.Log().Info("validated orders", "kept", len(out.Orders), "rejected", out.Rejected)

	return out, nil
}

func buildReport(ctx *df.Ctx, in ValidatedOrders) (Report, error) {
	report := Report{TotalOrders: len(in.Orders), Rejected: in.Rejected}
	for _, order := range in.Orders {
		report.GrossCents += order.AmountCents
	}

	ctx.Log().Info("built report", "orders", report.TotalOrders, "rejected", report.Rejected, "gross_cents", report.GrossCents)

	return report, nil
}

func main() {
	wf := df.NewWorkflow("order-pipeline", df.WorkflowOptions{})

	// wf.Node returns the node's handle. Passing a handle as another node's edge makes
	// it that node's one parent, and the compiler checks the parent's Out is the child's In.
	ingested := wf.Node(ingestOrders, df.Root, df.NodeOptions{Key: "ingest_orders"})
	validated := wf.Node(validateOrders, ingested, df.NodeOptions{Key: "validate_orders"})
	wf.Node(buildReport, validated, df.NodeOptions{Key: "build_report"})

	// Dispatches SDK commands: dumps the workflow graph during build,
	// and invokes the target node handler during execution.
	df.Main()
}
