gvm
===

gvm is a Go version manager. gvm installs a Go version and prints the commands
to configure your environment to use it. gvm can only install binary versions of
Go from https://golang.org/dl/. Below are examples for common shells.

bash:

`eval "$(gvm 1.9)"`

batch (windows cmd.exe):

`FOR /f "tokens=*" %i IN ('"gvm.exe" gvm 1.9') DO %i`

powershell:

`gvm --format=powershell 1.9 | Invoke-Expression`

Or using the project's Go version as determined by the .travis.yml file. For
example:

`eval "$(gvm --project-go)"`

Installation
------------

You can download a binary release of `gvm` for your specific platform from the
[releases](https://github.com/andrewkroh/gvm/releases) page. Then just put the
binary in your `PATH` and mark it as executable (`chmod +x gvm`).

You must adjust the version and platform info in URLs accordingly.

For Bash users:

``` bash
# Linux Example (assumes ~/bin is in PATH).
curl -sL -o ~/bin/gvm https://github.com/andrewkroh/gvm/releases/download/v0.0.3/gvm-linux-amd64
chmod +x ~/bin/gvm
eval "$(gvm 1.9)"
go version
```

``` bash
# macOS Example (assumes ~/bin is in PATH).
curl -sL -o ~/bin/gvm https://github.com/andrewkroh/gvm/releases/download/v0.0.3/gvm-darwin-amd64
chmod +x ~/bin/gvm
eval "$(gvm 1.9)"
go version
```

For Windows PowerShell users:

```
Invoke-WebRequest -URI https://github.com/andrewkroh/gvm/releases/download/v0.0.3/gvm-windows-amd64.exe -Outfile C:\Windows\System32\gvm.exe
gvm --format=powershell 1.9 | Invoke-Expression
go version
```

For existing Go users:

`go get github.com/andrewkroh/gvm`
