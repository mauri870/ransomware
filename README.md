# Ransomware

> Note: This project is purely academic, use at your own risk. I do not encourage in any way the use of this software illegally or to attack targets without their previous authorization

**Do not criticize me because I'm sharing a malware, the intent here is to disseminate and teach more about security in the actual world. Remember, security is always a double-edged sword**

### What is Ransomware?
Ransomware is a type of malware that prevents or limits users from accessing their system, either by locking the system's screen or by locking the users' files unless a ransom is paid. More modern ransomware families, collectively categorized as crypto-ransomware, encrypt certain file types on infected systems and forces users to pay the ransom through certain online payment methods to get a decrypt key.

### Project Summary
This project aims to build a "non hidden" crypto-ransomware written in Go. Basically, it will encrypt your files using AES-256-CFB, a strong encryption algorithm, using RSA-2048 to secure the key exchange with server. Yeah, a Cryptolocker like malware.

It is composed of two main parts, the server and the malware itself.

The server is responsible for store the Id and the respective encryption key, received from the malware binary during execution.

The malware encrypt with your RSA-2048 public key a payload containing the id/enckey generated on runtime, sending then to the server, where it is properly decrypted with the respective RSA private key, and then persisted for future usage.

### Installation
```
git clone https://github.com/mauri870/ransomware.git
go get -v github.com/akavel/rsrc
cd ransomware
go get -d -v
```

### Building the binaries

> DON'T RUN ransomware.exe IN YOUR PERSONAL MACHINE, EXECUTE ONLY IN A TEST ENVIRONMENT!

#### Autobuild

```
make
```
If you like build the server for windows from a unix machine, run `env GOOS=windows make` 

#### Manually
Generate a new RSA-2048 keypair:
```
openssl genrsa -out private.pem 2048
openssl rsa -in private.pem -outform PEM -pubout -out public.pem
```
> Keys protected with password is not supported

After that, a `private.pem` and `public.pem` will be created.
Copy the content of private.pem to `PRIV_KEY` on `server/main.go` and the public.pem content to `PUB_KEY` on `encrypt.go`.
Remember to put the content inside the []byte() conversion

We need build the server and the malware as follows, generating the binaries, the malware only for windows:
```
make simple-build
```
If you like build the server for windows from a unix machine, run `env GOOS=windows make simple-build`

After that, a binary called `ransomware.exe` and a `server`/`server.exe` will be generated on the build folder. The execution of `ransomware.exe` (even if it is compiled for linux/darwin) is locked to windows machines only.

By default, the server will listen on `localhost:8080`

> DON'T RUN ransomware.exe IN YOUR PERSONAL MACHINE, EXECUTE ONLY IN A TEST ENVIRONMENT!

## Usage and How it Works
Feel free to edit the parameters across the files for testing.
Put the binaries on a correct windows test environment, start the server by double click or run then on the terminal.
It will wait for the malware contact and persist the id/encryption keys

When double click on `ransomware.exe` binary it will walk interesting directories and encrypting all files that match the interesting file extensions using AES-256-CFB, recreating then with encrypted content and a custom extension(.encrypted by default) and create a READ_TO_DECRYPT.html file on desktop

In theory, for decrypt your files you need send an amount of BTC to the attacker's wallet, followed by a contact sending your ID(located on the file created on desktop). If your payment was confirmed, the attacker possibly will return your encryption key and you can use then to recover your files. This exchange can be accomplished in several ways.

Let's suppose you get your encryption key back, you can retrieve it pointing to the following url:

```
curl http://localhost:8080/api/keys/:id
```
Where `:id` is your identification stored on the file on desktop. After, run on a terminal:

```
ransomware.exe decrypt yourencryptionkeyhere
```
And that's it, got your files back :smile:

## Server endpoints

The server has only two endpoints

`POST api/keys/add` - Used by the malware to persist new keys. Some verifications are made, like the verification of the RSA autenticity. Returns 204 (empty content) in case of success or a json error.

`GET api/keys/:id` - Id is a 32 characters parameter, representing an Id already persisted. Returns a json containing the encryption key or a json error

As you can see, building a functional ransomware, with some of the best existing algorithms is not dificult, anyone with programming and security skills can build that.
