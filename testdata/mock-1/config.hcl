
var "foo" {
    value = "bar"
}

provider "mock-1" {
    foo      = var.foo
}
