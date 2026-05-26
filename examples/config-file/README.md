# Config File Example

Demonstrates config file loading with automatic precedence over tag defaults.

## Running

```bash
# With defaults
go run .

# With config file
mkdir -p ~/.config/greet
echo '{"name": "Config", "count": 3}' > ~/.config/greet/config.json
go run .

# Flag overrides config file
go run . hello --name Flag --count 5

# Custom config file path
go run . hello --config /path/to/config.json
```

## Precedence

1. Explicit flags (`--name`, `--count`)
2. Environment variables (`env:"VAR"` tags)
3. Config file values
4. Tag defaults (`default:"..."`)
