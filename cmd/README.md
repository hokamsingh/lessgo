# lessgo Framework

`lessgo` is a lightweight, modular HTTP framework built in Go, inspired by the architecture of NestJS. It aims to provide developers with a structured and scalable way to build server-side applications while leveraging the powerful features of Go.

## Features

- **Modular Design:** Organize your application into reusable modules for better maintainability.
- **Dependency Injection (DI):** Simplify dependency management and testing with built-in DI support.
- **Middleware Support:** Easily apply middleware to routes for handling cross-cutting concerns.
- **Routing:** Utilize powerful and flexible routing with Gorilla Mux.
- **Error Handling:** Centralized error handling to manage errors gracefully.
- **Configuration Management:** Load and manage application configurations via `.env` files.
- **Service Layer:** Separate business logic from controllers using services.

## Getting Started

### Installation

Clone the `lessgo` repository and include it in your Go project.

### Project Structure

```bash
git clone https://github.com/yourusername/lessgo.git
cd lessgo

lessgo/
├── app/
│   ├── middleware/          # Custom middleware implementations
│   └── module/              # Application modules
│       └── test/            # Example test module
├── cmd/                     # Main application entry point
├── internal/
│   ├── core/                # Core framework components
│   │   ├── config/          # Configuration management
│   │   ├── controller/      # Base controller interfaces and logic
│   │   ├── di/              # Dependency Injection container
│   │   ├── module/          # Module registration and management
│   │   └── router/          # Router initialization and setup
├── README.md                # This README file
└── .env                     # Environment configuration file.
```

### Config

```env
SERVER_PORT=8080
ENV=development
JWT_SECRET=mysecretkey
```

### Run 

``bash
go run cmd/main.go
```