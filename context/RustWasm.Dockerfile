FROM rust:1.31

RUN rustup update
RUN rustup target add wasm32-unknown-unknown