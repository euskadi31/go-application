
var "hostname" {
    value = "localhost"
}

app {
    environment = env.ENV
    env_prefix = "DEMO_APP"
}

logger {
    level = "debug"

    provider "console" {
        where {
            env = ["local"]
        }

    }

    provider "ecs" {
        where {
            env = ["prod", "dev"]
        }
    }
}

provider "logger" {
    level = "debug"

    writer "console" {
        where {
            env = ["local"]
        }

    }

    writer "ecs" {
        where {
            env = ["prod", "dev"]
        }
    }
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

