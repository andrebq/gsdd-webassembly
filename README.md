# What is this?

A small experiment written as part of GSDD days at N26. Basically it is just a convoluted way to use Go+Webassembly+Rust. Don't put blame on me if this eats up your browser memory or if it burns your server (should not be the case, but anyways....)

# How to use?

First of all you should have our good old friend `make` installed in your computer, then you can:

* Call `make go-install-deps` to install some utility tools
* Call `make` to watch your files and recompile if they change
* Call `make run-hchan` to run the `hchan` binary (exposes `Go channels` over `HTTP POST operations`)
* Call `make run-serve` to run the `wfsd web server` to serve your static files from dist
* Call `make run-runner` to compile the `webassembly runner` and a simple `rust/webassembly` library to test it. This runner is almost an exact copy of `perlin-network life runner`

Each `make run-???` operation will block while running, so you might need more than one terminal to run the entire thing.

# What it does?

It uses `virtual dom diffing` to update an HTML page with some `ping/pong` messages exchanged between `goroutines` running in the browser (executed from an `go-webassembly` module) and the `rust demo code`.

# Docs?

Not so much, this was just a simple experiment done in a couple of hours.
