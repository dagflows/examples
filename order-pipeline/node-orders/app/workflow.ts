/**
 * Customer enrichment and risk scoring nodes in TypeScript.
 * Contributes nodes to the workflow owned by py-orders, and all handlers must be exported.
 */

import { Workflow } from "@dagflows/sdk/authoring";

export const wf = new Workflow();

// ---------------------------------------------------------------------------
// Type definitions mirroring Python and Go schemas, validated across edges at build time.
// ---------------------------------------------------------------------------

/**
 * Order representation matching schemas in Python and Go.
 *
 * id uses bigint so the SDK emits an int64 JSON Schema integer. JSON Schema
 * allows an integer to satisfy a number, but never a number to satisfy an integer.
 * Using number here fails validation on the score_risk (TypeScript) to apply_pricing (Go)
 * edge where Go requires an int64 integer.
 */
interface Order {
  id: bigint;
  customer_id: string;
  // A number is fine here: amount_cents only ever travels go -> ts -> python,
  // and an integer satisfies a number in that direction.
  amount_cents: number;
  currency: string;
}

/** Raw order batch emitted by ingest_orders. */
interface Orders {
  orders: Order[];
}

/** Validated order batch emitted by validate_orders. */
interface ValidatedOrders {
  orders: Order[];
  rejected: number;
}

/** Order extended with customer tier classification. */
interface EnrichedOrder extends Order {
  tier: string;
}

/** Enriched order batch passed to build_report. */
interface EnrichedOrders {
  orders: EnrichedOrder[];
}

/** Risk assessment scores passed to apply_pricing. */
interface RiskScores {
  // bigint, because Go reads this as int64. See the note on Order.id.
  scores: { order_id: bigint; score: number }[];
}

// ---------------------------------------------------------------------------
// Declares external parent nodes by key with expected schemas for build-time validation.
// ---------------------------------------------------------------------------

const validateOrders = wf.externalNode<ValidatedOrders>("validate_orders");
const ingestOrders = wf.externalNode<Orders>("ingest_orders");

// ---------------------------------------------------------------------------
// C. Enrichment, on the branch that runs py -> go -> ts -> py.
// ---------------------------------------------------------------------------

/**
 * Classifies customer tiers for validated orders based on amount.
 */
export const enrichCustomers = wf.node(
  function enrichCustomers({ input }: { input: ValidatedOrders }): EnrichedOrders {
    return {
      orders: input.orders.map((order) => ({
        ...order,
        tier: order.amount_cents >= 100_00 ? "gold" : "standard",
      })),
    };
  },
  // Explicit key overrides function name with snake_case convention used across projects.
  { key: "enrich_customers", depends: [validateOrders] },
);

// ---------------------------------------------------------------------------
// D. Risk scoring, on the branch that runs py -> ts -> go -> py.
// ---------------------------------------------------------------------------

/**
 * Calculates risk scores directly from raw ingested orders for parallel branch processing.
 */
export const scoreRisk = wf.node(
  function scoreRisk({ input }: { input: Orders }): RiskScores {
    return {
      scores: input.orders.map((order) => ({
        // BigInt() rather than passing order.id through: the codec revives an
        // integer as a bigint only BEYOND the safe range, so a small id like
        // 1001 arrives as a plain number however the type reads. Converting
        // explicitly is right for either, and what leaves is the exact integer
        // Go decodes into an int64.
        order_id: BigInt(order.id),
        // Stable, spread across the range, and no floating point surprises:
        // the last digit of the id over ten.
        score: Number(BigInt(order.id) % 10n) / 10,
      })),
    };
  },
  // Explicit key matching external node references in go-orders.
  { key: "score_risk", depends: [ingestOrders] },
);
