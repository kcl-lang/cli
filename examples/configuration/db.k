env: str = option("env") or "dev"  # The attribute `env` has a default value "den"
database: str = option("database")
hosts = {
    dev = "postgres.dev"
    stage = "postgres.stage"
    prod = "postgres.prod"
}
dbConfig = {
    host = hosts[env]
    database = database
    port = "2023"
    conn = "postgres://${host}:${port}/${database}"
}
