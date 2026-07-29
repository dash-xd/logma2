variable "project_id" {
  description = "GCP project ID to deploy into."
  type        = string
}

variable "region" {
  description = "GCP region for the function and its source bucket."
  type        = string
}

variable "function_name" {
  description = "Name of the Cloud Function (also used as its source object prefix)."
  type        = string
}

variable "description" {
  description = "Human-readable description of the function."
  type        = string
  default     = "Deployed via Terraform composite action."
}

variable "entry_point" {
  description = "Exported Go function name Cloud Functions invokes, e.g. \"Function\"."
  type        = string
}

variable "source_dir" {
  description = "Path to the Go function source directory (relative to the caller's checkout) to zip and deploy."
  type        = string
}

variable "source_bucket_name" {
  description = "GCS bucket to hold the function's source archive. Created if it does not exist."
  type        = string
}

variable "runtime" {
  description = "Cloud Functions (1st gen) Go runtime identifier, e.g. go111, go116, go119, go121."
  type        = string
  default     = "go121"
}

variable "available_memory_mb" {
  description = "Memory allocated to the function, in MB. 128 is the Gen 1 minimum."
  type        = number
  default     = 128
}

variable "timeout_seconds" {
  description = "Function timeout, in seconds."
  type        = number
  default     = 60
}

variable "max_instances" {
  description = "Maximum concurrent instances."
  type        = number
  default     = 1
}

variable "allow_unauthenticated" {
  description = "Whether to allow unauthenticated (public) HTTP invocations."
  type        = bool
  default     = false
}

variable "environment_variables" {
  description = "Environment variables to set on the function."
  type        = map(string)
  default     = {}
}
