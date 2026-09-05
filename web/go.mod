// This file exists only to fence the web/ npm project (including its
// node_modules, which may contain stray vendored Go source, as flatted
// does) out of the root Go module's `./...` package pattern. web/ has no
// Go code of its own — see package.json for its actual toolchain.
module braess/web-not-a-go-module

go 1.26
