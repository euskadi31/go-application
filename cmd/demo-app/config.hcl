
var "hostname" {
    value = "localhost"
}

app {

}

provider "mysql" {
    host        = var.hostname
    port        = 3306
    username    = "root"
    password    = ""
    dbname      = "demo"
}

provider "mysql" {
    key         = "backup"
    host        = var.hostname
    port        = 3306
    username    = "root"
    password    = ""
    dbname      = "user"
}

