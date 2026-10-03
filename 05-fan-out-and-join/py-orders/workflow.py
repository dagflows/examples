from dataclasses import dataclass

from dagflows.authoring import Workflow
from dagflows.runtime import Ctx, Inputs

# A named workflow owns the workflow's settings. Exactly one project in the workspace names it.
wf = Workflow("order-pipeline")


# Dataclasses are reflected into the manifest as JSON Schema.
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
class EnrichedOrder(Order):
    tier: str


@dataclass
class EnrichedOrders:
    orders: list[EnrichedOrder]
    rejected: int


@dataclass
class Adjustment:
    order_id: int
    discount_cents: int


@dataclass
class PricedOrders:
    adjustments: list[Adjustment]


@dataclass
class Report:
    total_orders: int
    rejected: int
    gross_cents: int
    discount_cents: int
    net_cents: int
    by_tier: dict[str, int]


# Node keys default to the decorated function name unless overridden in @wf.node().
# The return annotation becomes the node's output schema.
@wf.node()
def ingest_orders(ctx: Ctx) -> Orders:
    orders = [
        Order(id=1001, customer_id="cus_alpha", amount_cents=1250, currency="usd"),
        Order(id=1002, customer_id="cus_beta", amount_cents=34000, currency="USD"),
        Order(id=1003, customer_id="cus_gamma", amount_cents=9999, currency="usd"),
        Order(id=1004, customer_id="cus_delta", amount_cents=500, currency="eur"),
        Order(id=1005, customer_id="cus_epsilon", amount_cents=0, currency="usd"),
    ]

    # ctx provides structured logging attached to this workflow run.
    ctx.log.info("run %s: ingested %d orders", ctx.run.workflow_run_id, len(orders))

    return Orders(orders=orders)


# external_node names a node of another project by its key. Its type is what this
# project expects of that node, checked against what that project declares at build time.
enriched = wf.external_node("enrich_customers", EnrichedOrders)
priced = wf.external_node("apply_pricing", PricedOrders)


# A node with several parents takes inputs. Each parent's output is reached by its
# handle and decoded into the handle's type.
@wf.node(depends=[enriched, priced])
def build_report(inputs: Inputs, ctx: Ctx) -> Report:
    orders = inputs[enriched].value()
    adjustments = inputs[priced].value()

    discounts = {adjustment.order_id: adjustment.discount_cents for adjustment in adjustments.adjustments}

    gross = sum(order.amount_cents for order in orders.orders)
    discount = sum(discounts.get(order.id, 0) for order in orders.orders)

    by_tier: dict[str, int] = {}
    for order in orders.orders:
        by_tier[order.tier] = by_tier.get(order.tier, 0) + 1

    report = Report(
        total_orders=len(orders.orders),
        rejected=orders.rejected,
        gross_cents=gross,
        discount_cents=discount,
        net_cents=gross - discount,
        by_tier=by_tier,
    )

    ctx.log.info(
        "report: %d orders, %d rejected, %d cents net of %d discount",
        report.total_orders,
        report.rejected,
        report.net_cents,
        report.discount_cents,
    )

    return report
