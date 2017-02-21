gvm
===

gvm is a Go version manager. gvm installs a Go version and prints the commands
to configure your environment to use it. gvm can only install binary versions of
Go from https://golang.org/dl/. Below are examples for common shells.

bash:

`eval "$(gvm 1.7.4)"`

batch (windows cmd.exe):

`FOR /f "tokens=*" %i IN ('"gvm.exe" gvm 1.7.4') DO %i`

powershell:

`gvm --format=powershell 1.7.4 | Invoke-Expression`

Or using the project's Go version as determined by the .travis.yml file. For
example:

`eval "$(gvm --project-go)"`

Installation
------------

You can download a binary release of `gvm` for your specific platform from the
[releases](https://github.com/andrewkroh/gvm/releases) page. Then just put the
binary in your `PATH` and mark it as executable (`chmod +x gvm`).

``` bash
# Example for macOS (assume ~/bin is in PATH).
curl -sL -o ~/bin/gvm https://github.com/andrewkroh/gvm/releases/download/v0.0.1/gvm-darwin-amd64
chmod +x ~/bin/gvm
```

For existing Go users:

`go get -u github.com/andrewkroh/gvm`
