# GitGo

GitGo is a simple, minimal version control system implemented in Go, inspired by Git. It allows tracking changes to a single file, committing snapshots, viewing commit logs, and managing basic branch references.

---

## Features

- Initialize GitGo repository for a single file.
- Commit file snapshots with descriptive messages.
- Maintain commit history with SHA-1 hashes.
- View commit logs with filters.
- Basic stash commands: list, apply, clear.
- Show repository status.
- Diff between files.
- Branch creation and listing.
- Revert to previous commits.

---

## Installation

Make sure you have [Go](https://golang.org/dl/) installed (version 1.18+ recommended).

Clone this repo and build:

```bash
git clone https://github.com/Anhadcoder/TrackIT.git
cd TrackIT
go build -o trackit main.go
