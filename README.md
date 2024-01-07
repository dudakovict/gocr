# GOCR - Golang OCR Server

GOCR is a simple Golang-based OCR (Optical Character Recognition) server that allows you to upload images and extract text from them.

## Configuration

1. **Generate TLS Certificates:**
   - Before running the application, generate TLS certificates using OpenSSL.
     ```
     openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem
     ```
     Follow the prompts to provide the necessary information.

2. **Configure the Server:**
   - Customize the server configuration by modifying the `conf.go` file or set the necessary environment variables.
     Adjust parameters such as port, timeouts, and maximum file size according to your requirements.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Gosseract](https://github.com/otiai10/gosseract) - OCR library for Golang.
