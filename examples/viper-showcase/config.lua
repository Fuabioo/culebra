-- Comprehensive configuration file showcasing all data types for Viper Get functions

-- Application information
app = {
    name = "Viper Showcase App",
    description = "A comprehensive demonstration of Viper Get functions",
    version = "2.1.0"
}

-- String values including edge cases
empty_string = ""
string_number = "42"
string_bool = "true"

-- Server configuration with various integer types
server = {
    port = 8080,
    max_connections = 1000,
    workers = 4,
    timeout_seconds = 30
}

-- Boolean values
debug = {
    enabled = true
}
ssl = {
    enabled = false
}
production_mode = false

-- Numeric values of different types
numbers = {
    zero = 0,
    negative = -123,
    large_int64 = 9223372036854775807,  -- Max int64
    pi = 3.14159265359,
    small_float = 0.001
}

-- Floating point values
timeouts = {
    request = 30.5,
    connection = 5.0,
    idle = 300.0
}

metrics = {
    success_rate = 99.9
}

-- Array values (will be converted to slices with ConvertArrays: true)
database = {
    hosts = {"primary.db.com", "secondary.db.com", "tertiary.db.com"},
    primary = {
        host = "primary.db.com",
        port = 5432,
        ssl = true,
        timeout = "30s"
    },
    replicas = {
        count = 3,
        hosts = {"replica1.db.com", "replica2.db.com", "replica3.db.com"}
    }
}

allowed_ports = {80, 443, 8080, 8443, 9000}
tags = {"web", "api", "database", "cache"}
empty_array = {}

-- Time and duration values (as strings for Viper parsing)
timestamps = {
    startup = "2023-01-01T12:00:00Z",
    last_backup = "2023-12-31T23:59:59Z"
}

cache = {
    ttl = "5m30s",      -- 5 minutes 30 seconds
    cleanup = "1h"       -- 1 hour
}

monitoring = {
    heartbeat = "15s",   -- 15 seconds
    metrics_interval = "1m",  -- 1 minute
    health_check = "10s"      -- 10 seconds
}

-- Map values for GetStringMap and GetStringMapString
environment = {
    LOG_LEVEL = "info",
    MAX_WORKERS = "10",
    CACHE_SIZE = "1000",
    DEBUG_MODE = "false"
}

-- Logging configuration
logging = {
    level = "info",
    file = "/var/log/app.log",
    rotate = true,
    max_size = "100MB"
}

-- Feature flags and settings
features = {
    authentication = true,
    rate_limiting = true,
    caching = false,
    metrics = true
}

-- Complex nested structures
services = {
    web = {
        instances = 3,
        ports = {8080, 8081, 8082},
        health_check = "/health",
        timeout = "30s"
    },
    api = {
        instances = 2,
        ports = {9000, 9001},
        health_check = "/api/health",
        timeout = "60s"
    },
    cache = {
        type = "redis",
        host = "cache.internal",
        port = 6379,
        timeout = "5s"
    }
}

-- Security configuration
security = {
    jwt = {
        secret = "super-secret-key",
        expiry = "24h",
        refresh_expiry = "168h"  -- 1 week
    },
    cors = {
        enabled = true,
        origins = {"https://example.com", "https://api.example.com"},
        methods = {"GET", "POST", "PUT", "DELETE"},
        headers = {"Content-Type", "Authorization", "X-API-Key"}
    },
    rate_limiting = {
        enabled = true,
        requests_per_minute = 1000,
        burst_size = 100
    }
}

-- Edge case values
nil_value = nil

-- Additional numeric edge cases
limits = {
    max_uint32 = 4294967295,      -- Max uint32
    max_int32 = 2147483647,       -- Max int32
    min_int32 = -2147483648,      -- Min int32
    tiny_float = 0.000001,
    big_float = 123456789.987654321
}

-- Mixed type arrays (testing type flexibility)
mixed_data = {
    values = {"string", 42, true, 3.14, nil},
    config_levels = {"debug", "info", "warn", "error"},
    retry_delays = {1.0, 2.5, 5.0, 10.0, 30.0}
}