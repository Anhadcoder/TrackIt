# TrackIT

A simple Git-like version control system for tracking a **single file**, built in Go.

---

## Overview

TrackIT is a lightweight, command-line file tracker inspired by Git. It supports fundamental version control operations like initializing a tracker, committing changes, managing branches, viewing logs, stashing, reverting commits, and diffing files — all for a **single file** in a directory.

It is designed to demonstrate the core concepts of version control systems, with a focus on ease of use and performance for individual file tracking.

---

## Features

- **Initialize** a tracker for a single file  
- **Commit** changes with messages and unique SHA-1 hashes  
- View commit **log** with commit history  
- **Stash** and **apply** changes (basic stash support)  
- View current **status** of tracked file  
- **Revert** to any previous commit  
- Create and **switch branches**  
- **Get** file content at any commit  
- **Diff** between two files to view changes  

---

## Installation

Make sure you have [Go](https://golang.org/dl/) installed (version 1.18+ recommended).

Clone this repo and build:

```bash
git clone https://github.com/yourusername/TrackIT.git
cd TrackIT
go build -o trackit main.go
