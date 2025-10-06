# Row Level Security Backend

Backend API dengan implementasi Row Level Security (RLS) menggunakan Go dan PostgreSQL.

## Fitur

- **Row Level Security (RLS)**: Implementasi keamanan tingkat baris di PostgreSQL
- **Multi-tenant Architecture**: Isolasi data berdasarkan tenant
- **JWT Authentication**: Sistem autentikasi menggunakan JSON Web Token
- **RESTful API**: API endpoints yang mengikuti standar REST
- **Database Migration**: Sistem migrasi database otomatis
- **Docker Support**: Containerization dengan Docker dan Docker Compose

## Teknologi

- **Backend**: Go (Golang)
- **Framework**: Gin
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT
- **Migration**: golang-migrate
- **Containerization**: Docker

## Struktur Project
├── cmd/server/          # Entry point aplikasi
├── config/              # Konfigurasi aplikasi
├── database/            # Database connection dan migration
├── handlers/            # HTTP handlers
├── middleware/          # Middleware (auth, RLS)
├── models/              # Data models
├── routes/              # Route definitions
├── utils/               # Utility functions
├── migrations/          # Database migrations
├── docker-compose.yml   # Docker compose configuration
├── Dockerfile          # Docker configuration
└── README.md           # Dokumentasi


## Instalasi dan Setup

### 1. Clone Repository
```bash
git clone <repository-url>
cd row-level-security-backend
```

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Setup Database
```bash
# Menggunakan Docker Compose
docker-compose up -d postgres

# Atau setup PostgreSQL manual dan update .env file
```

### 4. Environment Variables
Copy `.env.example` ke `.env` dan sesuaikan konfigurasi:
```bash
cp .env.example .env
```

### 5. Run Application
```bash
# Development
go run cmd/server/main.go

# Atau menggunakan Docker Compose
docker-compose up
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register user baru
- `POST /api/v1/auth/login` - Login user

### Posts (Protected)
- `GET /api/v1/posts` - Get all posts (dengan RLS)
- `POST /api/v1/posts` - Create new post
- `GET /api/v1/posts/:id` - Get specific post
- `PUT /api/v1/posts/:id` - Update post
- `DELETE /api/v1/posts/:id` - Delete post

## Row Level Security (RLS)

### Konsep RLS
RLS memungkinkan kontrol akses pada tingkat baris dalam database. Setiap query secara otomatis difilter berdasarkan policy yang didefinisikan.

### Implementasi
1. **Tenant Isolation**: Setiap user hanya bisa mengakses data dari tenant mereka
2. **User-specific Data**: User hanya bisa mengakses/memodifikasi data mereka sendiri
3. **Public/Private Content**: Sistem visibility untuk konten publik/privat

### Database Policies
```sql
-- Policy untuk posts: user hanya bisa melihat posts dari tenant mereka
CREATE POLICY post_tenant_isolation_policy ON posts
    FOR SELECT
    USING (
        tenant_id = current_setting('app.current_tenant_id')::INTEGER
        AND (is_public = true OR user_id = current_setting('app.current_user_id')::INTEGER)
    );
```

## Testing

### Manual Testing dengan cURL

1. **Register User**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "name": "Test User",
    "tenant_id": 1
  }'
```

2. **Login**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

3. **Create Post**
```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Test Post",
    "content": "This is a test post",
    "is_public": true
  }'
```

## Security Features

1. **JWT Authentication**: Semua endpoint protected menggunakan JWT
2. **Password Hashing**: Password di-hash menggunakan bcrypt
3. **Row Level Security**: Database-level security policies
4. **Tenant Isolation**: Complete data isolation between tenants
5. **Input Validation**: Request validation menggunakan Gin binding

## Development

### Adding New Features
1. Buat model di `models/`
2. Buat migration di `migrations/`
3. Buat handler di `handlers/`
4. Tambahkan route di `routes/`
5. Update RLS policies jika diperlukan

### Database Migration
```bash
# Create new migration
migrate create -ext sql -dir migrations -seq create_new_table

# Run migrations
go run cmd/server/main.go
```

## Production Deployment

1. Update environment variables untuk production
2. Use proper JWT secret
3. Configure database connection
4. Setup reverse proxy (nginx)
5. Enable HTTPS
6. Monitor logs dan performance

## Contributing

1. Fork repository
2. Create feature branch
3. Commit changes
4. Push to branch
5. Create Pull Request