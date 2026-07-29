provider "google" {
  project = var.project_id
  region  = var.region
  # Authenticates via Application Default Credentials -- the calling
  # action sets GOOGLE_APPLICATION_CREDENTIALS to a service account JSON
  # key file before running terraform plan/apply.
}

data "archive_file" "source" {
  type        = "zip"
  source_dir  = var.source_dir
  output_path = "${path.module}/.build/${var.function_name}.zip"
}

resource "google_storage_bucket" "source" {
  name                        = var.source_bucket_name
  location                    = var.region
  uniform_bucket_level_access = true
  force_destroy               = true
}

resource "google_storage_bucket_object" "source" {
  # Content-hash suffix so a code change uploads a new object and
  # triggers a redeploy, without needing a random/versioning provider.
  name   = "${var.function_name}/${data.archive_file.source.output_md5}.zip"
  bucket = google_storage_bucket.source.name
  source = data.archive_file.source.output_path
}

resource "google_cloudfunctions_function" "function" {
  name        = var.function_name
  description = var.description
  region      = var.region
  runtime     = var.runtime

  available_memory_mb = var.available_memory_mb
  timeout             = var.timeout_seconds
  max_instances       = var.max_instances

  entry_point           = var.entry_point
  source_archive_bucket = google_storage_bucket.source.name
  source_archive_object = google_storage_bucket_object.source.name

  trigger_http                 = true
  https_trigger_security_level = "SECURE_ALWAYS"

  environment_variables = var.environment_variables
}

resource "google_cloudfunctions_function_iam_member" "invoker" {
  count = var.allow_unauthenticated ? 1 : 0

  project        = var.project_id
  region         = var.region
  cloud_function = google_cloudfunctions_function.function.name

  role   = "roles/cloudfunctions.invoker"
  member = "allUsers"
}
