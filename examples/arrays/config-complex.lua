-- Complex configuration demonstrating deep nesting and various array types
-- This shows the full power of the array conversion feature

return {
    -- Microservices configuration
    microservices = {
        {
            name = "user-service",
            instances = 3,
            ports = {8001, 8002, 8003},
            routes = {"/users", "/auth", "/profile"},
            dependencies = {"database", "redis", "email-service"}
        },
        {
            name = "order-service", 
            instances = 2,
            ports = {8010, 8011},
            routes = {"/orders", "/cart", "/checkout"},
            dependencies = {"database", "payment-service", "inventory-service"}
        },
        {
            name = "notification-service",
            instances = 1,
            ports = {8020},
            routes = {"/notifications", "/webhooks"},
            dependencies = {"redis", "email-service", "sms-service"}
        }
    },
    
    -- Infrastructure configuration
    infrastructure = {
        load_balancers = {
            {
                name = "main-lb",
                algorithm = "round_robin",
                targets = {"10.0.1.10", "10.0.1.11", "10.0.1.12"},
                health_checks = {"/health", "/ready"},
                timeouts = {
                    connect = 5,
                    response = 30,
                    idle = 300
                }
            }
        },
        
        databases = {
            {
                type = "postgresql",
                primary = "db-primary.internal",
                replicas = {"db-replica-1.internal", "db-replica-2.internal"},
                schemas = {"users", "orders", "analytics"},
                connection_limits = {25, 15, 10}  -- per schema
            },
            {
                type = "redis", 
                cluster_nodes = {
                    "redis-1.internal:6379",
                    "redis-2.internal:6379", 
                    "redis-3.internal:6379"
                },
                key_patterns = {"session:*", "cache:*", "queue:*"}
            }
        }
    },
    
    -- Monitoring and alerting
    monitoring = {
        metrics = {
            collectors = {"prometheus", "statsd", "newrelic"},
            intervals = {15, 60, 300},  -- seconds
            retention_days = {7, 30, 365}
        },
        
        alerts = {
            {
                name = "high_cpu",
                thresholds = {80, 90, 95},  -- warning, critical, emergency
                services = {"user-service", "order-service"},
                channels = {"slack", "email", "pagerduty"}
            },
            {
                name = "database_connections",
                thresholds = {70, 85, 95},
                targets = {"primary", "replica-1", "replica-2"},
                escalation_minutes = {5, 15, 30}
            }
        }
    },
    
    -- Security configuration
    security = {
        allowed_origins = {
            "https://app.example.com",
            "https://admin.example.com", 
            "https://mobile.example.com"
        },
        
        api_keys = {
            {
                name = "admin",
                permissions = {"read", "write", "delete", "admin"},
                rate_limits = {1000, 100, 10}  -- per minute for read/write/admin
            },
            {
                name = "readonly",
                permissions = {"read"},
                rate_limits = {5000, 0, 0}
            }
        },
        
        encryption = {
            algorithms = {"AES-256-GCM", "ChaCha20-Poly1305"},
            key_rotation_days = {30, 90, 365},
            backup_locations = {"vault", "hsm", "kms"}
        }
    }
}