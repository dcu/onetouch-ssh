# OneTouch SSH

## Install

    go get github.com/dcu/onetouch-ssh

## Configure API key

Type the following command and follow the instructions:

    onetouch-ssh init

## Add Users

Type the following command:

    onetouch-ssh add-user <email> <country code> <phone number> <public key>

Next you can start adding the users using the form.
Type `Ctrl-c` to finish.

## Enable

To enable OneTouch for SSH just type:

    onetouch-ssh enable

And that's it, you can try to ssh to the server.

## Usage

When you try to connect to the ssh server it'll send you a push
notification with a limited period of time to approve:

    $ ssh ssh.server.com
    Sending approval request to your device... [sent]

If the user doesn't approve the request before the time expires a
security code is asked as a fallback.

    $ ssh ssh.server.com
    Sending approval request to your device... [sent]
    You didn't confirm the request. A text-message was sent to your phone.
    Enter security code:

## Executing Commands 

When you try to run commands it'll display info about the command, the
server IP and client IP.

## Git Integration

When you try to push or fetch from git it won't display anything but
you'll receive a push notification in your phone with the info.
The information includes the server IP, client IP, geo location, branch, repository
name.

