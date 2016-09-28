#!/bin/bash
DBNAME=treasure
ENV=development

# git install
which git || sudo yum install -y git

# mysql install
which mysql
if [ $? != 0 ]; then
    sudo yum install -y mysql
    sudo yum install -y mysql-server
fi
mysql -u root -h localhost --protocol tcp -e "drop database if exists treasure"
mysql -u root -h localhost --protocol tcp -e "create database treasure"
mysql -u root -h localhost --protocol tcp treasure < ./migration/sql

# nginx
sudo mv nginx.conf /etc/nginx/.
sudo service nginx restart

# go run
./footWear