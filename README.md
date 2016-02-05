# snap collector plugin - rabbitmq-collector
snap plugin for collecting metrics from a RabbitMQ node or cluster through the RabbitMQ
mamagement API. 

It's used in the [snap framework](http://github.com:intelsdi-x/snap).

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license-and-authors)
6. [Acknowledgements](#acknowledgements)


## Getting Started
### System Requirements

* RabbitMQ node or cluster
* RabbitMQ management plugin loaded on a node or a node in a cluster
* Credentials for reading from the RabbitMQ node or cluster

### Operating Systems
All operating systems currently supported by snap:
* Linux/amd64
* Darwin/amd64

### Installation
#### Download the rabbitmq-collector plugin binary
You can get the pre-built binaries for your OS and architecture at snap's [GitHub Releases](https://github.com/intelsdi-x/snap/releases) page.

#### To build the plugin binary
Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:

```
$ git clone https://github.com/intelsdi-x/snap-plugin-collector-rabbitmq.git
```

Build the plugin by running make within the cloned repo:

```
$ cd $GOPATH/src/intelsdi-x/snap-plugin-collector-rabbitmq
$ make
```

This builds the plugin in `build/rootfs/`


### Configuration and Usage
* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)
* The url, user, and password for accessing the RabbitMQ Management API to gather metrics from. The url,
user, and password can be passed via config to the snap framework on start or via the task manifest as part
of the config for the metrics to collect.

## Documentation

### Collected Metrics
This plugin as the ability to collect the following metrics:

#### Nodes

Namespace | Data Type | Description
----------|-----------|----------------------
/intel/rabbitmq/node/{name}/disk_free | int | Available disk space on the node
/intel/rabbitmq/node/{name}/disk_free_details/rate | float64 | Rate of available disk space usage on the node
/intel/rabbitmq/node/{name}/fd_used | int | Number of file descriptors in use
/intel/rabbitmq/node/{name}/fd_used_details/rate | float64 | Rate of file descriptors being used
/intel/rabbitmq/node/{name}/io_read_avg_time | float64 | Avg wall time (ms) for each disk read operation
/intel/rabbitmq/node/{name}/io_read_avg_time_details/rate | float64 | Rate of wall time for each disk read operation
/intel/rabbitmq/node/{name}/io_read_bytes | int | Total number of bytes read from disk by the persister
/intel/rabbitmq/node/{name}/io_read_bytes_details/rate | float64 | Rate of bytes read from disk by the persister
/intel/rabbitmq/node/{name}/io_read_count | int | Total number of read operations by the persister
/intel/rabbitmq/node/{name}/io_read_count_details/rate | float64 |  Rate of read operations by the persister
/intel/rabbitmq/node/{name}/io_seek_avg_time | float64 | Avg wall time (ms) for each seek operation
/intel/rabbitmq/node/{name}/io_seek_avg_time_details/rate | float64 | Rate of avg wall time (ms) for each seek operation
/intel/rabbitmq/node/{name}/io_seek_count | int | Total number of seek operations by the persister
/intel/rabbitmq/node/{name}/io_seek_count_details/rate | float64 | Rate of seek operations by the persister
/intel/rabbitmq/node/{name}/io_sync_avg_time | float64 | Avg wall time (ms) for each fsync() operation
/intel/rabbitmq/node/{name}/io_sync_avg_time_details/rate | float64 | Rate of avg wall time (ms) for each fsync() operation
/intel/rabbitmq/node/{name}/io_sync_count | int | Total number of fsync() operations by the persister
/intel/rabbitmq/node/{name}/io_sync_count_details/rate | float64 | Rate of fsync() operations by the persister
/intel/rabbitmq/node/{name}/io_write_avg_time | float64 | Avg wall time (ms) for each disk write operation
/intel/rabbitmq/node/{name}/io_write_avg_time_details/rate | float64 | Rate of avg wall time (ms) for each disk write operation
/intel/rabbitmq/node/{name}/io_write_byte | int | Total number of bytes written to disk by the persister
/intel/rabbitmq/node/{name}/io_write_byte_details/rate | float64 | Rate of bytes written to disk by the persister
/intel/rabbitmq/node/{name}/io_write_count | int | Total number of write operations by the persister
/intel/rabbitmq/node/{name}/io_write_count_details/rate | float64 |  Rate of write operations by the persister
/intel/rabbitmq/node/{name}/mem_used | int | Memory used in bytes
/intel/rabbitmq/node/{name}/mem_used_details/rate | float64 | Rate of memory used
/intel/rabbitmq/node/{name}/memory/atom | int | Total memory in use by atom
/intel/rabbitmq/node/{name}/memory/binary | int | Total binary memory in use
/intel/rabbitmq/node/{name}/memory/code | int | Total memory in use by code 
/intel/rabbitmq/node/{name}/memory/connection_channels | int | Total memory in use by connection channels 
/intel/rabbitmq/node/{name}/memory/connection_other | int | Total memory in use by other connections 
/intel/rabbitmq/node/{name}/memory/connection_readers | int | Total memory in use by connection readers 
/intel/rabbitmq/node/{name}/memory/connection_writers | int | Total memory in use by connection writers 
/intel/rabbitmq/node/{name}/memory/mgmt_db | int | Total memory in use by the Management database 
/intel/rabbitmq/node/{name}/memory/mnesia | int | Total memory in use by Mnesia database 
/intel/rabbitmq/node/{name}/memory/msg_index | int | Total memory in use by the message index 
/intel/rabbitmq/node/{name}/memory/other_ets | int | Total memory in use by other ets 
/intel/rabbitmq/node/{name}/memory/other_proc | int | Total memory in use by other processes
/intel/rabbitmq/node/{name}/memory/plugins | int | Total memory in use by plugins 
/intel/rabbitmq/node/{name}/memory/queue_procs | int | Total memory in use by queue processes
/intel/rabbitmq/node/{name}/memory/queue_slave_procs | int | Total memory in use by queue processes on slave
/intel/rabbitmq/node/{name}/memory/total | int | Total memory in use by the node 
/intel/rabbitmq/node/{name}/mnesia_disk_tx_count | int | Number of Mnesia transactions performed that required writes to disk
/intel/rabbitmq/node/{name}/mnesia_disk_tx_count_details/rate | float64 | Rate of Mnesia transactions performed that required writes to disk
/intel/rabbitmq/node/{name}/mnesia_ram_tx_count | int | Number of Mnesia transactions which have been performed that did not require writes to disk
/intel/rabbitmq/node/{name}/mnesia_ram_tx_count_details/rate | float64 | Rate of Mnesia transactions which have been performed that did not require writes to disk
/intel/rabbitmq/node/{name}/proc_used | int | Number of Erlang processes in use
/intel/rabbitmq/node/{name}/proc_used_details/rate | float64 | Rate of Erlang processes in use
/intel/rabbitmq/node/{name}/queue_index_read_count | int | Number of records read from the queue index
/intel/rabbitmq/node/{name}/queue_index_read_count_details/rate | float64 | Rate of records read from the queue index
/intel/rabbitmq/node/{name}/queue_index_write_count | int | Number of records written to the queue index
/intel/rabbitmq/node/{name}/queue_index_write_count_details/rate | float64 | Rate of records written to the queue index
/intel/rabbitmq/node/{name}/sockets_used | int | File descriptors used as sockets
/intel/rabbitmq/node/{name}/sockets_used_details/rate | float64 | Rate of File descriptors used as sockets

####  Exchanges

Namespace | Data Type | Description
----------|-----------|----------------------
/intel/rabbitmq/exchanges/{vhost}/{name}/message_stats/confirm | int | Count of messages confirmed
/intel/rabbitmq/exchanges/{vhost}/{name}/message_stats/publish_in | int | Count of messages published "in" to an exchange (not taking account of routing)
/intel/rabbitmq/exchanges/{vhost}/{name}/message_stats/publish_out | int | Count of messages published "out" of an exchange (takking account of routing)
/intel/rabbitmq/exchanges/{vhost}/{name}/message_stats/confirm_details/rate | float64 | Rate at which messages are being confirmed
/intel/rabbitmq/exchanges/{vhost}/{name}/message_stats/publish_in_details/rate | float64 | Rate of messages published "in" to an exchange
/intel/rabbitmq/exchanges/{vhost}/{name}/message_stats/publish_out_details/rate | float64 | Rate of  messages published "out" of an exchange

#### Queues

Namespace | Data Type | Description
----------|-----------|----------------------
/intel/rabbitmq/queues/{vhost}/{name}/consumers | int | Number of consumers connected to the queue
/intel/rabbitmq/queues/{vhost}/{name}/disk_reads | int | Number of disk reads made by the queue
/intel/rabbitmq/queues/{vhost}/{name}/disk_writes | int | Number of disk writes made by the queue
/intel/rabbitmq/queues/{vhost}/{name}/memory | int | Amount of memory in use by the queue
/intel/rabbitmq/queues/{vhost}/{name}/message_bytes | int | Total size (in bytes) of messages in the queue
/intel/rabbitmq/queues/{vhost}/{name}/message_bytes_persistent | int | Total size (in bytes) of messages persisted to disk
/intel/rabbitmq/queues/{vhost}/{name}/message_bytes_ram | int | Total size (in bytes) of messages in ram
/intel/rabbitmq/queues/{vhost}/{name}/message_bytes_ready | int | Total size (in bytes) of messages in the queue that are ready
/intel/rabbitmq/queues/{vhost}/{name}/message_bytes_unacknowledged | int | Total size (in bytes) of messages in the queue that are unacknowledged
/intel/rabbitmq/queues/{vhost}/{name}/message_stats/disk_reads | int | Count of disk reads for the queue
/intel/rabbitmq/queues/{vhost}/{name}/message_stats/disk_reads_details/rate | float64 | Rate of disk reads for the queue
/intel/rabbitmq/queues/{vhost}/{name}/message_stats/publish | int | Count of messages published
/intel/rabbitmq/queues/{vhost}/{name}/message_stats/publish_details/rate | float64 | Rate of messages published
/intel/rabbitmq/queues/{vhost}/{name}/messages | int | Number of messages in the queue
/intel/rabbitmq/queues/{vhost}/{name}/messages_details/rate | float64 | Rate of messages in the queue
/intel/rabbitmq/queues/{vhost}/{name}/messages_persistent | int | Number of messages persisted to disk
/intel/rabbitmq/queues/{vhost}/{name}/messages_ram | int | Number of messages in ram
/intel/rabbitmq/queues/{vhost}/{name}/messages_ready | int | Number of messages ready in the queue
/intel/rabbitmq/queues/{vhost}/{name}/messages_ready_details/rate | float64 | Rate of messages ready in the queue
/intel/rabbitmq/queues/{vhost}/{name}/messages_ready_ram | int | Number of messages in ram that are ready
/intel/rabbitmq/queues/{vhost}/{name}/messages_unacknowledged | int | Number of messages unacknowledged in the queue
/intel/rabbitmq/queues/{vhost}/{name}/messages_unacknowledged_details/rate | float64 | Rate of messages unacknowledged in the queue
/intel/rabbitmq/queues/{vhost}/{name}/messages_unacknowledged_ram | int | Number of messages in ram that are unacknowledged

#### vhosts

Namespace | Data Type | Description
----------|-----------|----------------------
/intel/rabbitmq/vhosts/{name}/messages | int | Sum of all message fields for all queues in vhost
/intel/rabbitmq/vhosts/{name}/messages_ready | int | Sum of all messages_ready fields for all queues in vhost
/intel/rabbitmq/vhosts/{name}/messages_acknowledged | int | Sum of all messages_acknowledged fields for all queues in vhost
/intel/rabbitmq/vhosts/{name}/message_stats/confirm | int | Global count of messages confirmed for vhost
/intel/rabbitmq/vhosts/{name}/message_stats/deliver_get | int | Global count of messages delivered for vhost
/intel/rabbitmq/vhosts/{name}/message_stats/get_no_ack | int | Global count of messages delivered in no-acknowledgement mode in response to basic.get across vhost
/intel/rabbitmq/vhosts/{name}/message_stats/publish | int | Global count of messages published across vhost
/intel/rabbitmq/vhosts/{name}/messages_details/rate | float64 | Rate of messages for all queues in vhost
/intel/rabbitmq/vhosts/{name}/messages_ready_details/rate | float64 | Rate of  all messages_ready fields for all queues in vhost
/intel/rabbitmq/vhosts/{name}/messages_acknowledged_details/rate | float64 | Rate of all messages_acknowledged fields for all queues in vhost
/intel/rabbitmq/vhosts/{name}/message_stats/confirm_details/rate | float64 | Global rate of messages confirmed for vhost
/intel/rabbitmq/vhosts/{name}/message_stats/deliver_get_details/rate | float64 | Global rate of messages delivered for vhost
/intel/rabbitmq/vhosts/{name}/message_stats/get_no_ack_details/rate | float64 | Global rate of messages delivered in no-acknowledgement mode in response to basic.get across vhost
/intel/rabbitmq/vhosts/{name}/message_stats/publish_details/rate | float64 | Global rate of messages published across vhost
/intel/rabbitmq/vhosts/{name}/messages_unacknowledged | int | Global count of messages unacknowledged across all queues
/intel/rabbitmq/vhosts/{name}/messages_unacknowledged_details/rate | float64 | Global rate of messages unacknowledged across all queues

### Examples
Here is an example task manifest for collecting metrics from a RabbitMQ node or cluster.

```json
{
    "version": 1,
    "schedule": {
        "type": "simple",
        "interval": "5s"
    },
    "workflow": {
        "collect": {
            "metrics": {
                "/intel/rabbitmq/nodes/*/memory/total": {},
                "/intel/rabbitmq/queues/*/*/messages": {},
                "/intel/rabbitmq/exchanges/*/*/message_stats/publish_in_details/rate": {},
                "/intel/rabbitmq/vhosts/*/messages_ready": {}
            },
            "config": {
                "/intel/rabbitmq": {
                    "url": "http://localhost:15672",
                    "user": "guest",
                    "password": "guest"
                }
            },
            "publish": [
                {
                    "plugin_name": "influx",
                    "config": {
                        "host": "localhost",
                        "port": 8086,
                        "database": "snap",
                        "user": "admin",
                        "password": "admin"
                    }
                }
            ]
        }
    }
}
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release. If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-rabbitmq/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-rabbitmq/pulls).

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support).

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [@geauxvirtual](https://github.com/geauxvirtual/)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.
