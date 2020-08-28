# Proofable Image

![Proofable Image Screenshot](docs/proofable-image-screenshot.png)

ProofableImage builds trust into your image by creating a blockchain certificate for it. The image certificate can not only prove the image as a whole but also prove the pixel boxes and the metadata inside it. For more details, please read through [this Medium post]().

If you want to prove your file system, please try out the [Proofable CLI](https://docs.proofable.io/cmd/proofable-cli/).

If you want to build trust into your own application, please check out the [Proofable Framework](https://proofable.io/).

## Installation

### Download a prebuilt binary

Following these steps to install the latest prebuilt binary into your current working directory, which is recommended.

#### For macOS and Linux users

Copy, paste and run the following bash command in a [macOS Terminal](https://support.apple.com/en-au/guide/terminal/welcome/mac):

```zsh
bash -c "$(eval "$(if [[ $(command -v curl) ]]; then echo "curl -fsSL"; else echo "wget -qO-"; fi) https://raw.githubusercontent.com/SouthbankSoftware/proofable-image/master/install.sh")"
```

#### For Windows users

Copy, paste and run the following PowerShell command in a [PowerShell prompt](https://docs.microsoft.com/en-us/powershell/scripting/overview?view=powershell-7):

```zsh
& ([ScriptBlock]::Create((New-Object Net.WebClient).DownloadString('https://raw.githubusercontent.com/SouthbankSoftware/proofable-image/master/install.ps1')))
```

### Build your own binary

Install a global binary using `go get`:

```zsh
GO111MODULE=on go get github.com/SouthbankSoftware/proofable-image
```

Or clone this repo and build one:

```zsh
git clone https://github.com/SouthbankSoftware/proofable-image.git
cd proofable-image
make
```

## Usage

```zsh
./proofable-image path/to/your/image.png
```

This will create an image certificate at `path/to/your/image.png.imgcert` if it doesn't exist yet, and verify the image against it. Then an image viewer will pop up to show any tampering. You can use the option `-imgcert-path` to test the certificate on another image:

```zsh
./proofable-image -imgcert-path=path/to/your/image.png.imgcert path/to/another/image.png
```

For all available options, please use:

```zsh
./proofable-image -h
```
