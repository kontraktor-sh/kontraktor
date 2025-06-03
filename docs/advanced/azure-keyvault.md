# Azure Key Vault Integration

This guide explains how to integrate Kontraktor with Azure Key Vault for secure secret management.

## Prerequisites

1. **Azure Subscription**
   - Active Azure subscription
   - Proper permissions to create and manage Key Vaults

2. **Azure CLI**
   - Install the [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli)
   - Log in to your Azure account:
     ```bash
     az login
     ```

3. **Azure Key Vault**
   - Create a Key Vault:
     ```bash
     az keyvault create --name my-keyvault --resource-group my-resource-group --location eastus
     ```

## Authentication

Kontraktor uses Azure's Default Credential Chain for authentication. This means it will try the following methods in order:

1. Environment variables
2. Managed Identity
3. Azure CLI credentials
4. Visual Studio Code credentials
5. Azure PowerShell credentials

### Environment Variables

Set the following environment variables:

```bash
export AZURE_TENANT_ID="your-tenant-id"
export AZURE_CLIENT_ID="your-client-id"
export AZURE_CLIENT_SECRET="your-client-secret"
```

### Managed Identity

If running in Azure (e.g., Azure VM, App Service), you can use Managed Identity:

1. Enable Managed Identity on your resource
2. Grant the identity access to Key Vault:
   ```bash
   az keyvault set-policy --name my-keyvault --object-id <identity-object-id> --secret-permissions get list
   ```

## Configuration

### Taskfile Configuration

Configure Azure Key Vault in your taskfile:

```yaml
version: "0.3"

vaults:
  azure_keyvault:
    my-vault:
      keyvault_name: my-keyvault
      secrets:
        API_KEY: api-secret
        DB_PASSWORD: db-secret
```

### Secret Management

1. **Adding Secrets**
   ```bash
   az keyvault secret set --vault-name my-keyvault --name api-secret --value "secret-value"
   ```

2. **Updating Secrets**
   ```bash
   az keyvault secret set --vault-name my-keyvault --name api-secret --value "new-value"
   ```

3. **Deleting Secrets**
   ```bash
   az keyvault secret delete --vault-name my-keyvault --name api-secret
   ```

## Using Secrets in Tasks

### Basic Usage

```yaml
tasks:
  deploy:
    desc: Deploy with secrets
    cmds:
      - echo "Using API key: ${API_KEY}"
      - echo "Using DB password: ${DB_PASSWORD}"
```

### Secret Rotation

```yaml
tasks:
  rotate-secrets:
    desc: Rotate secrets
    cmds:
      - |
        # Generate new secret
        NEW_SECRET=$(openssl rand -base64 32)
        
        # Update in Azure Key Vault
        az keyvault secret set --vault-name my-keyvault --name api-secret --value "$NEW_SECRET"
        
        # Verify update
        echo "Secret rotated successfully"
```

## Security Best Practices

1. **Access Control**
   - Use least privilege principle
   - Regularly audit access
   - Use Managed Identities when possible

2. **Secret Management**
   - Rotate secrets regularly
   - Use strong secret values
   - Monitor secret access

3. **Network Security**
   - Use private endpoints
   - Enable firewall rules
   - Use VNET integration

4. **Monitoring**
   - Enable diagnostic logging
   - Set up alerts
   - Monitor access patterns

## Troubleshooting

### Common Issues

1. **Authentication Failures**
   ```bash
   # Check Azure CLI login
   az account show
   
   # Verify environment variables
   env | grep AZURE
   ```

2. **Permission Issues**
   ```bash
   # Check Key Vault access
   az keyvault show --name my-keyvault
   
   # List access policies
   az keyvault show --name my-keyvault --query properties.accessPolicies
   ```

3. **Network Issues**
   ```bash
   # Test connectivity
   curl -v https://my-keyvault.vault.azure.net
   ```

### Debugging

Enable debug logging:

```bash
KONTRAKTOR_DEBUG=1 kontraktor run task-name
```

## Advanced Topics

### Multiple Key Vaults

```yaml
version: "0.3"

vaults:
  azure_keyvault:
    prod-vault:
      keyvault_name: prod-keyvault
      secrets:
        PROD_API_KEY: api-secret
    dev-vault:
      keyvault_name: dev-keyvault
      secrets:
        DEV_API_KEY: api-secret
```

### Secret Versioning

```yaml
tasks:
  deploy-specific-version:
    desc: Deploy with specific secret version
    cmds:
      - echo "Using specific version of secret"
      - az keyvault secret show --vault-name my-keyvault --name api-secret --version "specific-version"
```

### Backup and Restore

```yaml
tasks:
  backup-secrets:
    desc: Backup Key Vault secrets
    cmds:
      - |
        # Backup secrets
        for secret in $(az keyvault secret list --vault-name my-keyvault --query "[].id" -o tsv); do
          az keyvault secret backup --vault-name my-keyvault --id "$secret" --file "backup-$(basename $secret).bak"
        done
```

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
      - |
        # Rotate API key
        NEW_API_KEY=$(openssl rand -base64 32)
        az keyvault secret set --vault-name prod-keyvault --name api-secret --value "$NEW_API_KEY"
        
        # Rotate DB password
        NEW_DB_PASSWORD=$(openssl rand -base64 32)
        az keyvault secret set --vault-name prod-keyvault --name db-secret --value "$NEW_DB_PASSWORD"
        
        # Rotate JWT secret
        NEW_JWT_SECRET=$(openssl rand -base64 32)
        az keyvault secret set --vault-name prod-keyvault --name jwt-secret --value "$NEW_JWT_SECRET"
        
        echo "All secrets rotated successfully"
``` 