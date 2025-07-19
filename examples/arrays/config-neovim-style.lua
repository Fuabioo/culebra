-- Modern Neovim-style configuration with return statement
-- This demonstrates the enhanced array conversion capabilities

return {
    app = {
        name = "ModernApp",
        version = "2.0.0",
        environment = os.getenv("APP_ENV") or "development"
    },
    
    database = {
        primary = "postgres://localhost:5432/myapp",
        replicas = {
            "postgres://replica1:5432/myapp",
            "postgres://replica2:5432/myapp",
            "postgres://replica3:5432/myapp"
        },
        connection_pool = {
            min_size = 5,
            max_size = 20,
            timeouts = {30, 60, 120}  -- connection, query, idle timeouts
        }
    },
    
    services = {
        {
            name = "api",
            port = 8080,
            endpoints = {"/health", "/metrics", "/api/v1"}
        },
        {
            name = "worker", 
            port = 8081,
            queues = {"high", "normal", "low"}
        },
        {
            name = "scheduler",
            port = 8082,
            cron_jobs = {
                "0 */6 * * *",    -- every 6 hours
                "0 0 * * 0",      -- weekly
                "0 0 1 * *"       -- monthly
            }
        }
    },
    
    logging = {
        level = "info",
        outputs = {"stdout", "file", "syslog"},
        formatters = {
            {name = "json", enabled = true},
            {name = "text", enabled = false}
        }
    },
    
    feature_flags = {
        new_ui = true,
        beta_api = false,
        experimental_cache = true
    },
    
    -- Mixed array types
    mixed_config = {
        items = {1, "two", true, 4.5},
        complex = {
            {id = 1, active = true, tags = {"urgent", "bug"}},
            {id = 2, active = false, tags = {"feature", "enhancement"}}
        }
    }
}