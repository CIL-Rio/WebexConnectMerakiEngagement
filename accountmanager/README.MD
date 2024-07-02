# Account Manager

The Account Manager Service is a Go-based HTTP server designed to manage account creation and password management within an LDAP directory for external Meraki splash pages. It serves as an intermediary, processing requests from the web service to create guest accounts and subsequently authenticate user identity through Webex Connect.


## Team

Cisco Innovation Labs

## Features



- **User Account Creation**: Automatically creates a new user account in the LDAP directory if the user does not exist.
- **Password Management**: Generates a new password for account creation and updates the password for existing users.
- **Webex Connect Integration**: Sends user data to Webex Connect for authentication purposes.
- **REST API Endpoint**: Exposes an API endpoint (`/api/user`) for user creation or password update requests via JSON payloads.

## Prerequisites

Before deploying the Account Manager Service, ensure you have the following:

- An LDAP server accessible via the provided LDAP URL.
- Credentials for an LDAP admin user with permissions to add and modify user entries.
- A Webex Connect webhook URL and key for sending user data.
- Go (Version 1.21 or later recommended)

## Configuration

The service requires several environment variables to be set for its operation:

- `LDAP_URL`: The URL of the LDAP server.
- `LDAP_ROOT`: The root DN (Distinguished Name) for LDAP operations.
- `LDAP_ADMIN_USERNAME`: The username of the LDAP admin.
- `LDAP_ADMIN_PASSWORD`: The password for the LDAP admin.
- `LDAP_USERS_OU`: The organizational unit (OU) under which user accounts will be created or updated.
- `WEBEXCONNECT_HOOK_URL`: The Webex Connect webhook URL.
- `WEBEXCONNECT_KEY`: The key for authenticating to the Webex Connect webhook.

## Setup

1. Install Go, set your GOPATH, and make sure it's in your PATH.
2. Clone the repository:

```sh
git clone [repository-url]
cd [repository-directory]
```

Get the go-ldap/ldap module (if not already done):

```sh
go get github.com/go-ldap/ldap/v3
go get github.com/sirupsen/logrus
```

Configure the LDAP server connection details and admin credentials in main.go.

## Running the Server

To start the server, run the following command from the root directory of the project:

```sh
go run main.go
```

The server will start on port 8080 and listen for API requests.


# API Usage

To create a new user or update an existing user's password, send a POST request to `/api/user` with a JSON payload containing the user's details. The service will handle LDAP operations and communicate with Webex Connect as necessary.

Example POST request using curl:


```sh
curl -X POST http://localhost:8080/api/user \
     -H "Content-Type: application/json" \
     -d '{
           "phoneNumber": "1234567890",
           "loginURL": "https://example.com/login",
           "continueURL": "https://example.com/continue",
           "apMAC": "00:11:22:33:44:55",
           "apName": "AP1",
           "apTags": ["tag1", "tag2"],
           "clientMAC": "66:77:88:99:AA:BB",
           "clientIP": "192.168.1.100"
         }'

```
The server will respond with a JSON payload.

```
{
  "success": true,
  "message": "Operation completed successfully"
}
```