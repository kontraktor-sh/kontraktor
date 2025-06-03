# Secret Management

Kontraktor provides secure secret management through integration with various secret vaults. This guide explains how to use secret management features.

## Supported Vaults

Currently, Kontraktor supports the following secret vaults:

- Azure Key Vault
- HashiCorp Vault (coming soon)
- AWS Secrets Manager (coming soon)

## Azure Key Vault Integration

### Prerequisites

1. An Azure subscription
2. An Azure Key Vault instance
3. Proper authentication configured (see [Azure Key Vault Integration](advanced/azure-keyvault.md) for details)

### Configuration

Configure Azure Key Vault in your taskfile:

```yaml
version: "0.3"

vaults:
  azure_keyvault:
    my-vault:                    # Vault configuration name
      keyvault_name: my-vault    # Azure Key Vault name
      secrets:                   # Map of environment variables to secret names
        API_KEY: api-secret
        DB_PASSWORD: db-secret
```

### Using Secrets

Secrets are automatically loaded as environment variables and can be used in your tasks:

```yaml
tasks:
  deploy:
    desc: Deploy with secrets
    cmds:
      - echo "Using API key: ${API_KEY}"
      - echo "Using DB password: ${DB_PASSWORD}"
```

## Security Best Practices

1. **Never Store Secrets in Taskfiles**
   - Keep all sensitive data in vaults
   - Don't commit secrets to version control
   - Use environment variables for non-sensitive configuration

2. **Access Control**
   - Use least privilege principle
   - Regularly rotate secrets
   - Monitor secret access

3. **Secret Naming**
   - Use descriptive names
   - Follow a consistent naming convention
   - Document secret purposes

4. **Error Handling**
   - Handle missing secrets gracefully
   - Log secret access errors
   - Implement fallback mechanisms

## Secret Access Patterns

### Direct Access

```yaml
tasks:
  direct-access:
    cmds:
      - echo "Secret value: ${SECRET_NAME}"
```

### Conditional Access

```yaml
tasks:
  conditional-access:
    cmds:
      - |
        if [ -n "${SECRET_NAME}" ]; then
          echo "Secret is available"
        else
          echo "Secret is not available"
        fi
```

### Secret Rotation

```yaml
tasks:
  rotate-secret:
    desc: Rotate a secret
    cmds:
      - echo "Current secret: ${OLD_SECRET}"
      - echo "New secret: ${NEW_SECRET}"
      - echo "Updating secret..."
```

## Troubleshooting

### Common Issues

1. **Authentication Failures**
   - Check Azure credentials
   - Verify vault permissions
   - Check network connectivity

2. **Missing Secrets**
   - Verify secret names
   - Check vault configuration
   - Ensure proper access rights

3. **Environment Variable Issues**
   - Check variable names
   - Verify taskfile syntax
   - Check vault configuration

### Debugging

Enable debug logging to troubleshoot secret access:

```bash
KONTRAKTOR_DEBUG=1 kontraktor run task-name
```

## Advanced Topics

- [Azure Key Vault Integration](advanced/azure-keyvault.md)
- [Secret Rotation Strategies](advanced/secret-rotation.md)
- [Access Control and Permissions](advanced/access-control.md)

## Examples

### Complete Example

```yaml
version: "0.3"

vaults:
  azure_keyvault:
    prod-vault:
      keyvault_name: prod-keyvault
      secrets:
        API_KEY: api-secret
        DB_PASSWORD: db-secret
        JWT_SECRET: jwt-secret

tasks:
  deploy:
    desc: Deploy to production
    environment:
      ENVIRONMENT: prod
    cmds:
      - echo "Deploying to ${ENVIRONMENT}"
      - echo "Using API key: ${API_KEY}"
      - echo "Using DB password: ${DB_PASSWORD}"
      - echo "Using JWT secret: ${JWT_SECRET}"

  rotate-secrets:
    desc: Rotate production secrets
    cmds:
      - echo "Rotating API key..."
      - echo "Rotating DB password..."
      - echo "Rotating JWT secret..."
``` 