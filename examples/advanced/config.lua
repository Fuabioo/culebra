-- Advanced Lua Configuration Example
-- This demonstrates the power of Lua over static formats like JSON/YAML

-- Environment detection (from environment variables or defaults)
local environment = os.getenv("APP_ENV") or "development"
local is_production = environment == "production"
local is_development = environment == "development"

-- Base configuration
local config = {
    environment = environment,
    debug = not is_production,
    
    -- Application metadata
    app = {
        name = "advanced-example",
        version = "1.2.3"
    }
}

-- Conditional database configuration based on environment
if is_production then
    config.database = {
        url = os.getenv("DATABASE_URL") or "postgres://prod-server:5432/myapp",
        pool_size = 20,
        ssl_mode = "require",
        timeout = 30
    }
elseif environment == "staging" then
    config.database = {
        url = "postgres://staging-server:5432/myapp_staging",
        pool_size = 10,
        ssl_mode = "prefer",
        timeout = 15
    }
else -- development
    config.database = {
        url = "postgres://localhost:5432/myapp_dev",
        pool_size = 5,
        ssl_mode = "disable",
        timeout = 5
    }
end

-- Server configuration with environment-specific ports
local base_port = 8080
if environment == "test" then
    base_port = 8081
elseif is_production then
    base_port = 80
end

config.server = {
    host = is_production and "0.0.0.0" or "localhost",
    port = base_port,
    workers = is_production and 8 or 2,
    
    -- TLS configuration (only in production)
    tls = is_production and {
        enabled = true,
        cert_file = "/etc/ssl/certs/app.crt",
        key_file = "/etc/ssl/private/app.key"
    } or {
        enabled = false
    }
}

-- Feature flags with conditional logic
local features = {"auth", "logging"}

-- Add environment-specific features
if is_production then
    table.insert(features, "metrics")
    table.insert(features, "tracing")
    table.insert(features, "rate_limiting")
elseif is_development then
    table.insert(features, "debug_toolbar")
    table.insert(features, "hot_reload")
    table.insert(features, "profiling")
end

-- Add experimental features based on environment variable
if os.getenv("ENABLE_EXPERIMENTAL") == "true" then
    table.insert(features, "new_ui")
    table.insert(features, "ai_features")
end

config.features = features

-- Service discovery with computed endpoints
local services = {}

-- Helper function to generate service URLs
local function service_url(service, port)
    local host_suffix = is_production and ".prod.internal" or ".dev.local"
    return string.format("http://%s%s:%d", service, host_suffix, port)
end

services.user_service = service_url("user-service", 8001)
services.auth_service = service_url("auth-service", 8002)
services.notification_service = service_url("notification-service", 8003)

-- Add additional services in production
if is_production then
    services.analytics_service = service_url("analytics-service", 8004)
    services.billing_service = service_url("billing-service", 8005)
end

config.services = services

-- Rate limiting based on environment
config.rate_limit = is_production and 1000 or 10000  -- requests per minute

-- Logging configuration with computed log levels
config.logging = {
    level = is_production and "info" or "debug",
    format = is_production and "json" or "console",
    
    -- Conditional file logging
    file = is_production and {
        enabled = true,
        path = "/var/log/app/app.log",
        rotate = true,
        max_size = "100MB"
    } or nil
}

-- Monitoring configuration with computed endpoints
local monitoring_endpoints = {}

-- Helper function to create monitoring URLs
local function monitoring_url(path)
    local base = is_production and "https://monitoring.company.com" or "http://localhost:3000"
    return base .. path
end

table.insert(monitoring_endpoints, monitoring_url("/health"))
table.insert(monitoring_endpoints, monitoring_url("/metrics"))

if is_production then
    table.insert(monitoring_endpoints, monitoring_url("/alerts"))
    table.insert(monitoring_endpoints, monitoring_url("/dashboard"))
end

config.monitoring = {
    enabled = true,
    endpoints = monitoring_endpoints,
    interval = is_production and 30 or 60  -- seconds
}

-- Security configuration with environment-specific settings
config.security = {
    cors = {
        enabled = true,
        origins = is_production and {"https://myapp.com", "https://api.myapp.com"} 
                                 or {"*"}  -- Allow all in development
    },
    
    -- API key validation
    api_keys = {
        enabled = is_production,
        required_for = is_production and {"admin", "analytics"} or {}
    },
    
    -- Session configuration
    session = {
        secret = os.getenv("SESSION_SECRET") or "dev-secret-key",
        duration = is_production and 24 * 60 * 60 or 7 * 24 * 60 * 60,  -- seconds
        secure = is_production
    }
}

-- Cache configuration with computed settings
config.cache = {
    type = is_production and "redis" or "memory",
    
    -- Redis configuration (only if using Redis)
    redis = is_production and {
        host = os.getenv("REDIS_HOST") or "redis.prod.internal",
        port = 6379,
        db = 0,
        password = os.getenv("REDIS_PASSWORD")
    } or nil,
    
    -- TTL based on environment
    default_ttl = is_production and 300 or 60  -- seconds
}

-- Worker configuration with computed values
config.workers = {
    enabled = true,
    count = is_production and 4 or 1,
    
    -- Queue configuration
    queue = {
        type = is_production and "redis" or "memory",
        max_jobs = is_production and 1000 or 100
    }
}

-- Experimental: Configuration validation function
local function validate_config()
    -- Ensure required environment variables in production
    if is_production then
        local required_vars = {"DATABASE_URL", "SESSION_SECRET"}
        for _, var in ipairs(required_vars) do
            if not os.getenv(var) then
                error("Missing required environment variable: " .. var)
            end
        end
    end
    
    -- Validate port ranges
    if config.server.port < 1 or config.server.port > 65535 then
        error("Invalid port number: " .. config.server.port)
    end
    
    return true
end

-- Run validation
validate_config()

-- Dynamic configuration based on current time (example)
local current_hour = tonumber(os.date("%H"))
config.maintenance_mode = is_production and current_hour >= 2 and current_hour <= 4

-- Return the final configuration
return config