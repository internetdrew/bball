# bball (NBA CLI)

Catch up on NBA games from your terminal. Fast, simple, and just for fun.

This side project pulls live scoreboard data and season schedules from the NBA's public JSON endpoints and formats them nicely for the command line.

> Hobby project: no official releases or code signing. Build from source or `go install` to use.

---

## Features

- View today's NBA games: status, scores, and clocks
- Filter for only live or only final games
- See a team's upcoming or recent games
- Quick "catch-up" summary for a team's current (live) game

---

## Install

Requires Go 1.20+.

### Easiest: install via `go install`

```bash
go install github.com/internetdrew/bball@latest
```

Make sure Go's bin directory is on your PATH:

```bash
# Add GOPATH/bin (or GOBIN) to PATH for zsh
echo 'export PATH="$(go env GOPATH)/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc

# Verify
which bball
```

### Build from source

```bash
git clone https://github.com/internetdrew/bball.git
cd bball
go build -o bball .

# Optionally move it somewhere on your PATH
# e.g., ~/bin avoids sudo:
mkdir -p ~/bin
mv bball ~/bin
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

**Note for macOS:** Gatekeeper generally allows locally built binaries. Downloaded binaries from unknown sources may require manual approval or recompilation.

---

## Usage

```bash
# Show all of today's games
bball games

# Only live or only final games
bball games --live    # or -l
bball games --final   # or -f

# Show a team's schedule (defaults to upcoming)
bball schedule knicks
bball schedule nyk --upcoming   # or -u
bball schedule nyk --recent     # or -r

# Quick recap / summary of a team's current (live) game
bball catch lakers
```

Run `bball --help` or `bball <command> --help` for all options.

---

## Examples

Output will vary based on live games. Example formatting:

```
$ bball games -l

🏀 NBA Games - 3 game(s)
────────────────────────────────────────────────────────────

🔴 LIVE - Q3 05:21 - Jan 2, 2026 at 7:30 PM EST
  BOS 72 vs MIA 68

⏰ Scheduled - Jan 2, 2026 at 10:00 PM EST
  DAL  @ LAL
```

```
$ bball schedule nyk --recent

📅 Recent Games - 5 game(s)
────────────────────────────────────────────────────────────

W Final: vs BOS - 112 vs 109

L Final: @ MIA - 98 vs 104
```

```
$ bball catch warriors

🏀 Golden State Warriors vs Los Angeles Lakers — Final
📅 Jan 2, 2026 at 7:30 PM EST

GSW 109 - LAL 125

Top Performers:
Stephen Curry (GSW) - 34 PTS, 5 REB, 7 AST
LeBron James (LAL) - 28 PTS, 8 REB, 9 AST

Last updated: just now
```

---

## Data Sources

- Scoreboard: `https://cdn.nba.com/static/json/liveData/scoreboard/todaysScoreboard_00.json`
- Season schedule: `https://cdn.nba.com/static/json/staticData/scheduleLeagueV2.json`

These are unofficial public JSON endpoints.

---

## Development

- Build: `go build .`
- Run tests: `go test ./...`
- Recommended: `go mod tidy` after cloning to fetch dependencies

### Project layout

```
main.go
```

$ bball catch warriors

🏀 Golden State Warriors vs Los Angeles Lakers — In Progress
📅 Monday, January 2, 2006 at 7:30 PM EST

GSW 72 - LAL 68

Top Performers:
Stephen Curry (GSW) - 24 PTS, 3 REB, 4 AST
LeBron James (LAL) - 18 PTS, 6 REB, 5 AST

Last updated: just now

```
**Why no Homebrew/apt releases or signed macOS binaries?**

- Hobby project; no paid Apple certs and no release automation. Use `go install` or build locally.

**Windows support?**

- Should work in a modern terminal. Colors use `fatih/color`. If colors look odd, try a different terminal.

---

## Contributing

PRs and issues welcome. Keep it simple and fast—this is just for fun.

---

## License

No license selected yet. If you plan to use or redistribute commercially, open an issue to discuss.
```
