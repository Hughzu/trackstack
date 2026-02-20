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

variable "users_url" {
  type    = string
  default = getenv("TURSO_USERS_URL")
}

variable "users_token" {
  type    = string
  default = getenv("TURSO_USERS_TOKEN")
}

env "calories" {
  url = "${var.calories_url}?authToken=${var.calories_token}"
  migration {
    dir = "file://calories"
  }
  exclude = ["_litestream*"]
}

env "expenses" {
  url = "${var.expenses_url}?authToken=${var.expenses_token}"
  migration {
    dir = "file://expenses"
  }
  exclude = ["_litestream*"]
}

env "heat" {
  url = "${var.heat_url}?authToken=${var.heat_token}"
  migration {
    dir = "file://heat"
  }
  exclude = ["_litestream*"]
}

env "users" {
  url = "${var.users_url}?authToken=${var.users_token}"
  migration {
    dir = "file://users"
  }
  exclude = ["_litestream*"]
}
