# Capsule Host SDK

!!! info "What's new?"
    - `v0.0.1`: 🎉 first release

## What is the Capsule Host SDK alias **Capsule HDK**?

**Capsule HDK** is a SDK to develop Golang Host applications able to run WASM modules developped with the [Capsule MDK (WASM Module SDK)](https://github.com/bots-garden/capsule-module-sdk). A **Capsule** application is a **WebAssembly Module(or Function) Runner**.

The [Capsule application**s**](https://github.com/bots-garden/capsule) **capsule-cli** and **capsule-http** are both with this Capsule HDK:
- **capsule-cli**, **CLI**. With capsule-cli, you can simply execute a **WebAssembly Capsule module** in a terminal
- **capsule-http**, an **HTTP server** that serves **WebAssembly Capsule modules**

> The Capsule Host SDK is developed in GoLang and uses the **💜 [Wazero](https://github.com/tetratelabs/wazero)** project.

!!! info "Good to know"
    - 🤗 a capsule application is **"small"** (capsule-http weighs 12M)
    - 🐳 a Capsule application is statically compiled: you can easily run it in a **Distroless** Docker container.


## What are the **added values** of a Capsule application?

A Capsule application brings superpowers to the WASM Capsule modules with **host functions**. Thanks to these **host functions**, a **WASM Capsule module** can, for example, prints a message, reads files, writes to files, makes HTTP requests, ... See the [host functions section](host-functions.md).

!!! info "Useful information for this project"
    - 🖐 Issues: [https://github.com/bots-garden/capsule-host-sdk/issues](https://github.com/bots-garden/capsule-host-sdk/issues)
    - 🚧 Milestones: [https://github.com/bots-garden/capsule-host-sdk/milestones](https://github.com/bots-garden/capsule-host-sdk/milestones)
    - 📦 Releases: [https://github.com/bots-garden/capsule-host-sdk/releases](https://github.com/bots-garden/capsule-host-sdk/releases)

