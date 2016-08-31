# Ransomware

> Note: This project is purely academic, use at your own risk. I do not encourage in any way the use of this software illegally or to attack targets without their previous authorization

I'm not a security analyst so do not expect something perfect, but any help will be apreciated. This project is composed of two main parts for now, the server and the malware itself

The server is responsible for the keypair generation, not RSA,but an identification key and an encryption key that is kept safe on the server, granting a way to recover the encryption key in the future (I plan to change this system to a pair of RSA-2048 keys)

### Installation
```
go get -v github.com/mauri870/ransomware
git clone https://github.com/mauri870/ransomware.git
cd ransomware
```

### Building the binaries
> DON'T RUN ANY BINARY IN YOUR PERSONAL MACHINE, EXECUTE ONLY IN A TEST ENVIRONMENT!

We need build the server and the malware as follows, generating binaries `only` for windows:
```
make
```
After that, a binary called `ransomware.exe` and a `server.exe` will be generated on root

> DON'T RUN ANY BINARY IN YOUR PERSONAL MACHINE, EXECUTE ONLY IN A TEST ENVIRONMENT!

## Usage and How it Works
Feel free to edit the parameters across the files for testing.
Put the binaries on a correct windows test environment, start the server by double click or run then on the terminal.
It will wait for the malware contact and generate/persist the id/encryption keypairs

When double click on `ransomware.exe` binary it will walk interesting directories and encrypting all files that match the interesting file extensions, recreating then with encrypted content and a custom extension(.encrypted by default) and create a READ_TO_DECRYPT.txt file on desktop

In theory, for decrypt your files you need send an amount of BTC to the attacker's wallet, followed by a contact sending your ID(located on the file created on desktop). If your payment was confirmed, the attacker possibly will return your encryption key and you can use then to recover your files. This exchange can be accomplished in several ways(Possibly use an RSA algorithm will change this order).

Let's suppose you get your encryption key back (for testing it is on the file on desktop) you can use then on a terminal:
```
ransomware.exe decrypt yourencryptionkeyhere
```
And that's it, got your files back :smile:

As you can see, building a functional ransomware is not dificult, anyone with programming skills can build that
