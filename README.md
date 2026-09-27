# Compliance Monitor

A backend service written in Go that checks financial transactions for suspicious patterns in real time.

**-- This is a simulation not a certified compliance product, and it doesn't use real sanctions or watchlist data.**

**Status: running locally with a live browser dashboard, demo seeding, and alert resolution flows. Not yet deployed. The service still uses in-memory storage, so data resets when the server restarts.**

## What problem this is based on

Banks and financial companies are required to watch transactions for money laundering and fraud, and report anything suspicious. This project is a simplified version of that kind of system. It doesn't handle real money or connect to any real bank, but it checks synthetic transactions using the same kinds of patterns real systems look for.

## What it does

You send it a transaction (a user ID, an amount, and who they paid). It checks that transaction against three rules:

- **Structuring** — flags someone making several transactions just under a reporting limit, which is a known way people try to avoid scrutiny on one big transfer.
- **Velocity anomaly** — flags a transaction that's way bigger than a user's normal spending, using basic statistics (standard deviation) to spot the outlier.
- **Watchlist** — flags a transaction if the person being paid is on a blocklist.

Every check gets logged, whether it flags something or not, so you can always see why a transaction was or wasn't flagged.

## Dashboard

The project now includes a lightweight browser dashboard for operational visibility.

- Open the app root endpoint in a browser: `http://localhost:8080/`
- View the live alert list as transactions are evaluated
- Seed demo data to trigger structuring, velocity, and watchlist scenarios
- Hide resolved alerts, refresh the feed, and resolve alerts directly from the UI

This is a simple monitoring dashboard rather than a full analytics suite, but it gives a usable operational view while the service is running locally.

## How it's built

- Transactions come in through an HTTP API and get put in a queue instead of being checked immediately, so the API responds fast.
- A set of background workers pulls transactions off the queue and runs the checks. Multiple transactions get processed at the same time, but transactions from the same user are always processed in the order they came in.
- Data is currently stored in memory (a Go map, for now).

## Run it locally

```bash
go build -o server ./cmd/server
./server
# runs on :8080, or set the PORT environment variable
