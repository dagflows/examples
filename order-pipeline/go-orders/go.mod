// The builder compiles the MODULE ROOT and runs the resulting binary to emit
// the manifest, so package main lives beside this file. Pointing `module:` at
// a subpackage in dagflows.yaml is refused.
//
// The path is never resolved by anyone - nothing imports this module - so it
// only has to be a valid module path.
module dagflows.example/order-pipeline/go-orders

go 1.27

require github.com/dagflows/sdk-go v0.5.0
