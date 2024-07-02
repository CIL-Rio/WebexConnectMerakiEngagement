#!/bin/sh

# Create a directory to store the certificates
mkdir -p ./certs

# Generate a new private key
openssl genrsa -out ./certs/merakiauth.key 2048

# Create a new Certificate Signing Request (CSR)
openssl req -new -key ./certs/merakiauth.key -out ./certs/merakiauth.csr -subj "${CERT_SUBJECT}"

# Generate the self-signed certificate
openssl x509 -req -days 365 -in ./certs/merakiauth.csr -signkey ./certs/merakiauth.key -out ./certs/merakiauth.crt

# Change the certificate ownership to the user
chmod 664 ./certs/*

#print that the script is complete
echo "Certificate generation complete"