#!/bin/bash

set +x

cat > /home/vagrant/kafka/config/kraft/server.properties << 'EOF'
# Basic configuration
process.roles=broker,controller
node.id=1
controller.quorum.voters=1@192.168.1.151:9094

# Listeners
listeners=PLAINTEXT://0.0.0.0:9093,CONTROLLER://0.0.0.0:9094
advertised.listeners=PLAINTEXT://192.168.1.151:9093
inter.broker.listener.name=PLAINTEXT
controller.listener.names=CONTROLLER

# Security protocol mapping
listener.security.protocol.map=PLAINTEXT:PLAINTEXT,CONTROLLER:PLAINTEXT

# Log settings
log.dirs=/home/vagrant/kafka-logs
num.partitions=1
num.recovery.threads.per.data.dir=1

# Replication settings
offsets.topic.replication.factor=1
transaction.state.log.replication.factor=1
transaction.state.log.min.isr=1

# Retention settings
log.retention.hours=168
log.retention.check.interval.ms=300000

# Group settings
group.initial.rebalance.delay.ms=0

# Security settings
sasl.enabled.mechanisms=
ssl.client.auth=none
ssl.enabled.protocols=
ssl.keystore.type=
ssl.protocol=
ssl.truststore.type=

# Topic settings
auto.create.topics.enable=true
delete.topic.enable=true

# Network settings
socket.request.max.bytes=1073741824
num.network.threads=3
num.io.threads=8
socket.send.buffer.bytes=102400
socket.receive.buffer.bytes=102400

# Timeout settings
request.timeout.ms=30000
connections.max.idle.ms=60000
EOF

export KAFKA_SOCKET_REQUEST_MAX_BYTES=1073741824

KAFKA_CLUSTER_ID=$(/home/vagrant/kafka/bin/kafka-storage.sh random-uuid)
echo "Cluster ID: $KAFKA_CLUSTER_ID"

/home/vagrant/kafka/bin/kafka-storage.sh format --standalone -t $KAFKA_CLUSTER_ID -c /home/vagrant/kafka/config/kraft/server.properties

echo "Starting Kafka server on port 9093..."
nohup /home/vagrant/kafka/bin/kafka-server-start.sh /home/vagrant/kafka/config/kraft/server.properties > /home/vagrant/kafka.log 2>&1 &

sleep 5
echo "Checking if Kafka is running..."
ss -tlnp | grep 9093 || echo "Kafka may not be running yet, check logs: /home/vagrant/kafka.log"

echo "Kafka startup initiated. Check logs at /home/vagrant/kafka.log"
