# Trashcan Project

Go-based backend service with Gin framework, WebSocket support, and SQLite database.

## Tech Stack

- **Framework**: Gin (HTTP router)
- **Database**: SQLite (via glebarez/sqlite + GORM)
- **Auth**: JWT (golang-jwt/jwt/v5)
- **WebSocket**: gws (lxzan/gws)
- **Logging**: Zap
- **Config**: Viper

## Project Structure

- `ginServer/` - HTTP API handlers, middleware, models, router
- `internal/modules/` - WebSocket and FRP modules
- `initialize/` - Startup initialization (DB, config, logging)
- `utils/` - Shared utilities (JWT, UUID, image processing)
- `scripts/` - Database initialization and test data
- `config/` - Configuration structures
- `global/` - Global state

## Coding Guidelines (Karpathy)

### 1. Think Before Coding

- State assumptions explicitly. If uncertain, ask.
- Present multiple interpretations if they exist.
- Surface simpler approaches when they exist.
- Stop and ask when something is unclear.

### 2. Simplicity First

- Write minimum code that solves the problem.
- No features beyond what was asked.
- No abstractions for single-use code.
- No error handling for impossible scenarios.
- If 200 lines could be 50, rewrite it.

### 3. Surgical Changes

When editing existing code:
- Touch only what you must.
- Don't "improve" adjacent code or formatting.
- Match existing style.
- Remove only imports/variables/functions YOUR changes made unused.

### 4. Goal-Driven Execution

Transform tasks into verifiable goals:
- "Add validation" → "Write tests for invalid inputs, then make them pass"
- "Fix the bug" → "Write a test that reproduces it, then make it pass"

For multi-step tasks, state a brief plan with verification steps.
