## Build

### Local Build

```sh
go build -o gemini-invoice
```

### Release Build

```sh
./build.sh
```

## Signing

### Signing macOS Applications with a Self-Signed Certificate

#### 1. Create a Self-Signed Certificate

1. **Open Keychain Access**:
   Open Keychain Access from Applications > Utilities.

2. **Create a Certificate**:

   -  In the Keychain Access menu, go to `Keychain Access > Certificate Assistant > Create a Certificate`.
   -  Fill in the "Certificate Information":
      -  Name: Your certificate name.
      -  Identity Type: `Self Signed Root`.
      -  Certificate Type: `Code Signing`.

3. **Customize Settings**:

   -  Click "Continue" and customize settings as needed, then click "Create".

4. **Save the Certificate**:
   -  The new certificate will be saved in Keychain Access.

#### 2. Sign the Application Bundle

1. **Open Terminal**.
2. **Use `codesign` to Sign the Application**:

   ```sh
   codesign --deep --force --verify --verbose --sign "Your Certificate Name" path/to/your/app.app
   ```

3. **Verify the Signature**:

   ```sh
   spctl --assess --type execute --verbose path/to/your/app.app
   ```

### Signing Windows Applications on macOS

For Windows applications, you will need a tool like `osslsigncode` and OpenSSL.

#### 1. Install `osslsigncode`

1. **Install Homebrew** (if you don’t have it):

   ```sh
   /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
   ```

2. **Install `osslsigncode`**:

   ```sh
   brew install osslsigncode
   ```

#### 2. Create a Self-Signed Certificate for Windows

1. **Generate a Private Key and Certificate**:

   ```sh
   openssl req -new -newkey rsa:2048 -days 365 -nodes -x509 -keyout myprivatekey.key -out mycertificate.crt
   ```

2. **Convert the Certificate to PKCS12 Format**:

   ```sh
   openssl pkcs12 -export -inkey myprivatekey.key -in mycertificate.crt -out mycertificate.p12
   ```

#### 3. Sign the Windows Executable

1. **Use `osslsigncode` to Sign**:

   ```sh
   osslsigncode sign -pkcs12 mycertificate.p12 -pass privatekeypass -n "My Application" -i http://www.example.com -in myapp.exe -out myapp-signed.exe
   ```

   Breaking down the command:

   -  `-pkcs12 mycertificate.p12`: The certificate file.
   -  `-pass privatekeypass`: The password for the private key (if any).
   -  `-n "My Application"`: The name of the application.
   -  `-i http://www.example.com`: Your website or company URL.
   -  `-in myapp.exe`: The input executable file.
   -  `-out myapp-signed.exe`: The output signed executable file.

### Summary

-  **macOS**: Use Keychain Access to create a self-signed certificate, then use `codesign` and `spctl` to sign and verify your app.
-  **Windows**: Use OpenSSL to create a self-signed certificate, convert it to PKCS12 format, and use `osslsigncode` to sign the executable.
