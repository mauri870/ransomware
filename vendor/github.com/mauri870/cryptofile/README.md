# Cryptofile

Encrypt or decrypt files using AES-256-CFB, AES-192-CFB or AES-128-CFB

> Caution when using with the -delete flag, if you miss the key used for encrypt your file probably you never decrypt it

## Installation
```
 go get -v github.com/mauri870/cryptofile
 go install github.com/mauri870/cryptofile
```

# Usage
Use `cryptofile -h` for a complete list of flags

> For AES-256, use a 32 digits key, for AES-192 and AES-128, use 24 or 16 digits respectivelly

```
echo Hello World! > filetoencrypt.txt
cryptofile -key yourkeyhere -in filetoencrypt.txt
cryptofile -decrypt -key yourkeyhere -in filetoencrypt.txt.encrypted
```
You can also use `-delete` to remove the input file
