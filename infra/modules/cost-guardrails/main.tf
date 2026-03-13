terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
}

locals {
  billing_alarm_enabled = var.billing_alarm_email != null && var.billing_alarm_email != ""
}

resource "aws_budgets_budget" "monthly_cost" {
  count = local.billing_alarm_enabled ? 1 : 0

  name         = "${var.resource_prefix}-monthly-budget"
  budget_type  = "COST"
  limit_amount = tostring(var.billing_budget_limit)
  limit_unit   = "USD"
  time_unit    = "MONTHLY"

  notification {
    comparison_operator = "GREATER_THAN"
    threshold           = 90
    threshold_type      = "PERCENTAGE"
    notification_type   = "ACTUAL"

    subscriber_email_addresses = [var.billing_alarm_email]
  }
}
