[package]
name = "run_7"
edition = "v0.9.0"
version = "0.0.1"

[dependencies]
helloworld = "0.1.0"

[profile]
entries = ["main.k", "${helloworld:KCL_MOD}/main.k"]
