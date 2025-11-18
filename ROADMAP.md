# Culebra Roadmap

## 🚀 Future Improvements

### HTTP/Network Support
Currently, Culebra supports local configuration logic but doesn't include HTTP capabilities. Future versions could add network support for advanced use cases:

```lua
-- Future capability: Remote configuration fetching
local feature_flags = http.get("https://api.company.com/feature-flags")
local service_discovery = http.get("https://consul.company.com/v1/catalog/services")

config.features = feature_flags.enabled
config.services = service_discovery.endpoints
```

This would enable:
- **Remote Feature Flags**: Dynamic feature toggling from external services
- **Service Discovery**: Automatic service endpoint configuration
- **External Configuration**: Loading partial config from remote sources
- **API Integration**: Direct integration with configuration management APIs

**Implementation Options**:
1. Custom HTTP module for gopher-lua
2. Integration with existing `gluahttp` library
3. Go-side pre-fetching with Lua globals
4. Async configuration updates

### Other Potential Enhancements
- **Configuration Caching**: Cache remote configurations with TTL
- **Hot Reloading**: Watch for configuration changes and reload automatically
- **Configuration Validation**: Schema validation for Lua configurations
- **Encrypted Configurations**: Support for encrypted configuration files
- **Configuration Templating**: Advanced templating beyond basic Lua logic
- **Execution Timeouts**: Add context-based timeouts for Lua script execution
- **Resource Monitoring**: Track memory usage and execution time metrics
