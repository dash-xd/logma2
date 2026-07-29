output "function_name" {
  description = "Name of the deployed function."
  value       = google_cloudfunctions_function.function.name
}

output "function_url" {
  description = "HTTPS trigger URL of the deployed function."
  value       = google_cloudfunctions_function.function.https_trigger_url
}
