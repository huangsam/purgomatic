# Purgomatic 💎

High-performance auditing for large-scale photo/video migrations.

## Concept
Purgomatic identifies **"Golden Winners"** (safely archived assets) to help you eliminate **"Toil"** (redundant copies). It provide clear strategic advice on what to move and where.

### Features
- **Zero CGO**: Pure Go using `sqlite`.
- **Hardware Scaled**: Dynamic concurrency via `runtime.NumCPU()`.
- **Multi-Point Hashing**: Sparse SHA-256 sampling for speed and safety.
- **Stat-First Sync**: Skips unchanged files via metadata.
- **Worst Offenders**: Flags the top 3 largest files per year.

## Setup
1. **Install**: `go install github.com/huangsam/purgomatic@latest`
2. **Configure**: Create `scans.json`:
```json
[{ "source": "Synology", "path": "/Volumes/Photos" }, { "source": "Phone", "path": "/Backup" }]
```
3. **Run**:
```bash
make init   # Initialize database
make audit  # Scan and report
```
