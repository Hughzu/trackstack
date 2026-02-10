variable "calories_url" {
  type    = string
  default = getenv("TURSO_CALORIES_URL")
}

variable "calories_token" {
  type    = string
  default = getenv("TURSO_CALORIES_TOKEN")
}

variable "expenses_url" {
  type    = string
  default = getenv("TURSO_EXPENSES_URL")
}

variable "expenses_token" {
  type    = string
  default = getenv("TURSO_EXPENSES_TOKEN")
}

variable "heat_url" {
  type    = string
  default = getenv("TURSO_HEAT_URL")
}

variable "heat_token" {
  type    = string
  default = getenv("TURSO_HEAT_TOKEN")
}

env "calories" {
  url = "${var.calories_url}?authToken=${var.calories_token}"
  migration {
    dir = "file://migrations/calories"
  }
  exclude = ["_litestream*"]
}

env "expenses" {
  url = "${var.expenses_url}?authToken=${var.expenses_token}"
  migration {
    dir = "file://migrations/expenses"
  }
  exclude = ["_litestream*"]
}

env "heat" {
  url = "${var.heat_url}?authToken=${var.heat_token}"
  migration {
    dir = "file://migrations/heat"
  }
  exclude = ["_litestream*"]
}
