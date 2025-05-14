variable "project_id" {
  description = "The Google Cloud project ID"
  type        = string
}

variable "region" {
  description = "The Google Cloud region"
  type        = string
  default     = "europe-central2"
}

variable "telegram_token" {
  description = "Telegram Bot API token"
  type        = string
  sensitive   = true
}

variable "telegram_allowed_chat_ids" {
  description = "Comma-separated list of allowed Telegram chat IDs"
  type        = string
}

variable "google_sheets_spreadsheet_id" {
  description = "Google Sheets spreadsheet ID"
  type        = string
}

variable "google_sheets_sheet_name" {
  description = "Google Sheets sheet name"
  type        = string
}

variable "openexchangerates_app_id" {
  description = "OpenExchangeRates API app ID"
  type        = string
  sensitive   = true
}

variable "default_currency" {
  description = "Default currency code"
  type        = string
}

variable "webhook_url" {
  description = "Webhook URL"
  type        = string
}

variable "image_sha256" {
  description = "Image SHA256 hash"
  type        = string
}
