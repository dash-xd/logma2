# deploy-gcf-gen1

Composite action that deploys a Go-runtime Google Cloud Function (1st
gen) via Terraform, at minimum spec (128MB memory, single instance,
HTTP trigger).

It is meant to be called from a **separate, private repository** that
owns the function's source code:

```yaml
- uses: dash-xd/logma2/.github/actions/deploy-gcf-gen1@main
  with:
    gcp_project_id: my-project
    ...
```

See `examples/workflows/deploy-cloud-function.yml` in this repo for a
full sample caller workflow (placeholder values only).

## Authentication

This action does not decrypt any secrets itself. The calling workflow
is expected to supply the contents of a GCP service account JSON key
via the `gcp_credentials_json` input, already decrypted -- e.g. through
a sops-age step in the caller's own bootstrap (as done in production in
`huram-abi`).

## Inputs

| Name | Required | Default | Description |
| --- | --- | --- | --- |
| `gcp_project_id` | yes | | Target GCP project ID |
| `gcp_region` | yes | | Region for the function and its source bucket |
| `gcp_credentials_json` | yes | | Decrypted service account JSON key contents |
| `function_name` | yes | | Cloud Function name |
| `entry_point` | yes | | Exported Go function name to invoke |
| `source_dir` | yes | | Path to the Go source dir, relative to the caller's checkout |
| `source_bucket_name` | yes | | GCS bucket for the source archive (created if missing) |
| `tf_state_bucket` | yes | | GCS bucket for Terraform remote state |
| `tf_state_prefix` | no | `cloud-functions` | State path prefix, combined with `function_name` |
| `runtime` | no | `go121` | Cloud Functions (1st gen) Go runtime identifier |
| `available_memory_mb` | no | `128` | Memory in MB (128 is the Gen 1 minimum) |
| `timeout_seconds` | no | `60` | Function timeout |
| `max_instances` | no | `1` | Maximum concurrent instances |
| `allow_unauthenticated` | no | `false` | Allow public HTTP invocation |
| `environment_variables` | no | `{}` | JSON object of env vars for the function |
| `terraform_action` | no | `apply` | `apply` or `plan` (plan-only, no changes applied) |

## Outputs

| Name | Description |
| --- | --- |
| `function_url` | HTTPS trigger URL of the deployed function |
