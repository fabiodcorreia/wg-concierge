# WG-Concierge
WG-Concierge allow to add new devices quickly with no need to connect to the server remotely, it's mostly and ingress tool for WireGuard networks

## Conditions
- WG-Concierge needs to be installed on the server with root access (this should change but for now needs the same access as wg)
- It should not be exposed outside the LAN
- The client needs to be on the LAN to make the registration since the email will contain and URL with internal address

## Workflow
1. The admin send an invitation by Email. The email contains an one time use URL (Basic Auth for start)
2. The server will record the URL and email
3. The user receives the email and opens the URL
4. The server checks if the URL is valid, if not reply error
5. The user gets a form to enter the device name and submits
6. Repeat step 4 for the new request
7. The server checks if the machine name already exists for that email, if so reply error
8. The server generates the private and public keys for the client (in memory)
9. The server grabs the last IP added, increment and lock it to that emails-device
10. The server generate the client configuration file (in memory)
11. The server update wg config with the new client
12. The server reply with QR Code and Configuration file
13. The server burns the URL so it can't be used anymore
14. The admin can see the table of peer clients (Basic Auth for start)

## Goals
1. Allow quick configuration of new devices without SSH to the server
2. Separation of keys, the server never store the client private key
3. Keep track of all devices registered

## Endpoints
|  Method | Auth |   Path    | Params | Body        |
|---------|------|-----------|--------|-------------|
|   GET   | Yes  | /map      |        |             |
|   GET   | Yes  | /invite   | email  |             |
|   GET   | No   | /register | token  |             |
|   POST  | No   | /register | token  | device_name |

## Out of Scope

1. Remove or Update peers (that needs to be done directly on the server)
2. Manage the wg server operations like start/stop...

## Development

Since this app needs to run on the same server of WireGuard a Vagrant Box is provided. This box start with WireGuard already installed and with a standard configuration.

The folder build will be synchronized with the box, so every time the project is compiled the result will be available inside the box.

### Start the Box
```
vagrant up
```

### Connect to the Box
```
vagrant ssh
```

### Shutdown the box
```
vagrant halt
```

### Destroy the box
```
vagrant destroy
```