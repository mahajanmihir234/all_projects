# Schedule JSON files

Each `*.json` file in this directory defines tasks loaded at startup. Fields:

| Field | Description |
|-------|-------------|
| `id` | Task identifier (printed on run) |
| `task_type` | `print` (default) or `fail` (simulated error) |
| `schedule.type` | `once`, `interval`, or `cron` |

## Schedule types

**once** — run a single time  
- `in`: duration from load time, e.g. `"1s"`, `"1500ms"`  
- `at`: RFC3339 absolute time, e.g. `"2026-05-19T15:00:00Z"`

**interval** — fixed delay between runs  
- `every`: interval, e.g. `"500ms"`  
- `start_in`: optional delay before first run

**cron** — standard cron (optional seconds), via [robfig/cron](https://github.com/robfig/cron)  
- `expression`: e.g. `"*/2 * * * * *"` (every 2s), `"@every 4s"`

Run from the `task_scheduler` directory:

```bash
go run .
go run . /path/to/other/schedules
```
