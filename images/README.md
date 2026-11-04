# Screenshots

The examples' READMEs show these screenshots of the Dagflows dashboard. Take them on
app.dagflows.com with the examples deployed, all at the same window width and in the same theme,
cropped to the dialog or panel they show. Blur anything personal or secret: emails, signing secrets
and the token at the end of a webhook's URL.

| File | Where | What it shows |
|---|---|---|
| `new-project.png` | Projects → **New project** | The Create project dialog, named `Examples` |
| `new-workflow.png` | The project → **New workflow** | The Create workflow dialog: `order-pipeline`, the repository's URL, branch `main` |
| `deploy.png` | The workflow → **Deploy** | The Deploy workflow dialog, **Commit hash** left empty |
| `build.png` | The workflow → **Builds** | A finished build of step 1 |
| `run.png` | The workflow → **Runs** → **Run** | The Run workflow dialog, the payload box empty |
| `run-detail.png` | A run of step 1 | `ingest_orders` succeeded, its log reading `ingested 5 orders` |
| `02-output.png` | A run of step 2 | The output of `ingest_orders`, the five orders |
| `03-run.png` | A run of step 3 | The three nodes, in order, succeeded |
| `04-run.png` | A run of step 4 | The three nodes, `validate_orders` in Go |
| `05-run.png` | A run of step 5 | The six nodes, the two branches and the join |
| `06-run-payload.png` | Step 6 → **Runs** → **Run** | The Run workflow dialog with `orders.json` as the payload |
| `06-run.png` | A run of step 6 | Started manually, `ingest_orders` logging `manual event ... received at ...` |
| `07-connect.png` | Step 7 → **Triggers** → **Connect** on `orders_received` | The dialog with **Create a webhook for this** and **Require a signature** on |
| `07-secret.png` | Right after connecting | The Save your signing secret dialog, the secret and the URL's token blurred |
| `07-test-event.png` | **Event sources** → `orders_received` → **Send test event** | The Send a test event dialog with the orders |
| `07-run.png` | A run started by the signed request | Started by a webhook, `ingest_orders` logging `webhook event ... received at ...` |
| `08-variables.png` | Step 8's workflow → **Settings** → **ENVs & Secrets** | `DISCOUNT_CENTS` `500` and `GOLD_TIER_CENTS` `10000` |
| `08-secret.png` | The project → **Settings** → **ENVs & Secrets** | `REPORT_SIGNING_KEY` under Secrets |
| `08-run.png` | A run of step 8 | The logs with `discount_cents=500`, `gold from 10000 cents` and `signed the report` |
