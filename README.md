# Just a SMTP forwarder daemon

## What is this?

Simple SMTP forwarding daemon.

If your SMTP server (ex. sendgrid.net) requires authentication, but your old application has no capable. This simple program relay SMTP connection from the application to the server.

## Usage

Define environment variables and run.

```bash
# Define environment variables.
SMTP_BRIDGE_HOST=smtp.sendgrid.net # Real SMTP server host name
SMTP_BRIDGE_PORT=587               # Real SMTP server port number
SMTP_BRIDGE_USER=apikey            # Real SMTP server user name
SMTP_BRIDGE_PASSWORD=password      # Real SMTP server user's password

# Run program
go run main.go
# or specify binding address
go run main.go 0.0.0.0:10025
```

## Testing

```bash
netcat -C localhost 1025
HELO localhost
MAIL FROM:<foo@example.com>
RCPT TO:<bar@example.com>
DATA
From: "Foo" <foo@example.com>
Subject: test

test
.
QUIT
```
