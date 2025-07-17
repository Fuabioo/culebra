-- Example Lua configuration file
app = {
    name = "My Awesome App",
    version = "1.0.0"
}

database = {
    host = "localhost",
    port = 5432,
    name = "myapp_db",
    user = "postgres"
}

debug = true

-- You can use Lua logic in your config
if os.getenv("ENVIRONMENT") == "production" then
    debug = false
    database.host = "prod-db.example.com"
end

-- Arrays/lists are supported
allowed_hosts = {"localhost", "127.0.0.1", "example.com"}

-- Nested configuration
server = {
    http = {
        port = 8080,
        timeout = 30
    },
    tls = {
        enabled = false,
        cert_file = "/path/to/cert.pem",
        key_file = "/path/to/key.pem"
    }
}