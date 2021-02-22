[![Commits.lol](/commits.lol.png?raw=true "Commits.lol")](https://commits.lol)

### Commits.lol

[![Build Status](https://travis-ci.com/TunedMystic/commits.lol.svg?branch=master)](https://travis-ci.com/TunedMystic/commits.lol)
[![codecov](https://codecov.io/gh/TunedMystic/commits.lol/branch/master/graph/badge.svg)](https://codecov.io/gh/TunedMystic/commits.lol)
[![Go Report Card](https://goreportcard.com/badge/github.com/tunedmystic/commits.lol)](https://goreportcard.com/report/github.com/tunedmystic/commits.lol)

Commits.lol is a service that collects funny commit messages from Github.

The commits are randomized on every page load, but you can also randomize commit messages with the emoji buttons!

😂 = funny commits

💩 = poop-related commits

😇 = cry-for-help commits

🤬 = _turn the censoring on/off_

<br />

### How does it work?

It uses Github's REST API to fetch commits with bits of profanity 😛.

The messages are then censored, color coded and saved to the db.

New commits are fetched from Github every hour.

<br />

### What's the tech stack?

The backend is built with Go _(Mux, Sqlx and friends)_

The frontend is plain HTML and JS. Tailwind for CSS.

The database is SQLite.
