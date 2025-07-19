-- Traditional Culebra configuration with global variables
-- This demonstrates backward compatibility

app_name = "MyApp"
version = "1.0.0"

-- Arrays in traditional style
database_hosts = {"localhost", "db1.example.com", "db2.example.com"}
ports = {8080, 8081, 8082}

-- Nested structure with arrays
features = {
    authentication = {"oauth", "jwt", "basic"},
    storage = {"redis", "postgres", "s3"},
    monitoring = {"prometheus", "grafana"}
}

-- Configuration flags
debug_enabled = true
max_connections = 100