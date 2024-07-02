# Dockerfile
FROM alpine:latest

# Install openssl
RUN apk add --no-cache openssl

# Create a directory to store the certificates
RUN mkdir -p /certs

# Copy the script into the image
COPY ./script.sh /script.sh

# Make the script executable
RUN chmod +x /script.sh

# Set the entrypoint to the script
ENTRYPOINT ["/script.sh"]