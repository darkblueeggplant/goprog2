#!/bin/bash

set +x

export DEBIAN_FRONTEND=noninteractive


sudo apt-get update
sudo apt-get install openjdk-17-jre -y
# to avoid and reuse downloaded tar ball
pushd /vagrant
wget --continue --show-progress --https-only --timestamping https://dlcdn.apache.org/kafka/3.9.1/kafka_2.12-3.9.1.tgz
popd

mv /vagrant/kafka_2.12-3.9.1.tgz .
tar -xvf kafka_2.12-3.9.1.tgz -C /home/vagrant
rm -rf kafka_2.12-3.9.1.tgz
mv /home/vagrant/kafka{_2.12-3.9.1,}

echo "PATH=$PATH:/home/vagrant/kafka/bin" > /home/vagrant/.bashrc
if [ $HOSTNAME == "kafka" ]
then
	KAFKA_CLUSTER_ID=$(/home/vagrant/kafka/bin/kafka-storage.sh random-uuid)
	echo $KAFKA_CLUSTER_ID > /vagrant/uuid
else
	KAFKA_CLUSTER_ID=$(cat /vagrant/uuid)
fi
echo $KAFKA_CLUSTER_ID

/home/vagrant/kafka/bin/kafka-storage.sh format --standalone -t $KAFKA_CLUSTER_ID -c /home/vagrant/kafka/config/kraft/server.properties

sudo chown -R vagrant:vagrant /home/vagrant/kafka

nohup /home/vagrant/kafka/bin/kafka-server-start.sh /home/vagrant/kafka/config/kraft/server.properties > kafka.log 2>&1 &


