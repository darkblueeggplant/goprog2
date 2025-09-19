#!/bin/bash

set +x

apt-get update
apt-get install openjdk-17-jre -y
# to avoid and reuse downloaded tar ball
pushd /vagrant
wget --continue --show-progress --https-only --timestamping https://dlcdn.apache.org/kafka/3.9.1/kafka_2.12-3.9.1.tgz
popd

mv /vagrant/kafka_2.12-3.9.1.tgz .
tar -xvf kafka_2.12-3.9.1.tgz -C /home/vagrant
rm -rf kafka_2.12-3.9.1.tgz
mv /home/vagrant/kafka{_2.12-3.9.1,}
chown -R vagrant:vagrant /home/vagrant/kafka