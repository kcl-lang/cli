[package]
name = "test_space"
edition = "v0.10.0"
version = "0.0.1"

[dependencies]
newhelloworld = { oci = "oci://ghcr.io/kcl-lang/helloworld", tag = "0.1.4", package = "subhelloworld", version = "0.0.1" }
