package config

// HttpConfig holds the configuration options for the HTTP server.
type HttpConfig struct {
	ReadTimeout   int
	WriteTimeout  int
	IdleTimeout   int
	MaxHeaderSize int
	TLSCertFile   string
	TLSKeyFile    string
	Security      SecurityConfig
	Session       SessionConfig
}

// SecurityConfig holds the security-related configuration options.
type SecurityConfig struct {
	EnableHSTS            bool
	ContentSecurityPolicy string
}

// SessionConfig holds the session-related configuration options.
type SessionConfig struct {
	Store   string
	Timeout int
}

// NewHttpConfig creates a new HttpConfig with optional settings.
// You can pass option functions to override default settings.
func NewHttpConfig(options ...func(*HttpConfig)) *HttpConfig {
	// Set default values
	cfg := &HttpConfig{
		ReadTimeout:   5,       // Default to 5 seconds
		WriteTimeout:  5,       // Default to 5 seconds
		IdleTimeout:   120,     // Default to 120 seconds
		MaxHeaderSize: 1 << 20, // Default to 1 MB
		TLSCertFile:   "",      // No default cert file
		TLSKeyFile:    "",      // No default key file
		Security: SecurityConfig{
			EnableHSTS:            true,                 // Default to enabling HSTS
			ContentSecurityPolicy: "default-src 'self'", // Default CSP
		},
		Session: SessionConfig{
			Store:   "memory", // Default to in-memory store
			Timeout: 3600,     // Default to 1 hour
		},
	}

	// Apply any provided options
	for _, option := range options {
		option(cfg)
	}

	return cfg
}

// Option functions

func WithReadTimeout(timeout int) func(*HttpConfig) {
	return func(cfg *HttpConfig) {
		cfg.ReadTimeout = timeout
	}
}

func WithWriteTimeout(timeout int) func(*HttpConfig) {
	return func(cfg *HttpConfig) {
		cfg.WriteTimeout = timeout
	}
}

func WithIdleTimeout(timeout int) func(*HttpConfig) {
	return func(cfg *HttpConfig) {
		cfg.IdleTimeout = timeout
	}
}

func WithMaxHeaderSize(size int) func(*HttpConfig) {
	return func(cfg *HttpConfig) {
		cfg.MaxHeaderSize = size
	}
}

func WithTLSCertFile(certFile string) func(*HttpConfig) {
	return func(cfg *HttpConfig) {
		cfg.TLSCertFile = certFile
	}
}

func WithTLSKeyFile(keyFile string) func(*HttpConfig) {
	return func(cfg *HttpConfig) {
		cfg.TLSKeyFile = keyFile
	}
}

func WithHSTS(enabled bool) func(*HttpConfig) {
	return func(cfg *HttpConfig) {
		cfg.Security.EnableHSTS = enabled
	}
}

func WithContentSecurityPolicy(policy string) func(*HttpConfig) {
	return func(cfg *HttpConfig) {
		cfg.Security.ContentSecurityPolicy = policy
	}
}

func WithSessionStore(store string) func(*HttpConfig) {
	return func(cfg *HttpConfig) {
		cfg.Session.Store = store
	}
}

func WithSessionTimeout(timeout int) func(*HttpConfig) {
	return func(cfg *HttpConfig) {
		cfg.Session.Timeout = timeout
	}
}
