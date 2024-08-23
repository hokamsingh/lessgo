# Define the source directory and output directory
SRC_DIR = app/src
OUT_DIR = app/plugins

# Ensure the output directory exists
$(OUT_DIR):
	mkdir -p $(OUT_DIR)

# Find all _module.go files and create corresponding .so files
MODULES = $(patsubst $(SRC_DIR)/%.go,$(OUT_DIR)/%.so,$(wildcard $(SRC_DIR)/*_module.go))

# Default target
all: $(OUT_DIR) $(MODULES)

# Rule to compile Go files into shared objects
$(OUT_DIR)/%.so: $(SRC_DIR)/%.go
	go build -buildmode=plugin -o $@ $<

# Clean up compiled files
clean:
	rm -f $(OUT_DIR)/*.so
