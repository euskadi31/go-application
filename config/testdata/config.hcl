
var "hostname" {
    value = "localhost"
}

app {
    environment = "dev"
}

provider "mysql" {
    host        = var.hostname
    port        = 3306
    username    = "root"
    password    = ""
    dbname      = "demo"
}

provider "mysql" {
    key         = "user"
    host        = var.hostname
    port        = 3306
    username    = "root"
    password    = ""
    dbname      = "user"
}

