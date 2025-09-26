# Tesla HVAC Configuration Guide

This guide explains how to configure the Tesla HVAC client using the configuration management system.

## Configuration File Location

The default configuration file is located at:
- Linux/macOS: `~/.config/tesla-hvac/config.json`
- Windows: `%APPDATA%\tesla-hvac\config.json`

You can specify a custom path using the `-config` flag or environment variable.

## Quick Start

1. **Create a configuration file:**
   ```bash
   tesla-config -action create
   ```

2. **Set your vehicle VIN:**
   ```bash
   tesla-config -action set-vin -vin 5YJ3E1EA4KF123456
   ```

3. **Set your private key file:**
   ```bash
   tesla-config -action set-key -key-file ~/.tesla/private_key.pem
   ```

4. **Set your OAuth token file:**
   ```bash
   tesla-config -action set-token -token-file ~/.tesla/oauth_token.json
   ```

5. **Validate your configuration:**
   ```bash
   tesla-config -action validate
   ```

## Configuration Structure

### Tesla Configuration (`tesla`)

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `vin` | string | Vehicle Identification Number | Required |
| `private_key_file` | string | Path to private key file | "" |
| `oauth_token_file` | string | Path to OAuth token file | "" |
| `connection_timeout` | duration | Connection timeout | 60s |
| `scan_timeout` | duration | Vehicle scan timeout | 30s |
| `max_concurrent_requests` | int | Max concurrent API requests | 5 |
| `request_timeout` | duration | Individual request timeout | 10s |
| `scan_retries` | int | Number of scan retry attempts | 3 |
| `scan_delay` | duration | Delay between scan attempts | 2s |

### Client Configuration (`client`)

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `client_name` | string | Client application name | "tesla-hvac-client" |
| `client_version` | string | Client version | "1.0.0" |
| `keep_alive_interval` | duration | Keep-alive interval | 30s |
| `health_check_interval` | duration | Health check interval | 60s |
| `enable_auto_reconnect` | bool | Enable automatic reconnection | true |
| `enable_health_checks` | bool | Enable health monitoring | true |
| `enable_metrics` | bool | Enable metrics collection | false |

### Retry Configuration (`retry`)

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `max_retries` | int | Maximum retry attempts | 3 |
| `initial_delay` | duration | Initial retry delay | 1s |
| `max_delay` | duration | Maximum retry delay | 30s |
| `backoff_factor` | float | Exponential backoff factor | 2.0 |
| `jitter` | bool | Add jitter to retry delays | true |

### Circuit Breaker Configuration (`circuit_breaker`)

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `max_failures` | int | Failures before opening circuit | 5 |
| `reset_timeout` | duration | Time before attempting reset | 60s |
| `half_open_max_calls` | int | Max calls in half-open state | 3 |

### Logging Configuration (`logging`)

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `level` | string | Log level (debug, info, warn, error) | "info" |
| `format` | string | Log format (json, text) | "text" |
| `output` | string | Output destination (stdout, stderr, file) | "stdout" |
| `file_path` | string | Log file path (when output is file) | "" |
| `max_size` | int | Max log file size in MB | 100 |
| `max_backups` | int | Max number of backup files | 3 |
| `max_age` | int | Max age in days | 7 |

## Environment Variables

You can override configuration values using environment variables:

| Environment Variable | Configuration Field |
|---------------------|-------------------|
| `TESLA_VIN` | `tesla.vin` |
| `TESLA_PRIVATE_KEY_FILE` | `tesla.private_key_file` |
| `TESLA_OAUTH_TOKEN_FILE` | `tesla.oauth_token_file` |
| `TESLA_CLIENT_NAME` | `client.client_name` |
| `TESLA_CLIENT_VERSION` | `client.client_version` |
| `TESLA_LOG_LEVEL` | `logging.level` |
| `TESLA_LOG_FORMAT` | `logging.format` |
| `TESLA_LOG_OUTPUT` | `logging.output` |
| `TESLA_LOG_FILE` | `logging.file_path` |

## Configuration Management

### Hot Reloading

The configuration system supports hot-reloading. When you modify the configuration file, the client will automatically reload the configuration and apply changes without restarting.

### Configuration Validation

The configuration system validates all settings on load and save. Invalid configurations will be rejected with detailed error messages.

### Configuration Builder

You can programmatically build configurations using the `ConfigBuilder`:

```go
config := tesla.NewConfigBuilder().
    WithVIN("5YJ3E1EA4KF123456").
    WithPrivateKeyFile("/path/to/private_key.pem").
    WithRetryConfig(tesla.RetryConfig{
        MaxRetries: 5,
        InitialDelay: 2 * time.Second,
        MaxDelay: 60 * time.Second,
        BackoffFactor: 2.0,
        Jitter: true,
    }).
    Build()
```

## Example Configuration

```json
{
  "tesla": {
    "vin": "5YJ3E1EA4KF123456",
    "private_key_file": "/home/user/.tesla/private_key.pem",
    "oauth_token_file": "/home/user/.tesla/oauth_token.json",
    "connection_timeout": "60s",
    "scan_timeout": "30s",
    "max_concurrent_requests": 5,
    "request_timeout": "10s",
    "scan_retries": 3,
    "scan_delay": "2s"
  },
  "client": {
    "client_name": "tesla-hvac-client",
    "client_version": "1.0.0",
    "keep_alive_interval": "30s",
    "health_check_interval": "60s",
    "enable_auto_reconnect": true,
    "enable_health_checks": true,
    "enable_metrics": false
  },
  "retry": {
    "max_retries": 3,
    "initial_delay": "1s",
    "max_delay": "30s",
    "backoff_factor": 2.0,
    "jitter": true
  },
  "circuit_breaker": {
    "max_failures": 5,
    "reset_timeout": "60s",
    "half_open_max_calls": 3
  },
  "logging": {
    "level": "info",
    "format": "text",
    "output": "stdout",
    "file_path": "",
    "max_size": 100,
    "max_backups": 3,
    "max_age": 7
  }
}
```

## Troubleshooting

### Common Issues

1. **Configuration file not found**: Use `tesla-config -action create` to create a default configuration.

2. **Invalid VIN**: Ensure your VIN is 17 characters long and contains only valid characters.

3. **Private key file not found**: Verify the path to your private key file is correct.

4. **OAuth token file not found**: Verify the path to your OAuth token file is correct.

5. **Permission denied**: Ensure the configuration directory is writable.

### Debug Mode

Enable debug logging to troubleshoot issues:

```bash
tesla-config -action set-log-level -level debug
```

### Configuration Validation

Always validate your configuration after making changes:

```bash
tesla-config -action validate
```

## Security Considerations

- Store private keys and OAuth tokens in secure locations
- Use appropriate file permissions (600) for sensitive files
- Consider using environment variables for sensitive configuration in production
- Regularly rotate OAuth tokens
- Keep private keys secure and never share them
