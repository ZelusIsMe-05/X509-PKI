# 🔒 Security Documentation

## Overview
This document outlines the security improvements made to the X.509 PKI system and provides guidance for secure deployment.

---

## 🔴 Critical Security Implementations

### 1. **JWT Secret from Environment Variables**
**Status:** ✅ FIXED

**Before:**
```go
const jwtSecret = "x509-pki-super-secret-key-2024"  // Hardcoded in code
```

**After:**
```go
var jwtSecret string

func init() {
    jwtSecret = os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        log.Fatal("JWT_SECRET environment variable not set")
    }
}
```

**Setup:**
1. Create `.env` file (copy from `.env.example`)
2. Set `JWT_SECRET` to a strong random string (min 32 characters)
3. **NEVER commit `.env` to Git** (already in `.gitignore`)

**Example:**
```bash
# Generate a secure random secret
openssl rand -base64 32
```

---

### 2. **Environment-Based CORS Configuration**
**Status:** ✅ FIXED

**Before:**
```go
w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
```

**After:**
```go
corsOrigin := os.Getenv("CORS_ORIGIN")
if corsOrigin == "" {
    corsOrigin = "http://localhost:5173"
}
w.Header().Set("Access-Control-Allow-Origin", corsOrigin)
```

**Setup in `.env`:**
```
# Development
CORS_ORIGIN=http://localhost:5173

# Production
CORS_ORIGIN=https://yourdomain.com
```

---

### 3. **Strong Input Validation**
**Status:** ✅ FIXED

#### Username Requirements:
- ✅ 3-50 characters
- ✅ Only alphanumeric, underscore, hyphen
- ✅ Regex validation applied

#### Password Requirements:
- ✅ Minimum 8 characters
- ✅ At least 1 uppercase letter (A-Z)
- ✅ At least 1 lowercase letter (a-z)
- ✅ At least 1 number (0-9)

**Backend validation** in `internal/service/auth_service.go`:
```go
func ValidateUsername(username string) error { ... }
func ValidatePassword(password string) error { ... }
```

**Frontend validation** in `pages/LoginPage.tsx` provides real-time feedback.

---

### 4. **Rate Limiting**
**Status:** ✅ FIXED

**Implementation:** IP-based rate limiting in `internal/middleware/ratelimit.go`

**Rules:**
- Maximum 5 attempts per 15 minutes per IP
- Applied to `/api/auth/register` and `/api/auth/login`
- Returns HTTP 429 (Too Many Requests)

**Usage:**
```go
// Applied automatically in router
http.HandleFunc("/api/auth/login", middleware.EnableCORS(
    middleware.RateLimit(handler.LoginHandler),
))
```

---

### 5. **Logout Endpoint**
**Status:** ✅ FIXED

**New Endpoint:** `POST /api/auth/logout`

**Features:**
- Revokes refresh token from database
- Requires valid access token (protected)
- Clears client-side tokens

**Request:**
```json
{
  "refresh_token": "eyJhbGc..."
}
```

**Frontend Function:**
```typescript
export const logoutUser = async (): Promise<void> => {
  // Calls /api/auth/logout to revoke server-side
  // Clears localStorage tokens
}
```

---

## 🟡 Recommended Future Improvements

### 1. **Use HttpOnly Cookies Instead of localStorage**
Currently using localStorage (XSS vulnerable). Consider:
- Set `HttpOnly` flag on cookies (JS cannot access)
- Set `Secure` flag (HTTPS only)
- Set `SameSite=Strict`

**Example Backend Code Needed:**
```go
http.SetCookie(w, &http.Cookie{
    Name:     "access_token",
    Value:    accessToken,
    HttpOnly: true,
    Secure:   true,  // HTTPS only
    SameSite: http.SameSiteLax,
    MaxAge:   900,   // 15 minutes
})
```

### 2. **HTTPS Enforcement**
- Use HTTPS in production (SSL/TLS certificate)
- Redirect HTTP to HTTPS
- Set `Strict-Transport-Security` header

### 3. **Additional Security Headers**
```go
w.Header().Set("X-Content-Type-Options", "nosniff")
w.Header().Set("X-Frame-Options", "DENY")
w.Header().Set("Content-Security-Policy", "default-src 'self'")
```

### 4. **Password Change Endpoint**
Implement `POST /api/auth/change-password` with:
- Old password verification
- New password validation
- Session invalidation

### 5. **Account Lockout**
After N failed login attempts, lock account for M minutes

### 6. **Audit Logging**
Log all authentication events (login, logout, token refresh)

### 7. **Two-Factor Authentication (2FA)**
Add TOTP/SMS verification for admin accounts

---

## 📋 Setup Instructions

### Development Environment
1. **Clone and setup:**
   ```bash
   cd X509-PKI
   ```

2. **Create `.env` file:**
   ```bash
   cp .env.example .env
   ```

3. **Edit `.env`:**
   ```env
   # Generate secure secret
   JWT_SECRET=<your-32-char-random-string>
   CORS_ORIGIN=http://localhost:5173
   SERVER_PORT=8080
   DB_PATH=data/users.db
   ```

4. **Start backend:**
   ```bash
   cd backend
   go run cmd/main.go
   ```

5. **Start frontend:**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

### Production Environment
1. **Build backend:**
   ```bash
   go build -o x509-pki cmd/main.go
   ```

2. **Build frontend:**
   ```bash
   npm run build
   ```

3. **Set environment variables** (use secure methods):
   - Docker: `.env` file or secrets
   - Kubernetes: ConfigMaps/Secrets
   - AWS: Systems Manager Parameter Store
   - Azure: Key Vault

4. **Enable HTTPS:**
   - Use reverse proxy (nginx, Caddy)
   - Install SSL certificate
   - Redirect HTTP → HTTPS

5. **Database:**
   - Use separate database with proper backups
   - Enable encryption at rest
   - Regular security audits

---

## 🔐 Security Checklist

- [ ] JWT_SECRET is random and 32+ characters
- [ ] .env file is NOT committed to Git
- [ ] CORS_ORIGIN matches your domain (production)
- [ ] HTTPS enabled (production)
- [ ] Database is backed up regularly
- [ ] Rate limiting is active
- [ ] Input validation enabled
- [ ] Logging and monitoring configured
- [ ] Regular security updates applied
- [ ] Penetration testing completed (optional)

---

## 🚨 Reporting Security Issues

If you discover a security vulnerability, please **DO NOT** create a public GitHub issue.

Instead:
1. Email security team privately
2. Include details of the vulnerability
3. Allow time for patching before public disclosure

---

## 📚 References

- [OWASP Password Storage Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
- [OWASP Authentication Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)
- [Argon2 Documentation](https://github.com/P-H-C/phc-winner-argon2)

---

**Last Updated:** April 19, 2026
**Version:** 2.0 (Security Hardening Release)
