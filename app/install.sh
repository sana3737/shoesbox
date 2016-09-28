#!/bin/bash
# build 
GOOS=linux GOARCH=amd64 go build footWear.go

# scp する処理
expect -c "
set timeout -1
spawn scp -r css dbconfig.yml footWear images init.sh templates migration nginx.conf $1@$2:.
expect \"Are you sure you want to continue connecting (yes/no)?\" {
    send \"yes\n\"
    expect \"$1@$2's password:\"
    send \"$3\n\"
} \"$1@$2's password:\" {
    send \"$3\n\"
} \"Permission denied (publickey,gssapi-keyex,gssapi-with-mic).\" {
    exit
}
interact
"

# ssh する処理
expect -c "
set timeout -1
spawn ssh $1@$2 ./init.sh&
expect \"Are you sure you want to continue connecting (yes/no)?\" {
    send \"yes\n\"
    expect \"$1@$2's password:\"
    send \"$3\n\"
} \"$1@$2's password:\" {
    send \"$3\n\"
} \"Permission denied (publickey,gssapi-keyex,gssapi-with-mic).\" {
    exit
}
interact
"
