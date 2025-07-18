-- Auto-loaded Lua configuration
-- This file will be automatically discovered by Viper when named "example.lua"
-- and placed in one of the configured search paths

app = {
    name = "Auto-loaded App",
    version = "2.0.0"
}

database = {
    host = "auto-db.example.com",
    port = 5432,
    name = "autoload_db",
    user = "autouser"
}

-- Demonstrate conditional configuration based on environment
local env = os.getenv("APP_ENV") or "development"
environment = env

if env == "production" then
    debug = false
    database.host = "prod-db.company.com"
    database.port = 5432
elseif env == "development" then
    debug = true
    database.host = "localhost"
    database.port = 5433
else
    debug = true  -- Default for other environments
end

-- Server configuration
server = {
    port = 8080,
    host = "0.0.0.0",
    workers = 4
}

-- Feature flags
features = {
    new_ui = true,
    beta_features = false,
    analytics = env ~= "development"
}

-- Logging configuration
logging = {
    level = debug and "debug" or "info",
    file = "/var/log/autoload-app.log",
    console = debug
}