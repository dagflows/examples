# Step 3: Several Nodes

The workflow from [step 2](../02-typed-node) grows to three typed nodes in a chain: `ingest_orders` → `validate_orders` → `build_report`. Each node receives its parent's output, decoded into its declared type.

Pick your language - each directory is a standalone runnable project:

- [Python](python)
- [Go](go)
- [Node](node) (TypeScript)

### What you'll learn
- How to declare a node's parent and receive its output as a typed input.
- How Dagflows runs nodes in dependency order and hands each output to the next node.
- How Dagflows checks each edge when it builds the workflow: the child's input schema must accept the parent's output schema.
