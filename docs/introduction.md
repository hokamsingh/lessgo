---
![License](https://img.shields.io/badge/license-MIT-green.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/hokamsingh/lessgo)](https://pkg.go.dev/github.com/hokamsingh/lessgo)
[![Go Version](https://img.shields.io/github/go-mod/go-version/hokamsingh/lessgo-cli)](https://golang.org/dl/)
![Version](https://img.shields.io/badge/version-v1.0.2-blue)
[![Go Report Card](https://goreportcard.com/badge/github.com/hokamsingh/lessgo)](https://goreportcard.com/report/github.com/hokamsingh/lessgo)

## 🛠️ LessGo Framework Release

We are excited to announce the latest release of the **LessGo** framework! This release introduces several powerful features and enhancements to make your Go development experience even more robust and efficient.

### 🚀 New Features

- **Enhanced Dependency Injection (DI)**:
  - **Conditional Bindings**: Register dependencies based on specific conditions.
  - **Scoped Services**: Manage services with different lifetimes and scopes.
  - **Named Dependencies**: Support for naming dependencies for more granular control.

- **Dynamic Configuration Management**:
  - **Factory Pattern**: Implement dynamic configuration management with factories to handle various application environments and requirements.

- **Inbuilt Error Handling**:
  - **Custom Error Responses**: Integrate a flexible error handling system that returns appropriate HTTP responses based on the error type.
  - **Detailed Logging**: Capture detailed error logs and stack traces for better debugging.

- **Job Scheduler Integration**:
  - **Built-in Scheduler**: Utilize an advanced job scheduling library for managing periodic tasks and background jobs seamlessly.

- **JSON Response Handling**:
  - **Automatic JSON Parsing**: Simplify the process of parsing and responding with JSON data.
  - **Custom JSON Encoding**: Encode and decode JSON responses with integrated error handling.

- **Advanced Data Validation**:
  - **Strict Validation**: Ensure data integrity by validating incoming requests against predefined models.
  - **Dynamic Key Checks**: Verify that all required fields are present in the request payload.

- **🛡️ CSRF Protection Middleware**:
  - **`LessGo.WithCsrf()`**: Integrated Cross-Site Request Forgery (CSRF) protection to safeguard your applications against unauthorized actions.

- **🔒 XSS Protection Middleware**:
  - **`LessGo.WithXss()`**: Built-in Cross-Site Scripting (XSS) protection to prevent malicious script injections.

- **⚡ Caching Middleware**:
  - **`LessGo.WithCaching("localhost:6379", 5*time.Minute)`**: Redis-based caching middleware for improved performance with configurable cache expiration
- **Response Handling**
  - Enhanced response handling to ensure multiple responses are not sent from the same context, preventing unexpected behaviors and potential crashes.
  
- **Template Rendering**
  - Introduced support for rendering templates dynamically from a specified directory, reducing the need to manually list each template file.

- **Middleware Improvements**
  - Refined existing middleware to ensure seamless integration and performance.
  - Addressed potential edge cases and improved overall robustness.

## ⚙️ Migration Notes

- The new security middlewares (`WithCsrf` and `WithXss`) are optional but highly recommended for applications that handle sensitive data.
- The caching middleware requires a running Redis instance. Ensure your Redis server is properly configured and accessible from your application.

## 🛡️ Security

- Security is a top priority for us, and this release emphasizes that with the introduction of CSRF and XSS protection middlewares. We strongly recommend enabling these features to safeguard your applications against common web vulnerabilities.


### 📚 Documentation & Examples

- **Enhanced Documentation**: Comprehensive guides and examples to help you get started with the new features and integrations.
- **Code Examples**: Practical examples demonstrating how to use the new features in real-world applications.

### 🛠️ Improvements & Fixes

- **Performance Enhancements**: Optimizations for better performance and reduced resource usage.
- **Bug Fixes**: Various bug fixes and stability improvements to ensure a smoother development experience.

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

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.

---

For any questions, issues, or feedback, please visit our [GitHub repository](https://github.com/hokamsingh/lessgo) or join our [community discussions](#).

Happy coding!

---
