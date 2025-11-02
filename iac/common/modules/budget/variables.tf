variable "budget_name" {
  description = "Name of the budget"
  type        = string
}

variable "monthly_budget_limit" {
  description = "Monthly budget limit in USD"
  type        = string
}

variable "alert_emails" {
  description = "List of email addresses to receive budget alerts"
  type        = list(string)
}

variable "cost_filters" {
  description = "Cost filters to track specific resources (e.g., by tags)"
  type        = map(list(string))
  default     = {}

  # Example:
  # cost_filters = {
  #   "tag:App" = ["blog"]
  #   "tag:Environment" = ["prod"]
  # }
}

variable "time_period_start" {
  description = "Budget start date (format: YYYY-MM-DD_HH:MM)"
  type        = string
  default     = "2025-01-01_00:00"
}

variable "time_period_end" {
  description = "Budget end date (format: YYYY-MM-DD_HH:MM)"
  type        = string
  default     = "2030-12-31_23:59"
}

variable "tags" {
  description = "Additional tags for the budget"
  type        = map(string)
  default     = {}
}