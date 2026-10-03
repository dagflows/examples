# Step 2: Typed Node

The workflow from [step 1](../01-first-node) with its root node (`ingest_orders`) typed: the same mock orders, now returned as a declared type that the SDK publishes in the workflow manifest.

Pick your language - each directory is a standalone runnable project:

- [Python](python)
- [Go](go)
- [Node](node) (TypeScript)

### What you'll learn
- How to declare a node's output type (dataclasses, structs, TypeScript interfaces).
- How the SDK reflects that type into the manifest as JSON Schema.
- Why it matters: from step 3, Dagflows checks each edge between typed nodes against these schemas when it builds the workflow.
