# timeruler
A time blocking service that can be run locally or remotely as a server

# Motivation
After some experimentation with multiple time management strategies, I found [timeblocking](https://todoist.com/productivity-methods/time-blocking) to be the most effective. Since I'm a software person who works on a computer, I want to avoid switching to a notebook to write, check, and update my time blocking plan for the day. So I created timeruler, with which I can build a time blocking schedule, update it throughout the day as short-notice tasks come up, and get notified through [ntfy](https://ntfy.sh) so that I can go as deep into my work as possible without constantly looking at the clock to check if I should be switching to a new task. I'm continuously adding new features (and UI improvements through [trctl](https://github.com/dethancosta/trctl)) as I discover pain points and potential improvements.

# Installation
Requires Go >= 1.21.<br>
Run `git clone https://github.com/dethancosta/timeruler`. From the cloned directory, run `go install -o timeruler ./cmd/tr-server`.<br>
Alternatively, if you don't want to download the source code, you can run `go install github.com/dethancosta/timeruler/cmd/tr-server@latest`.

# Usage/API Reference
API documentation coming soon. If using [trctl](https://github.com/dethancosta/trctl), run `trctl` for available commands

# Tests
Tests have been written for most of the internal funcationality. I still need to write tests for the API itself. To run the existing tests, run `go test ./...` from the `timeruler` directory.