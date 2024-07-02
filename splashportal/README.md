# Splash Web Service

This is a PHP file that serves as the splash web service for a web application. The web service receive the phone number from the client and sends it to the account manager to authenticate and validate the user phone number. 

## Functionality

The splash web service performs the following tasks:

1. Get the phone number from the client.
2. Check if the user has accepted the terms and conditions.
3. Send the phone number to the account manager.

## Dependencies

The splash web service depends on the following service:

- Account Manager Service: A service that manages user accounts and sends messages via Webex Connect.

## Endpoints

The splash web service has the following endpoints:

- GET /: Renders the Access request form.


## Environment Variables

The splash web service uses the following environment variables:

- ACCOUNTMANAGER: The URL of the account manager service. Defaults to http://ACCOUNTMANAGER.

## Scenario

Upon connecting to the WiFi network, users are prompted to sign in. They must enter their phone number and agree to the terms and conditions. The service then triggers the Account Manager to manage account creation and triggering Webex Connect.
