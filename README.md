# One Time Secret Service

A secure web service for storing and sharing one-time secrets (passwords, API keys, confidential information) with automatic deletion after first access.

## 📋 Overview

**One Time Secret Service** is a high-performance Go application that allows users to:
- Create encrypted secrets with unique links
- Securely share confidential information
- Guarantee single-use access to data
- Automatically delete secrets after retrieval

Data is protected using **AES-256-GCM** cryptography, with the encryption key embedded in the URL and never stored on the server.

## 🎯 Key Features

✅ **Single-Use Access** — secrets are deleted immediately after retrieval  
✅ **AES-256-GCM Encryption** — cryptographic data protection  
✅ **Key Embedded in URL** — no keys stored on server  
✅ **Fast Web Framework** — Gofiber for high performance  
✅ **Easy Deployment** — SQLite in-memory database, no complex dependencies  

## 🛠️ Technology Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.25.0 |
| Framework | Gofiber v3 |
| Database | SQLite + GORM |
| Cryptography | AES-GCM |
| UUID | google/uuid |

## 📦 Project Structure

```
.
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── cypher/
│   │   └── base.go          # Cryptographic module (AES-GCM)
│   └── storage/
│       ├── base.go          # Storage interface
│       └── database.go      # GORM implementation
├── templates/
│   ├── index.html           # Home page (secret creation)
│   └── get_secret.html      # Secret retrieval page
├── go.mod                   # Go module
└── README.md                # This file
```

## 🚀 Installation and Setup

### Requirements
- Go 1.25.0 or higher

### Install Dependencies
```bash
go mod download
```

### Run the Server
```bash
go run ./cmd/main.go
```

The server will run on `http://localhost:3000`

### Build
```bash
go build -o one_time_secret_service ./cmd/main.go
./one_time_secret_service
```

## 📡 API Endpoints

### 1. Home Page
```
GET /
```
Returns HTML page for creating secrets.

### 2. Create Secret
```
POST /api/create
```

**Parameters:**
- `secret` (form-data) — secret text to encrypt

**Response:**
```json
{
  "id": "3f4a2c1b-a1b2c3d4e5f6g7h8-hex_encrypted_key"
}
```

**Example:**
```bash
curl -X POST http://localhost:3000/api/create \
  -d "secret=my_password_123"
```

### 3. Retrieve Secret
```
POST /api/get
```

**Parameters:**
- `id` (form-data) — secret identifier (obtained during creation)

**Response:**
```json
{
  "secret": "my_password_123"
}
```

**Example:**
```bash
curl -X POST http://localhost:3000/api/get \
  -d "id=3f4a2c1b-a1b2c3d4e5f6g7h8-hex_encrypted_key"
```

### 4. Secret Retrieval Page
```
GET /s/:id
```
Returns HTML page for retrieving a secret by ID.

## 🔒 How Encryption Works

1. **Key Generation**: When creating a secret, a random 128-bit AES key is generated
2. **Encryption**: Data is encrypted using AES-GCM with a random nonce
3. **Key Embedding**: The key is hex-encoded and embedded in the URL identifier
4. **Storage**: Only encrypted data is stored on the server (without the key)
5. **Decryption**: Client sends the URL with the key, server decrypts and deletes the data
6. **Single-Use**: The secret is deleted from the database immediately after successful retrieval

## 💡 Usage Examples

### Web Interface
1. Open `http://localhost:3000/`
2. Enter your secret in the form
3. Copy the generated link
4. Share the link with the recipient
5. When the recipient visits the link, the secret is displayed and deleted

### Using cURL
```bash
# Create secret
RESPONSE=$(curl -s -X POST http://localhost:3000/api/create \
  -d "secret=super_secret_password")
ID=$(echo $RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)

# Retrieve secret
curl -X POST http://localhost:3000/api/get \
  -d "id=$ID"
```

## ⚠️ Important Notes

- **Single-Use Access**: After the first retrieval, the secret is deleted and no longer accessible
- **Temporary Storage**: By default, SQLite in-memory database is used (data is lost on restart)
- **URL Security**: The entire URL contains the encryption key — do not log it and transmit over secure channels
- **Scalability**: For production use PostgreSQL and deploy multiple instances

## 📝 License

See the [LICENSE](LICENSE) file

