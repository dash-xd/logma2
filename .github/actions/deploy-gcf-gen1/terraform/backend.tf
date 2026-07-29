terraform {
  # Bucket and prefix are supplied at `terraform init` time via
  # -backend-config so each deployed function gets its own state path.
  # See action.yml.
  backend "gcs" {}
}
