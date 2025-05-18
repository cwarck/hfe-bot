locals {
  name                    = "hfe-bot"
  github_repository_name  = "cwarck/hfe-bot"
  github_actions_identity = "principalSet://iam.googleapis.com/projects/${data.google_project.this.number}/locations/global/workloadIdentityPools/github/attribute.repository/${local.github_repository_name}"
}

data "google_project" "this" {}

resource "google_project_service" "required_apis" {
  for_each = toset([
    "artifactregistry.googleapis.com",
    "compute.googleapis.com",
    "drive.googleapis.com",
    "iam.googleapis.com",
    "run.googleapis.com",
    "sheets.googleapis.com",
  ])

  project = var.project_id
  service = each.key

  disable_dependent_services = false
  disable_on_destroy         = false
}

# Add this service account to the target Google Sheets spreadsheet with the "Editor" role
resource "google_service_account" "app" {
  account_id   = local.name
  display_name = "Service Account for ${local.name}"
  project      = var.project_id
}

resource "google_service_account_iam_binding" "iam_serviceaccount_user" {
  service_account_id = google_service_account.app.name
  role               = "roles/iam.serviceAccountUser"
  members            = [local.github_actions_identity]
}

resource "google_artifact_registry_repository" "app" {
  location      = var.region
  repository_id = local.name
  description   = "Docker repository for ${local.name}"
  format        = "DOCKER"
  project       = var.project_id

  depends_on = [google_project_service.required_apis]
}

# Allow to push from Github Actions using Workload Identity Federation
resource "google_artifact_registry_repository_iam_binding" "binding" {
  project    = google_artifact_registry_repository.app.project
  location   = google_artifact_registry_repository.app.location
  repository = google_artifact_registry_repository.app.name
  role       = "roles/artifactregistry.writer"
  members    = [local.github_actions_identity]
}

resource "google_cloud_run_v2_service" "app" {
  name                = local.name
  location            = var.region
  deletion_protection = false
  ingress             = "INGRESS_TRAFFIC_ALL"

  template {
    containers {
      image = "${google_artifact_registry_repository.app.location}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.app.name}/hfe-bot@sha256:${var.image_sha256}"

      env {
        name  = "TELEGRAM_TOKEN"
        value = var.telegram_token
      }

      env {
        name  = "TELEGRAM_ALLOWED_CHAT_IDS"
        value = var.telegram_allowed_chat_ids
      }

      env {
        name  = "GOOGLE_SHEETS_SPREADSHEET_ID"
        value = var.google_sheets_spreadsheet_id
      }

      env {
        name  = "GOOGLE_SHEETS_SHEET_NAME"
        value = var.google_sheets_sheet_name
      }

      env {
        name  = "OPENEXCHANGERATES_APP_ID"
        value = var.openexchangerates_app_id
      }

      env {
        name  = "DEFAULT_CURRENCY"
        value = var.default_currency
      }

      env {
        name  = "WEBHOOK_URL"
        value = var.webhook_url
      }
    }

    service_account = google_service_account.app.email
  }
}

# Disable authentication for the Cloud Run service
resource "google_cloud_run_v2_service_iam_binding" "noauth" {
  project  = google_cloud_run_v2_service.app.project
  location = google_cloud_run_v2_service.app.location
  name     = google_cloud_run_v2_service.app.name
  role     = "roles/run.invoker"
  members  = ["allUsers"]
}

# Allow to uptdate the Cloud Run service from Github Actions
resource "google_cloud_run_v2_service_iam_binding" "admin" {
  project  = google_cloud_run_v2_service.app.project
  location = google_cloud_run_v2_service.app.location
  name     = google_cloud_run_v2_service.app.name
  role     = "roles/run.admin"
  members  = [local.github_actions_identity]
}
