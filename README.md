gvm
===

gvm is a Go version manager. gvm installs a Go version and prints the commands
to configure your environment to use it. gvm can install Go binary versions from
https://storage.googleapis.com/golang or build it from source. Below are
examples for common shells.

bash:

`eval "$(gvm 1.23.0)"`

cmd.exe (for batch scripts `%i` should be substituted with `%%i`):

`FOR /f "tokens=*" %i IN ('"gvm.exe" 1.23.0') DO %i`

powershell:

`gvm --format=powershell 1.23.0 | Invoke-Expression`

gvm flags can be set via environment variables by setting `GVM_<flag>`. For
example `--http-timeout` can be set via `GVM_HTTP_TIMEOUT=10m`.

Installation
------------

You can download a binary release of `gvm` for your specific platform from the
[releases](https://github.com/andrewkroh/gvm/releases) page. Then just put the
binary in your `PATH` and mark it as executable (`chmod +x gvm`).

You must adjust the version and platform info in URLs accordingly.

Linux (amd64):

``` bash
# Linux Example (assumes ~/bin is in PATH).
curl -sL -o ~/bin/gvm https://github.com/andrewkroh/gvm/releases/download/v0.5.2/gvm-linux-amd64
chmod +x ~/bin/gvm
eval "$(gvm 1.23.0)"
go version
```

Linux (arm64):

``` bash
# Linux Example (assumes ~/bin is in PATH).
curl -sL -o ~/bin/gvm https://github.com/andrewkroh/gvm/releases/download/v0.5.2/gvm-linux-arm64
chmod +x ~/bin/gvm
eval "$(gvm 1.23.0)"
go version
```

macOS (universal):

``` bash
# macOS Example
curl -sL -o /usr/local/bin/gvm https://github.com/andrewkroh/gvm/releases/download/v0.5.2/gvm-darwin-all
chmod +x /usr/local/bin/gvm
eval "$(gvm 1.23.0)"
go version
```

Windows (PowerShell):

```
[Net.ServicePointManager]::SecurityProtocol = "tls12"
Invoke-WebRequest -URI https://github.com/andrewkroh/gvm/releases/download/v0.5.2/gvm-windows-amd64.exe -Outfile C:\Windows\System32\gvm.exe
gvm --format=powershell 1.23.0 | Invoke-Expression
go version
```

Fish Shell:

Use `gvm` with fish shell by executing `gvm 1.23.0 | source` in lieu of using `eval`.

For existing Go users:

`go install github.com/andrewkroh/gvm/cmd/gvm@v0.5.2`
