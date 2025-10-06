# Security Guidelines

## Environment Variables Setup

### Required Environment Variables
Aplikasi ini memerlukan environment variables berikut:

- `DATABASE_URL`: Connection string ke PostgreSQL
- `JWT_SECRET`: Secret key untuk JWT (minimum 32 karakter)
- `API_KEY`: API key untuk external services

### Setup untuk Development

1. **Copy template environment file:**
```bash
cp .env.example .env
```

2. **Generate secure secrets:**
```bash
chmod +x scripts/generate-secrets.sh
./scripts/generate-secrets.sh
```

3. **Update .env file dengan values yang di-generate**

### Setup untuk Production

1. **Jangan pernah commit file .env ke repository**
2. **Set environment variables di server/container:**
```bash
export DATABASE_URL="postgres://user:pass@host:5432/db"
export JWT_SECRET="your-32-char-secret-key"
export API_KEY="your-api-key"
```

3. **Atau gunakan secret management service** (AWS Secrets Manager, HashiCorp Vault, dll)

## Security Best Practices

### JWT Secret
- Minimum 32 karakter
- Generate menggunakan cryptographically secure random generator
- Rotate secara berkala
- Jangan hardcode di aplikasi

### Database
- Gunakan strong password
- Enable SSL/TLS untuk production
- Limit database user permissions
- Regular backup dan test restore

### API Key
- Generate random hex/base64 string
- Store securely
- Implement rate limiting
- Monitor usage

## Environment Variables Validation

Aplikasi akan melakukan validasi saat startup:
- Memastikan semua required env vars ada
- Validasi format dan panjang minimum
- Gagal startup jika ada yang missing

## Deployment Checklist

- [ ] All secrets generated securely
- [ ] No hardcoded secrets in code
- [ ] Environment variables set in deployment
- [ ] .env file not in version control
- [ ] SSL/TLS enabled for database
- [ ] Rate limiting implemented
- [ ] Monitoring and logging configured