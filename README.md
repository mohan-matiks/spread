<div align="center">
  <img src="https://github.com/user-attachments/assets/c1c06796-43b3-47a3-a748-3cdc3fe75358" width="500" />
  <p><strong>Spread your delightful releases to the world!✨</strong></p>
  
  [![Go Version](https://img.shields.io/badge/Go-1.22+-blue.svg)](https://golang.org)
  [![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
  [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://github.com/SwishHQ/spread) 
</div>

## Overview

Spread is an OTA (Over-the-Air) update server designed specifically for React Native applications. It enables developers to push JavaScript bundle updates to their React Native apps without requiring users to download new versions from app stores. A [react-native-code-push](https://github.com/microsoft/react-native-code-push) compatible server.

Spread uses Cloudflare R2 bucket for storing bundles. You can learn more about Cloudflare R2 [here](https://developers.cloudflare.com/r2/)

You can use [this](https://aexomir1.medium.com/configuring-react-native-code-push-using-custom-server-e40e87697a26) article to set Spread Host URL in react-native-code-push.

### Key Features

- **Fully Self-hostable** - Complete control over your update infrastructure
- **CodePush compatible** - Drop-in replacement for Microsoft CodePush
- **Multi-platform support** - iOS and Android bundle management
- **Environment management** - Separate configurations for development, staging, and production
- **Version control** - Track and manage different app versions
- **Rollback capabilities** - Quickly revert to previous versions
- **Web dashboard** - React-based admin interface

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/SwishHQ/spread.git
cd spread
```

### 2. Environment Configuration

Copy the example environment file and configure your settings:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
ENV=local
APP_NAME=spread
PORT=4000

# MongoDB Configuration
MONGODB_URL=mongodb://localhost:27017
MONGODB_DATABASE=spread

# Cloudflare R2 Configuration (for bundle storage)
CLOUDFLARE_R2_ACCOUNT_ID=your_account_id
CLOUDFLARE_R2_BUCKET=your_bucket_name
CLOUDFLARE_R2_ACCESS_KEY_ID=your_access_key
CLOUDFLARE_R2_SECRET_ACCESS_KEY=your_secret_key
```

### 3. Install Dependencies

Install Go dependencies:
```bash
go mod download
```

Install frontend dependencies:
```bash
cd web
npm install
cd ..
```

### 4. Build the Application

Build the backend:
```bash
make build
```

Build the frontend:
```bash
cd web
npm run build
cd ..
```

### 5. Run the Server

Start the Spread server:
```bash
./spread serve
```

The server will be available at `http://localhost:4000`

## Project Structure
Project Structure

```
├── cmd/                # Command-line interface (CLI) entrypoints
│   ├── root.go         # Root command configuration
│   ├── serve.go        # Server command
│   └── client.go       # Client/release commands
├── cli/                # CLI utilities
│   └── bundle_cli.go   # Bundle management CLI
├── config/             # Configuration management
├── middleware/         # HTTP middleware (auth, logging, etc.)
├── pkg/                # External package integrations
│   ├── cloudflare.go   # Cloudflare R2 integration
│   └── db.go           # Database connection
├── src/                # Main application source code
│   ├── controller/     # HTTP controllers (route handlers)
│   ├── model/          # Data models and schemas
│   ├── repository/     # Data access layer
│   └── service/        # Business logic and orchestration
├── types/              # Shared type definitions
├── utils/              # Utility/helper functions
├── logger/             # Logging configuration and utilities
├── exception/          # Centralized error handling
├── web/                # Frontend React application
│   ├── src/            # React source code
│   ├── public/         # Static assets
│   ├── dist/           # Built frontend assets
│   └── package.json    # Frontend dependencies
├── script/             # Build and deployment scripts
├── main.go             # Application entry point
├── go.mod              # Go module definition
├── go.sum              # Go dependency checksums
├── Makefile            # Build automation
├── Dockerfile          # Container configuration
├── .env.example        # Environment variable template
└── README.md           # This file (the one you’re reading)
```
Pro tip: Each directory is lovingly crafted to keep things decoupled and maintainable. Whether you’re a backend buff, a frontend fanatic, or just here for the scripts, you’ll find your happy place.
## 🔧 Development Setup

### Local Development

1. **Start MongoDB** (if not already running):
   ```bash
   # macOS with Homebrew
   brew services start mongodb-community
   
   # Ubuntu/Debian
   sudo systemctl start mongod
   
   # Windows
   net start MongoDB
   ```

2. **Run in development mode**:
   ```bash
   # Terminal 1: Run backend
   go run *.go serve
   
   # Terminal 2: Run frontend (in web directory)
   cd web
   npm run dev
   ```

3. **Access the application**:
   - Backend API: `http://localhost:4000`
   - Frontend: `http://localhost:5173` (Vite dev server)

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `ENV` | Environment (local, development, production) | `local` | Yes |
| `APP_NAME` | Application name | `spread` | No |
| `PORT` | Server port | `4000` | No |
| `MONGODB_URL` | MongoDB connection string | - | Yes |
| `MONGODB_DATABASE` | MongoDB database name | `spread` | Yes |
| `CLOUDFLARE_R2_ACCOUNT_ID` | Cloudflare R2 account ID | - | Yes |
| `CLOUDFLARE_R2_BUCKET` | Cloudflare R2 bucket name | - | Yes |
| `CLOUDFLARE_R2_ACCESS_KEY_ID` | Cloudflare R2 access key | - | Yes |
| `CLOUDFLARE_R2_SECRET_ACCESS_KEY` | Cloudflare R2 secret key | - | Yes |

## 🛠️ Building and Deployment

### Building from Source

```bash
# Build the entire application
make build

# Build with specific GOOS and GOARCH
GOOS=linux GOARCH=amd64 go build -o spread-linux-amd64
GOOS=darwin GOARCH=amd64 go build -o spread-darwin-amd64
GOOS=windows GOARCH=amd64 go build -o spread-windows-amd64.exe
```

### Docker Deployment

```bash
# Build Docker image
docker build -t spread .

# Run with Docker
docker run -p 4000:4000 --env-file .env spread
```

### Production Deployment

1. Set up a production MongoDB instance
2. Configure Cloudflare R2 storage
3. Set environment variables for production
4. Build and deploy the application

Example production deployment with Docker Compose:

```yaml
version: '3.8'
services:
  spread:
    build: .
    ports:
      - "4000:4000"
    environment:
      - ENV=production
      - MONGODB_URL=mongodb://mongo:27017
      - MONGODB_DATABASE=spread
      - CLOUDFLARE_R2_ACCOUNT_ID=${CLOUDFLARE_R2_ACCOUNT_ID}
      - CLOUDFLARE_R2_BUCKET=${CLOUDFLARE_R2_BUCKET}
      - CLOUDFLARE_R2_ACCESS_KEY_ID=${CLOUDFLARE_R2_ACCESS_KEY_ID}
      - CLOUDFLARE_R2_SECRET_ACCESS_KEY=${CLOUDFLARE_R2_SECRET_ACCESS_KEY}
    depends_on:
      - mongo
    restart: unless-stopped

  mongo:
    image: mongo:6.0
    ports:
      - "27017:27017"
    volumes:
      - mongo_data:/data/db
    restart: unless-stopped

volumes:
  mongo_data:
```

## Using the CLI

### Installation

Install the Spread CLI globally:

```bash
# From source
go install github.com/SwishHQ/spread/cmd/spread@latest

# Or build locally
make build
sudo cp spread /usr/local/bin/

# Or use the install script
curl -fsSL https://cdn-swish.justswish.in/spread-install.sh | sh
```
### Creating a Release

```bash
spread release \
  --remote https://your-spread-server.com \
  --auth-key YOUR_AUTH_KEY \
  --app-name my-react-native-app \
  --environment production \
  --target-version 1.2.0 \
  --os-name ios \
  --project-dir /path/to/react-native/project \
  --is-typescript true \
  --description "Bug fixes and performance improvements"
```

### CLI Options

| Flag | Description | Required | Default |
|------|-------------|----------|---------|
| `--remote` | Spread server URL | Yes | - |
| `--auth-key` | Authentication key | Yes | - |
| `--app-name` | Application name | Yes | - |
| `--environment` | Environment (development, staging, production) | Yes | - |
| `--target-version` | Target app version | Yes | - |
| `--os-name` | Operating system (ios, android) | Yes | - |
| `--project-dir` | React Native project directory | No | Current directory |
| `--is-typescript` | Is TypeScript project | No | false |
| `--description` | Release description | No | - |
| `--disable-minify` | Disable bundle minification | No | false |
| `--hermes` | Enable Hermes engine | No | false |


## Contributing

We welcome contributions from the community! Here's how you can help:

### Development Workflow

1. **Fork the repository**
2. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. **Make your changes**
4. **Run tests**:
   ```bash
   go test ./...
   ```
5. **Commit your changes**:
   ```bash
   git commit -m "feat: add your feature description"
   ```
6. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```
7. **Create a Pull Request**

### Code Style Guidelines

- Follow Go conventions and use `gofmt` for formatting
- Write meaningful commit messages following conventional commits
- Add tests for new functionality
- Update documentation for API changes
- Ensure all tests pass before submitting PR

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./src/service -v
```


<div align="center">
  <p>Made with 💚 by the Swish Engineering</p>
</div>

