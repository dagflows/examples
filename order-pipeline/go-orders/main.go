// Go pipeline nodes for validation and pricing.
// Emits manifest metadata via df.Main() and executes nodes when invoked by the worker.
// Uses an unnamed workflow instance because py-orders owns global workflow settings.
package main

import df "github.com/dagflows/sdk-go"

// --------------------------------------------------------------------------
// Struct schemas matching Python and TypeScript payloads, validated at build time.
// --------------------------------------------------------------------------

// Order represents an order emitted by ingest_orders (Python).
type Order struct {
	ID          int64  `json:"id"`
	CustomerID  string `json:"customer_id"`
	AmountCents int64  `json:"amount_cents"`
	Currency    string `json:"currency"`
}

// Orders represents a batch of orders from ingest_orders.
type Orders struct {
	Orders []Order `json:"orders"`
}

// ValidatedOrders represents validated orders passed to enrich_customers.
type ValidatedOrders struct {
	Orders   []Order `json:"orders"`
	Rejected int     `json:"rejected"`
}

// RiskScore represents a risk assessment emitted by score_risk (TypeScript).
type RiskScore struct {
	OrderID int64   `json:"order_id"`
	Score   float64 `json:"score"`
}

// RiskScores represents a batch of risk scores from score_risk.
type RiskScores struct {
	Scores []RiskScore `json:"scores"`
}

// Adjustment represents an order discount amount.
type Adjustment struct {
	OrderID       int64 `json:"order_id"`
	DiscountCents int64 `json:"discount_cents"`
}

// PricedOrders represents pricing adjustments passed to build_report.
type PricedOrders struct {
	Adjustments []Adjustment `json:"adjustments"`
}

// --------------------------------------------------------------------------
// B. Validation, on the branch that runs py -> go -> ts -> py.
// --------------------------------------------------------------------------

// validateOrders drops invalid currency or non-positive amounts and normalizes currency code.
// The input parameter automatically receives decoded data validated against the parent node schema.
func validateOrders(_ *df.Ctx, in Orders) (ValidatedOrders, error) {
	out := ValidatedOrders{Orders: make([]Order, 0, len(in.Orders))}

	for _, order := range in.Orders {
		// Compares case-insensitively and normalizes currency to lowercase usd.
		normalised := order
		normalised.Currency = "usd"

		if !equalFold(order.Currency, "usd") || order.AmountCents <= 0 {
			out.Rejected++
			continue
		}

		out.Orders = append(out.Orders, normalised)
	}

	return out, nil
}

// --------------------------------------------------------------------------
// E. Pricing, on the branch that runs py -> ts -> go -> py.
// --------------------------------------------------------------------------

// applyPricing computes order discounts based on risk scores from the parallel risk branch.
func applyPricing(_ *df.Ctx, in RiskScores) (PricedOrders, error) {
	out := PricedOrders{Adjustments: make([]Adjustment, 0, len(in.Scores))}

	for _, scored := range in.Scores {
		// Orders with low risk scores (< 0.34) receive a 500 cent discount.
		discount := int64(0)
		if scored.Score < 0.34 {
			discount = 500
		}

		out.Adjustments = append(out.Adjustments, Adjustment{
			OrderID:       scored.OrderID,
			DiscountCents: discount,
		})
	}

	return out, nil
}

// equalFold provides zero-dependency case-insensitive ASCII comparison.
func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range len(a) {
		x, y := a[i], b[i]
		if 'A' <= x && x <= 'Z' {
			x += 'a' - 'A'
		}
		if 'A' <= y && y <= 'Z' {
			y += 'a' - 'A'
		}
		if x != y {
			return false
		}
	}

	return true
}

func main() {
	// Unnamed workflow instance contributing nodes without overriding global settings.
	wf := df.NewWorkflow("", df.WorkflowOptions{})

	// External parent nodes referenced by key and validated against producer schemas.
	wf.Node(validateOrders,
		wf.External[Orders]("ingest_orders"),
		df.NodeOptions{Key: "validate_orders"},
	)

	wf.Node(
		applyPricing,
		wf.External[RiskScores]("score_risk"),
		df.NodeOptions{Key: "apply_pricing"},
	)

	// Handles manifest generation or dispatches node execution based on runtime environment.
	df.Main()
}
