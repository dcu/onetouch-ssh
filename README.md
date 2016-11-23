![Authy](https://raw.githubusercontent.com/dcu/onetouch-ssh/master/authy-logo.png)

# OneTouch SSH
OneTouch SSH protects a users's SSH login via a OneTouch approval.  If anyone tries to SSH with your account into a protected machine, you'll receive a OneTouch notification allowing you to Approve or Deny access.

If there is no response to the OneTouch request within a set time, your Authy registered device will instead request confirmation via an SMS delivered Authy OneCode.

Without OneTouch or SMS verification, access will not be granted.

### Prerequisites
**Install Go**
[https://golang.org/doc/install](https://golang.org/doc/install) and make sure both your GOROOT and GOPATH environmetal variables are set.

**Create an SSH key**

    Open a terminal on your local computer and enter the following:
    ssh-keygen -t rsa -C "your_email@example.com"
    Just press <Enter> to accept the default location and file name
    Enter, and re-enter, a passphrase when prompted.

This keypair should be saved in your ~/.ssh/ folder with the filename you chose.

### Install OneTouch SSH

    go get github.com/dcu/onetouch-ssh

### Configure API key

Get an Authy key from your [Authy dashboard](https://dashboard.authy.com).

The next step is to run this command to setup your environment.

    onetouch-ssh init

### Add Users

Type the following command:

    onetouch-ssh add-user <email> <country code> <phone number> <public key>

NOTE:  You can add keys in one of two ways.

**File Path**

    onetouch-ssh add-user you@your.com 1 4155551234 ~/.ssh/id_rsa.pub

**Pasted Key**

    onetouch-ssh add-user you@your.com 1 4155551234 ssh-rsa AAM8sBlW9CmrCQRFAAB3NzaC1yc2EAAHELPAADAQABAAABAQCyFQwZ2pVKfNS5iztqwaoNFaGpbLGvngQIMZgIsf+AUfGFt3c9Y4STUCKd0642miDvb6XPLINgAVPVJGzEZbZoU/+gUGGlNb+UNIVERSEFACTORY/NsE/sWqx2wuK93nvIoJXP7V+4jet9mKITt0B5aBH0mdmtY3AZS2JsksrzIcjDYldLwo+nIVFE4c4f+T7m9M8sBlW9CmrCQRF7nMbkVgSQ3Npt2IiMJaJ/1gWBxycSgMVMFiUS1Q2P3znUsBGp7p9CGssq02+NavML3sXFASyBSZ you@your.com


Next you can start adding the users using the form.
Type `Ctrl-c` to finish.

### Enable

To enable OneTouch for SSH just type:

    onetouch-ssh enable

And that's it, you can try to ssh to the server.

### Usage

When you try to connect to the ssh server it'll send you a push
notification with a limited period of time to approve:

    $ ssh ssh.server.com
    Sending approval request to your device... [sent]

If the user doesn't approve the request before the time expires an Authy OneCode  delivered via SMS is asked as a fallback.

    $ ssh ssh.server.com
    Sending approval request to your device... [sent]
    You didn't confirm the request. A text-message was sent to your phone.
    Enter security code:

### Executing Commands

When you try to run commands it'll display info about the command, the
server IP and client IP.

![OneTouch](https://raw.githubusercontent.com/dcu/onetouch-ssh/master/onetouch-ssh.png)

### Git Integration

When you try to push or fetch from git it won't display anything but
you'll receive a push notification in your phone with the info.
The information includes the server IP, client IP, geo location, branch, repository
name.

### Troubleshooting
Make sure your key and AuthyID are listed in the authorized_keys file

    cat ~/.ssh/authorized_keys

Make sure the users you want to allow access to are listed in your users.list

    cat ~/.authy-onetouch/users.list
