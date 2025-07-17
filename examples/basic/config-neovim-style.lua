-- Neovim-style configuration with functions and modular setup
local config = {}

-- Helper functions (like Neovim's vim.fn)
local function env(key, default)
    return os.getenv(key) or default
end

local function setup_database()
    local db_config = {
        host = env("DB_HOST", "localhost"),
        port = tonumber(env("DB_PORT", "5432")),
        name = env("DB_NAME", "myapp_db"),
        user = env("DB_USER", "postgres"),
        ssl_mode = env("DB_SSL", "disable")
    }

    -- Environment-specific overrides
    if env("ENVIRONMENT") == "production" then
        db_config.host = "prod-db.example.com"
        db_config.ssl_mode = "require"
        db_config.pool_size = 20
    elseif env("ENVIRONMENT") == "development" then
        db_config.pool_size = 5
        db_config.log_queries = true
    end

    return db_config
end

local function setup_server()
    local server_config = {
        http = {
            port = tonumber(env("PORT", "8080")),
            timeout = 30,
            read_timeout = 10,
            write_timeout = 10
        },
        tls = {
            enabled = env("TLS_ENABLED", "false") == "true",
            cert_file = env("TLS_CERT", "/etc/ssl/cert.pem"),
            key_file = env("TLS_KEY", "/etc/ssl/key.pem")
        }
    }

    -- Auto-enable TLS in production
    if env("ENVIRONMENT") == "production" then
        server_config.tls.enabled = true
        server_config.http.port = 443
    end

    return server_config
end

local function setup_logging()
    local log_level = env("LOG_LEVEL", "info")
    local is_prod = env("ENVIRONMENT") == "production"

    return {
        level = log_level,
        format = is_prod and "json" or "text",
        output = is_prod and "/var/log/app.log" or "stdout",
        rotate = is_prod,
        max_size = "100MB",
        max_files = 10
    }
end

-- Plugin-like configuration modules
local plugins = {
    metrics = {
        enabled = env("METRICS_ENABLED", "true") == "true",
        endpoint = "/metrics",
        port = tonumber(env("METRICS_PORT", "9090"))
    },

    auth = {
        provider = env("AUTH_PROVIDER", "jwt"),
        secret = env("JWT_SECRET", "dev-secret-change-me"),
        expiry = "24h"
    },

    cache = {
        enabled = env("CACHE_ENABLED", "true") == "true",
        backend = env("CACHE_BACKEND", "redis"),
        ttl = tonumber(env("CACHE_TTL", "3600"))
    }
}

-- Main configuration assembly (like Neovim's init.lua)
config.app = {
    name = "My Awesome Neovim-like App",
    version = "1.0.0",
    environment = env("ENVIRONMENT", "development")
}

config.database = setup_database()
config.server = setup_server()
config.logging = setup_logging()
config.plugins = plugins

-- Feature flags (very Neovim-like)
config.features = {
    debug_mode = env("ENVIRONMENT") ~= "production",
    hot_reload = env("ENVIRONMENT") == "development",
    profiling = env("ENABLE_PROFILING", "false") == "true",
    experimental_features = env("EXPERIMENTAL", "false") == "true"
}

-- Security settings
config.security = {
    cors = {
        enabled = true,
        origins = env("ENVIRONMENT") == "production"
            and { "https://myapp.com", "https://api.myapp.com" }
            or { "*" }
    },
    rate_limit = {
        enabled = env("ENVIRONMENT") == "production",
        requests_per_minute = tonumber(env("RATE_LIMIT", "100"))
    }
}

-- Validation function (like Neovim's vim.validate)
local function validate_config()
    assert(config.database.host, "Database host is required")
    assert(config.database.port > 0, "Database port must be positive")
    assert(config.server.http.port > 0, "Server port must be positive")

    if config.server.tls.enabled then
        assert(config.server.tls.cert_file, "TLS cert file required when TLS enabled")
        assert(config.server.tls.key_file, "TLS key file required when TLS enabled")
    end
end

-- Run validation
validate_config()

-- Export configuration (return statement like Neovim plugins)
return config
