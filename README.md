# cronlens

A terminal UI for inspecting, testing, and visualizing cron job schedules and their next execution times.

---

## Installation

```bash
go install github.com/yourusername/cronlens@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/cronlens.git
cd cronlens
go build -o cronlens .
```

---

## Usage

Launch the interactive TUI:

```bash
cronlens
```

Pass a cron expression directly to inspect it:

```bash
cronlens "*/5 * * * *"
```

### In the TUI

- Type or paste a cron expression into the input field
- View the next 10 scheduled execution times in real time
- Use arrow keys to browse and compare multiple expressions
- Press `q` or `Ctrl+C` to quit

---

## Features

- Instant parsing and validation of cron expressions
- Visual timeline of upcoming executions
- Support for standard 5-field and extended 6-field (with seconds) formats
- Human-readable schedule descriptions
- Lightweight with no external dependencies beyond the standard library

---

## Requirements

- Go 1.21 or later

---

## License

MIT © 2024 yourusername