$targets = @(
    @{ OS="windows"; ARCH="amd64"; OUT="windows/starter-kit.exe" }
    @{ OS="linux"; ARCH="amd64"; OUT="linux/starter-kit" }
    @{ OS="linux"; ARCH="arm64"; OUT="linux/starter-kit-arm64" }
    @{ OS="darwin"; ARCH="amd64"; OUT="darwin/starter-kit" }
    @{ OS="darwin"; ARCH="arm64"; OUT="darwin/starter-kit-arm64" }
)

foreach ($target in $targets) {
    $env:GOOS = $target.OS
    $env:GOARCH = $target.ARCH
    Write-Output "Building for $($target.OS) $($target.ARCH)..."
    go build -o "build/$($target.OUT)"
}

Write-Output "Build selesai! Cek folder 'build/'."
