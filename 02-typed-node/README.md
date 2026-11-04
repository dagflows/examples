# Step 2: Typed Node

The node from [step 1](../01-first-node), with its output declared as a type. The orders are the
same, but Dagflows now knows their shape: when it builds the workflow, it reads the type from your
code and records it as the node's output.

Pick a language. Each folder is a complete project:

- [python](python): `workflow.py` and `requirements.txt`
- [go](go): `main.go` and `go.mod`
- [node](node), in TypeScript: `workflow.ts`, `dagflows.yaml` and `package.json`

## The Code

Each language declares an `Order` and an `Orders` type, and the node returns an `Orders`.

**Python**, [workflow.py](python/workflow.py), with dataclasses:

```python
@dataclass
class Order:
    id: int
    customer_id: str
    amount_cents: int
    currency: str


@dataclass
class Orders:
    orders: list[Order]


@wf.node()
def ingest_orders(ctx: Ctx) -> Orders:
    orders = [Order(id=1001, customer_id="cus_alpha", amount_cents=1250, currency="usd"), ...]
    return Orders(orders=orders)
```

**Go**, [main.go](go/main.go), with structs and their `json` tags:

```go
type Order struct {
	ID          int64  `json:"id"`
	CustomerID  string `json:"customer_id"`
	AmountCents int64  `json:"amount_cents"`
	Currency    string `json:"currency"`
}

type Orders struct {
	Orders []Order `json:"orders"`
}

func ingestOrders(ctx *df.Ctx, _ df.None) (Orders, error) {
	// ...
	return Orders{Orders: orders}, nil
}
```

**TypeScript**, [workflow.ts](node/workflow.ts), with interfaces:

```ts
interface Order {
  id: number;
  customer_id: string;
  amount_cents: number;
  currency: string;
}

interface Orders {
  orders: Order[];
}

export const ingestOrders = wf.node(
  function ingestOrders({ ctx }): Orders {
    // ...
    return { orders };
  },
  { key: "ingest_orders" },
);
```

- **The return type is the node's output type.** Python reads the return annotation, Go the
  function's result, TypeScript the declared return type.
- **Field names are what goes over the wire.** In Go, the `json` tags name the fields, so the
  output reads `customer_id` in every language.
- **Why it matters:** from step 3 on, one node's output is another's input, and Dagflows checks
  that the two types agree before the workflow ever runs.
- **TypeScript** needs two things JavaScript did not: [dagflows.yaml](node/dagflows.yaml) names
  `workflow.ts` as the project's entry, and `typescript` in `package.json` lets Dagflows read the
  interfaces.

## Deploy It

Follow [Deploy an Example on Dagflows](../README.md#deploy-an-example-on-dagflows) with the folder
of your language. If you deployed step 1, replace the repository's files with this step's, push,
and click **Deploy**. The run takes no payload.

## What You'll See

The run logs the same line as in step 1. Open `ingest_orders` to see its output, the five orders
as an `Orders`:

![The output of ingest_orders](../images/02-output.png)

## Try It on Your Machine

Building the workflow's manifest shows what Dagflows reads from your code, the node and its types:

**Python**:

```bash
python -m venv .venv && source .venv/bin/activate
pip install -r requirements.txt
python -m dagflows build manifest workflow -o dagflows-manifest.json
```

**Go**:

```bash
go run . build manifest -o dagflows-manifest.json
```

**TypeScript**:

```bash
npm ci
npx dagflows-sdk build manifest workflow.ts -o dagflows-manifest.json
```

In `dagflows-manifest.json`, the node's `io.output.schema` is the `Orders` type as JSON Schema: an
object with an `orders` array, each item an `Order` with its four fields. To run the node, use the
commands of [step 1](../01-first-node#try-it-on-your-machine), with `workflow.ts:ingestOrders` in
TypeScript.

## Next

[Step 3](../03-several-nodes) adds two nodes after this one, each fed the output of the one before.
