# Step 8: Variables and Secrets

The workflow from [step 7](../07-webhook), with some values moved out of the code and onto
Dagflows. The discount and the amount from which an order is gold become variables, which you
change without a deploy. And `build_report` signs the report with a key, a secret that never goes
into the repository.

```
               orders_received (trigger)
                          |
                  ingest_orders (Python)
                   /                  \
    validate_orders (Go)        score_risk (TypeScript)
            |                            |
   enrich_customers (TypeScript)   apply_pricing (Go)
     GOLD_TIER_CENTS                DISCOUNT_CENTS
                   \                  /
                   build_report (Python)
                    REPORT_SIGNING_KEY
```

This folder is the whole example. Copy all of it:

```
dagflows.yaml          the workspace: the workflow's name and its three projects
orders.json            the five orders, a run's payload
py-orders/             Python: the trigger, ingest_orders and build_report, which signs
go-orders/             Go: validate_orders and apply_pricing, which reads DISCOUNT_CENTS
node-orders/           TypeScript: enrich_customers, which reads GOLD_TIER_CENTS, and score_risk
```

## The Code

A node names the variables it needs, and reads them from its environment when it runs. Dagflows
holds their values and gives each node only the ones it names.

**Go**, `apply_pricing` in [go-orders/main.go](go-orders/main.go):

```go
discountCents, err := strconv.ParseInt(os.Getenv("DISCOUNT_CENTS"), 10, 64)
if err != nil || discountCents < 0 {
	return PricedOrders{}, fmt.Errorf("DISCOUNT_CENTS is %q, want a whole number of cents", os.Getenv("DISCOUNT_CENTS"))
}

// ...

wf.Node(applyPricing, wf.External[RiskScores]("score_risk"),
	df.NodeOptions{Key: "apply_pricing", Env: []string{"DISCOUNT_CENTS"}})
```

**TypeScript**, `enrich_customers` in [node-orders/workflow.ts](node-orders/workflow.ts):

```ts
const goldTier = process.env.GOLD_TIER_CENTS ?? "";

// ...

{ key: "enrich_customers", depends: [validateOrders], env: ["GOLD_TIER_CENTS"] },
```

**Python**, `build_report` in [py-orders/workflow.py](py-orders/workflow.py), signs the report:

```python
@wf.node(depends=[enriched, priced], env=["REPORT_SIGNING_KEY"])
def build_report(inputs: Inputs, ctx: Ctx) -> SignedReport:
    # ...
    key = os.environ.get("REPORT_SIGNING_KEY", "")
    body = json.dumps(asdict(report), sort_keys=True, separators=(",", ":"))
    signature = hmac.new(key.encode(), body.encode(), hashlib.sha256).hexdigest()

    return SignedReport(report=report, signature=signature)
```

- **`env` names a node's variables**: `env=[...]` in Python, `Env: []string{...}` in Go,
  `env: [...]` in TypeScript. A node gets those, as environment variables, and no others.
- **A variable is read inside the node**, while it runs. Dagflows builds the workflow without any.
- **A value is set on the workflow, its project or its organization.** A node gets the workflow's
  own value of a name, else its project's, else its organization's.
- **A secret** is a variable whose value is never shown again once saved.

## Deploy It

Set the variables first: a deploy fails while a node names a variable that has no value, and says
which, for example `node "apply_pricing" names DISCOUNT_CENTS`.

1. **The two values, on the workflow.** Open the workflow's **Settings → ENVs & Secrets**. Under
   **Environment Variables**, click **New variable** and add `DISCOUNT_CENTS` = `500`, then
   `GOLD_TIER_CENTS` = `10000`:

   ![The workflow's two variables](../images/08-variables.png)

2. **The key, on the project.** Open the project's **Settings → ENVs & Secrets**. Under **Secrets**,
   click **New secret** and add `REPORT_SIGNING_KEY` with a long random value, such as the output of
   `openssl rand -hex 32`. Set on the project, the key signs for every workflow of the project that
   names it.

   ![The project's secret](../images/08-secret.png)

3. **Deploy and run.** Follow [Deploy an Example on Dagflows](../README.md#deploy-an-example-on-dagflows)
   with this whole folder, or replace your repository's files with it, push, and click **Deploy**.
   Then run it with [orders.json](orders.json) as the payload, as in
   [step 6](../06-trigger#deploy-it), or send the orders to the webhook of step 7.

## What You'll See

`enrich_customers` and `apply_pricing` log the values they were given, and `build_report` the
report's signature:

![The run with the variables and the signature in its logs](../images/08-run.png)

```
INFO manual event <event id> received at <time>
INFO run <run id>: ingested 5 orders
time=<time> level=INFO msg="validated orders" kept=3 rejected=2
INFO enriched 3 orders, gold from 10000 cents node=enrich_customers
INFO scored 5 orders node=score_risk
time=<time> level=INFO msg="priced orders" scored=5 discounted=3 discount_cents=500
INFO report: 3 orders, 2 rejected, 43749 cents net of 1500 discount
INFO signed the report: <signature>
```

`build_report`'s output is the report and its signature. Whoever receives the report checks the
signature with the same key.

### Change a Value Without a Deploy

A node reads its variables when it starts, so a new value reaches the next run. On the workflow's
**ENVs & Secrets**, set `DISCOUNT_CENTS` to `700` and run again. `apply_pricing` logs
`discount_cents=700`, and the report reads `43149 cents net of 2100 discount`. A node already
running keeps the values it started with.

### Keep a Secret Secret

- Wherever a node that gets a secret prints it, in its log or in an error, Dagflows replaces it with
  `***`. A value shorter than 8 bytes is not replaced, since it would hide ordinary text.
- That catches accidents only. A node can still print a secret in another form, so print none.
- A secret is stored encrypted, and travels encrypted to the machine that runs the node.

## Try It on Your Machine

Locally, you set the variables yourself. From `go-orders`:

```bash
DISCOUNT_CENTS=700 go run . dev run apply_pricing --input score_risk='{"scores": [{"order_id": 1001, "score": 0.1}, {"order_id": 1004, "score": 0.4}]}'
```

Order 1001 gets the 700 cent discount. Without `DISCOUNT_CENTS`, the node fails and says what is
missing: `DISCOUNT_CENTS is "", want a whole number of cents`.

## Where to Go From Here

That is the whole series: typed nodes in three languages, joined into one workflow, started by a
webhook and configured from Dagflows. From here, change the example into your own workflow, or start
a new one from the step closest to what you need.
