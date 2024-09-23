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
# Gemini Invoice

Gemini Invoice is a simple invoice generation application built with Go and Fyne.

## Installation Instructions

### macOS

#### Using Pre-built Application
1. Download the latest release for macOS from the [Releases](https://github.com/yourusername/gemini-invoice/releases) page.
2. Unzip the downloaded file.
3. Drag the "Gemini Invoice.app" to your Applications folder.
4. Double-click the app to run it.

If you encounter a security warning when trying to open the app, follow these steps:
1. Right-click (or Control-click) on the app icon.
2. Select "Open" from the context menu.
3. Click "Open" in the dialog box that appears.

#### Building from Source
1. Ensure you have Go installed on your system. If not, download and install it from [golang.org](https://golang.org/).
2. Clone the repository:
   ```
   git clone https://github.com/yourusername/gemini-invoice.git
   cd gemini-invoice
   ```
3. Run the build script:
   ```
   ./build.sh
   ```
4. The built application will be available in the `dist` directory as "Gemini Invoice.app".

### Windows

#### Using Pre-built Application
1. Download the latest release for Windows from the [Releases](https://github.com/yourusername/gemini-invoice/releases) page.
2. Unzip the downloaded file to a location of your choice.
3. Navigate to the extracted "Gemini Invoice Windows" folder.
4. Double-click the "Gemini Invoice.bat" file to run the application.

#### Building from Source
1. Ensure you have Go installed on your system. If not, download and install it from [golang.org](https://golang.org/).
2. Install MinGW-w64 for CGo support. You can download it from [here](https://sourceforge.net/projects/mingw-w64/).
3. Clone the repository:
   ```
   git clone https://github.com/yourusername/gemini-invoice.git
   cd gemini-invoice
   ```
4. Run the build script:
   ```
   ./build.sh
   ```
5. The built application will be available in the `dist/Gemini Invoice Windows` directory.

## Usage

1. Launch the Gemini Invoice application.
2. Fill in the invoice details in the application:
   - ID: Automatically generated, but can be modified
   - To: Enter the customer's name or company
   - Item Name: Enter the name of the product or service
   - Item Price: Enter the price for the item
3. Click the "Select Output Directory" button to choose where to save the generated invoice.
4. Click "Generate Invoice" to create a PDF invoice in the selected directory.

## Configuration

You can customize some default values by editing the `config.yaml` file before building or running the application. The config file allows you to set:

- Title: The default title for your invoices
- From: Your company or personal details
- Logo: The path to your company logo image

## Development

To contribute to the development of Gemini Invoice:

1. Fork the repository on GitHub.
2. Clone your forked repository locally.
3. Make your changes and test them thoroughly.
4. Push your changes to your fork.
5. Submit a pull request with a clear description of your changes.

## License

[Specify your license here, e.g., MIT, GPL, etc.]

## Support

For support, please open an issue on the GitHub repository.
