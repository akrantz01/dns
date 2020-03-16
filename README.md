# DNS
A custom DNS server that can dynamically add and remove records through a web UI backed by a REST API. 
It uses [Bolt](https://github.com/etcd-io/bbolt) for storage and [Miekg DNS](https://github.com/miekg/dns) for, well, DNS.

## Configuration
All configuration is done through a YAML file.
An example configuration file with descriptions of each field can be found at [`config.sample.yaml`](/config.sample.yaml).

## Deployment
The server can be deployed via either Docker or a standalone binary.
All assets are bundled with the binary, so all you need to do is compile it.
The Docker image is on [Docker Hub](https://hub.docker.com/r/akrantz/dns) and the binary can be download from the [releases](https://github.com/akrantz01/krantz.dev/releases) page.
The server looks for a configuration file named `config.yaml` in either the user's home directory or the working directory.
To pass the configuration file to the Docker container run it with the argument: `-v /path/to/config.yaml:/config.yaml:ro`.
