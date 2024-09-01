---
![License](https://img.shields.io/badge/license-MIT-green.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/hokamsingh/lessgo)](https://pkg.go.dev/github.com/hokamsingh/lessgo)
[![Go Version](https://img.shields.io/github/go-mod/go-version/hokamsingh/lessgo-cli)](https://golang.org/dl/)
![Version](https://img.shields.io/badge/version-v1.0.2-blue)
[![Go Report Card](https://goreportcard.com/badge/github.com/hokamsingh/lessgo)](https://goreportcard.com/report/github.com/hokamsingh/lessgo)

**🚀 LessGo v1.0.2 - The High-Performance Minimal Go Web Framework**

Welcome to **LessGo** ⚡—a high-performance, minimal Go web framework crafted for building scalable, maintainable, and modern applications. LessGo empowers developers with advanced features like dynamic configuration management, inbuilt error handling, robust data validation, and Dependency Injection (DI).

Inspired by the modular architecture of frameworks like NestJs, LessGo offers a flexible structure that allows you to build your applications your way. Whether you prefer a strict controller-service-module architecture or a more fluid design, LessGo has you covered. With built-in support for CORS, CSRF, XSS protection, logging, error handling, rate limiting, caching, and job scheduling, LessGo makes it easier to develop secure, efficient, and performant web applications.

### ✨ Key Features:

1. **🔄 Concurrency & Task Management**  
   LessGo includes a powerful task manager that allows you to run tasks in parallel or sequence, giving you full control over task management and program flow. We are committed to enhancing its robustness and extensibility.

2. **⚙️ Dynamic Configuration Management**  
   With built-in support for multiple environments (development, production, testing), LessGo simplifies configuration management. It provides a user-friendly API for working with environment variables, offering validation and easy access to typed values like numbers and booleans.

3. **🔧 LessGo Context**  
   The LessGo Context is a powerful abstraction over the request and response objects. It simplifies API development by providing methods for handling headers, cookies, query parameters, JSON parsing, and more. Error management, HTTP redirection, and file attachments are made effortless.

4. **🧩 Flexible Dependency Injection (DI)**  
   LessGo offers basic DI functionality that we're continuously improving. You can choose whether to bind your entire application into a single container or work with a more traditional approach, giving you flexibility in how you manage dependencies.

5. **🚀 App Initialization**  
   The App is the core of LessGo, initializing various configurations and middleware. It exposes the main server with a simple `Listen` method to start your application.

6. **⏰ Job Scheduling**  
   LessGo supports job scheduling, leveraging the powerful Cron package to handle recurring tasks seamlessly.

7. **🔒 Comprehensive Middleware**  
   LessGo comes with essential built-in middleware, including CORS, JSON parser, cookie parser, CSRF and XSS protection, caching, file uploads, and rate limiting (both in-memory and Redis-backed).

8. **🌐 Powerful Router**  
   The built-in router handles different HTTP methods with custom handlers, supporting method chaining, sub-routers, and custom middleware.

9. **📦 Controller Interface**  
   Controllers in LessGo allow you to create dedicated routers for specific endpoints, integrating seamlessly with the service layer for efficient request handling.

10. **🔍 Service Layer**  
   The Service interface in LessGo is designed for encapsulating business logic. It easily binds to databases and other dependencies, streamlining the development process.

11. **🔗 Modular Design**  
   Modules in LessGo bring together controllers and services, managing routes and handlers in a cohesive manner.

12. **💬 WebSocket Support (Upcoming)**  
   We are actively working on integrating WebSocket support to make real-time communication in your applications even easier.

---

### 📦 Installation & Upgrade

To get started with the latest version of LessGo, update your dependencies using:

```sh
go get github.com/hokamsingh/lessgo@latest
```

### 🌟 Get Started Quickly with LessGo CLI

We're also introducing the **LessGo CLI**, a command-line tool to help you scaffold and manage your LessGo projects with ease! With the CLI, you can:

- **Create a New Project**: Quickly set up a new LessGo project using `lessgo-cli new myapp`.
- **Check Version**: Keep track of the CLI version with `lessgo-cli --version`.
- **Cross-Platform Support**: Works seamlessly on both Windows and Unix-based systems.

Install the LessGo CLI with:

```sh
go install github.com/hokamsingh/lessgo-cli@latest
```

Make sure to try out the CLI to streamline your project setup and start building with LessGo in no time!

### 🙌 Acknowledgments

We would like to thank our contributors and community for their support and feedback. Your contributions have been invaluable in shaping the LessGo framework.

### 📜 License

This project is licensed under the MIT License - see the LICENSE file for details.

For any questions, issues, or feedback, please visit our [GitHub repository](https://github.com/hokamsingh/lessgo) or join our community discussions.

**Happy coding!** 🎉

--- 
