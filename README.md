# X.509 PKI Dashboard Web Application

## Overview
This is a secure internal Public Key Infrastructure (PKI) web dashboard managing X.509 certificates. The application features user authentication using high-security **Argon2id** password hashing and **JWT** session tokens. It serves different customized and responsive UI dashboards depending on Role-Based Access Control (RBAC): `Admin` and `Client`.

## Tech Stack
-   **Frontend:** React, TypeScript, Vite, Pure CSS (Glassmorphism + Dark Mode themes).
-   **Backend:** Go (Golang), `net/http`.
-   **Database:** SQLite (`modernc.org/sqlite` without CGO).
-   **Security:** Argon2id (OWASP standard), SHA-256 for token storage, HttpOnly cookies for JWT refresh mechanisms.

## Project Structure
-   `/frontend`: Contains the SPA built entirely with `React` configured by `Vite`.
    -   `App.tsx`: Routes user sessions to `<AdminDashboardPage>` or `<DashboardPage>`.
    -   `LoginPage.tsx`: Provides account registration/login capabilities powered by fluid component transitions.
-   `/backend`: Exposes the central API.
    -   `/cmd`: Run entry-point `main.go`.
    -   `/internal`: Defines domain models, database abstractions, JWT and Argon2id cryptographic functions.

## Authentication Mechanisms
- By default, whenever you boot up the web service, the backend will auto-seed a primary administrator account into the SQLite table if it is missing:
    -   **Username:** `admin`
    -   **Password:** `Admin@x509-pki`
- Registration forms used from the Application will exclusively issue **Client** accounts.

## Running the Application

### 0. Setup Environment Variables (First Time Only)
```bash
# Copy the example .env file
cp .env.example .env

# Edit .env and set a strong JWT_SECRET (min 32 characters)
# You can generate one with:
# openssl rand -base64 32
```

**Important:** Never commit `.env` to Git (already in `.gitignore`).

### 1. Boot the Backend Server
```bash
cd backend
go run cmd/main.go
```
The Go server will start up `http://localhost:8080` and provision `/data/users.db`.

**You should see:**
```
✅ SQLite DB ready at: data/users.db
✅ Seeded admin user (admin / Admin@x509-pki)
📌 API List:
POST /api/auth/register → Register User
...
🚀 Server running at http://localhost:8080
```

### 2. Boot the Frontend Client
(In a separate terminal)
```bash
cd frontend
npm install
npm run dev
```
Navigate to the provided `http://localhost:5173` URI.
You can immediately login with `admin` to preview the admin features.

---

## 🔒 Security Features

✅ **Argon2id Password Hashing** - OWASP recommended algorithm  
✅ **JWT Token Rotation** - Single-use refresh tokens  
✅ **Rate Limiting** - 5 attempts per 15 minutes on login/register  
✅ **Strong Input Validation** - Username & password requirements  
✅ **CORS Protection** - Configurable per environment  
✅ **Token Revocation** - Logout endpoint to revoke sessions  

For detailed security information, see [SECURITY.md](SECURITY.md).

---
