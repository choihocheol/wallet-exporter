# wallet-exporter

`wallet-exporter` is a Prometheus exporter that queries blockchain nodes and exposes wallet balance metrics for Prometheus to scrape.

## About

This project is forked from [`QuokkaStake/cosmos-wallets-exporter`](https://github.com/QuokkaStake/cosmos-wallets-exporter) v4.0.1.

In addition to the original Cosmos SDK support, this fork adds:

- Support for scraping wallet balance metrics from EVM-based blockchain nodes.
- Compatibility with both Cosmos SDK and EVM-based networks through a single exporter.
