### `config` Package Documentation

#### Overview:
The `config` package is responsible for loading and managing application configuration from environment variables and `.env` files.

#### Functions:

- **`LoadConfig()`**: Loads configuration from `.env` file (if found) and environment variables into a `Config` map. Searches for `.env` in the current or root directory.

- **`Get()`**: Retrieves a string value from the config, with a fallback default.

- **`GetInt()` / `GetBool()` / `GetFloat64()`**: Retrieves integer, boolean, or float64 values with fallback defaults, logging errors for invalid values.

- **`Validate()`**: Ensures required keys are present, logging a fatal error if any are missing.

- **`Reload()`**: Reloads the configuration, useful for dynamic environments.

- **`MergeWithDefaults()`**: Combines the current config with defaults, prioritizing current values.

- **`FilterByPrefix()`**: Returns a new config map with keys matching a given prefix, removing the prefix from keys.

- **`findRootDir()`**: Locates the root directory containing the `.env` file, starting from the current directory.

#### Usage:
1. **Loading Configuration:**
   ```go
   cfg := LoadConfig()
   ```

2. **Accessing Configuration:**
   ```go
   port := cfg.Get("PORT", "8080")
   timeout := cfg.GetInt("TIMEOUT", 30)
   ```

3. **Validating Configuration:**
   ```go
   cfg.Validate("DB_HOST", "DB_USER", "DB_PASS")
   ```

4. **Merging with Defaults:**
   ```go
   defaultConfig := Config{"PORT": "8080", "ENV": "development"}
   cfg = cfg.MergeWithDefaults(defaultConfig)
   ```

5. **Filtering by Prefix:**
   ```go
   dbConfig := cfg.FilterByPrefix("DB_")
   ```

### Error Handling:
The package provides error logging for invalid data types and missing configurations, ensuring that the application fails gracefully.

### Summary:
This package is a flexible, extendable configuration manager for Go applications, allowing for dynamic environment management, easy access to configuration variables, and integration of defaults and prefixes for modular setups.