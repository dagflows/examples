"""The trigger that starts a run, and the initial and terminal nodes (ingest_orders and build_report).

Owns global workflow settings for the multi-project workspace by instantiating a named Workflow.
Secondary projects contribute nodes via unnamed workflows.
"""

from dataclasses import dataclass, field

from dagflows.authoring import Workflow
from dagflows.runtime import Inputs

wf = Workflow("order-pipeline", max_concurrent_nodes=4)


# --------------------------------------------------------------------------
# Dataclasses reflected as JSON Schema in manifests for pre-run type validation across edges.
# --------------------------------------------------------------------------


@dataclass
class Order:
    id: int
    customer_id: str
    amount_cents: int
    currency: str


@dataclass
class Orders:
    orders: list[Order]


@dataclass
class Report:
    total_orders: int
    gross_cents: int
    discount_cents: int
    net_cents: int
    # Tier name -> how many orders landed in it.
    by_tier: dict[str, int] = field(default_factory=dict)


# --------------------------------------------------------------------------
# A. The root.
# --------------------------------------------------------------------------


# A run starts with a batch of orders as its event. A manual run, or a webhook bound
# to this trigger, delivers the batch, and the platform checks it against Orders first.
orders_received = wf.trigger("orders_received", event=Orders)


@wf.node(depends=[orders_received])
def ingest_orders(batch: Orders) -> Orders:
    """Root node handing the delivered batch of orders to both branches.

    Its one parent is the trigger, so `batch` is the event, decoded into Orders.
    """
    return batch


# --------------------------------------------------------------------------
# F. Where the two branches meet.
# --------------------------------------------------------------------------


@wf.node(
    depends=[
        # References external project nodes by key (TypeScript and Go), resolved and validated at build time.
        wf.external_node("enrich_customers"),
        wf.external_node("apply_pricing"),
    ]
)
def build_report(inputs: Inputs) -> Report:
    """Combines outputs from the enrichment and pricing branches into a final report.

    Uses `inputs` to access multiple parent outputs by node key, decoded from JSON into dictionaries.
    """
    enriched = inputs["enrich_customers"].value()
    priced = inputs["apply_pricing"].value()

    # order id -> discount, so the two branches can be joined by order.
    discounts = {
        int(adjustment["order_id"]): int(adjustment["discount_cents"])
        for adjustment in priced["adjustments"]
    }

    by_tier: dict[str, int] = {}
    gross = 0
    discount_total = 0

    for order in enriched["orders"]:
        gross += int(order["amount_cents"])
        # Applies zero discount if an order was not scored on the risk branch.
        discount_total += discounts.get(int(order["id"]), 0)
        tier = order["tier"]
        by_tier[tier] = by_tier.get(tier, 0) + 1

    return Report(
        total_orders=len(enriched["orders"]),
        gross_cents=gross,
        discount_cents=discount_total,
        net_cents=gross - discount_total,
        by_tier=by_tier,
    )
