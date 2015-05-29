# BOSH OpenStack CPI [![Build Status](https://travis-ci.org/frodenas/bosh-openstack-cpi.png)](https://travis-ci.org/frodenas/bosh-openstack-cpi)

This is an **experimental** external [BOSH CPI](http://bosh.io/docs/bosh-components.html#cpi) for [OpenStack](https://www.openstack.org/).

## Disclaimer

This is **NOT** presently a production ready CPI. This is a work in progress. It is suitable for experimentation and may not become supported in the future.

## Usage

### Deployment
This CPI can be deployed using the [BOSH OpenStack CPI release](https://github.com/frodenas/bosh-openstack-cpi-boshrelease).

### Installation

Using the standard go get:

```
$ go get github.com/frodenas/bosh-openstack-cpi/main
```

### Configuration

Create a configuration file:

```
{
  "openstack": {
    "identity_endpoint": "http://identity-endpoint.example.com/v2.0",
    "username": "username",
    "password": "password",
    "tenant_name": "tenant",
    "region": "region",
    "default_keypair": "keyname",
    "default_security_groups": [
      "security_group"
    ],
    "disable_config_drive": false,
    "disable_neutron": false,
    "ignore_server_availability_zone": false
  },
  "actions": {
    "agent": {
      "mbus": "https://mbus:mbus@0.0.0.0:6868",
      "ntp": [
        "0.north-america.pool.ntp.org"
      ],
      "blobstore": {
        "type": "local",
        "options": {}
      }
    },
    "registry": {
      "protocol": "http",
      "host": "127.0.0.1",
      "port": 25777,
      "username": "admin",
      "password": "admin",
      "tls": {
        "_comment": "TLS options only apply when using HTTPS protocol",
        "insecure_skip_verify": true,
        "certfile": "/path/to/public.pem",
        "keyfile": "/path/to/private.pem",
        "cacertfile": "/path/to/ca.pem"
      }
    }
  }
}
```

| Option                                    | Required   | Type          | Description
|:------------------------------------------|:----------:|:------------- |:-----------
| openstack.identity_endpoint               | Y          | String        | OpenStack Identify endpoint URI
| openstack.username                        | Y          | String        | OpenStack Username (Username is required if using Identity V2. In Identity V3, either User ID or a combination of Username and Domain ID or Domain Name are needed)
| openstack.user_id                         | Y          | String        | OpenStack UserID (Identity V3)
| openstack.password                        | Y          | String        | OpenStack Password (Exactly one of Password or API Key is required for the Identity V2 and V3)
| openstack.api_key                         | Y          | String        | OpenStack API Key (Exactly one of Password or API Key is required for the Identity V2 and V3)
| openstack.tenant_name                     | Y          | String        | OpenStack Tenant Name (Some providers allow you to specify a Tenant Name instead of the Tenant ID. Some require both)
| openstack.tenant_id                       | Y          | String        | OpenStack Tenant ID (Some providers allow you to specify a Tenant Name instead of the Tenant ID. Some require both)
| openstack.domain_name                     | Y          | String        | OpenStack Domain Name (At most one of Domain ID and Domain Name must be provided if using Username with Identity V3. Otherwise, either are optional)
| openstack.domain_id                       | Y          | String        | OpenStack Domain ID (At most one of Domain ID and Domain Name must be provided if using Username with Identity V3. Otherwise, either are optional)
| openstack.region                          | N          | String        | OpenStack Region
| openstack.default_keypair                 | N          | String        | Default OpenStack Key Pair to be used when creating servers
| openstack.default_security_groups         | N          | Array&lt;String&gt; | Default OpenStack Security Groups to be used when creating servers
| openstack.disable_config_drive            | N          | Boolean       | Disable injecting OpenStack user data via the Config Drive (`false` by default)
| openstack.disable_neutron                 | N          | Boolean       | Disable OpenStack Neutron interactions (`false` by default)
| openstack.ignore_server_availability_zone | N          | Boolean       | Ignore OpenStack Server's Availability Zone when creating OpenStack volumes. Commonly used if Ceph is used for block storage (`false` by default)
| actions.agent.mbus.endpoint               | Y          | String        | [BOSH Message Bus](http://bosh.io/docs/bosh-components.html#nats) URL used by deployed BOSH agents
| actions.agent.ntp                         | Y          | Array&lt;String&gt; | List of NTP servers used by deployed BOSH agents
| actions.agent.blobstore.type              | Y          | String        | Provider type for the [BOSH Blobstore](http://bosh.io/docs/bosh-components.html#blobstore) used by deployed BOSH agents (e.g. dav, s3)
| actions.agent.blobstore.options           | Y          | Hash          | Options for the [BOSH Blobstore](http://bosh.io/docs/bosh-components.html#blobstore) used by deployed BOSH agents
| actions.registry.protocol                 | Y          | String        | [BOSH Registry](http://bosh.io/docs/bosh-components.html#registry) Protocol (`http` or `https`)
| actions.registry.host                     | Y          | String        | [BOSH Registry](http://bosh.io/docs/bosh-components.html#registry) Host
| actions.registry.port                     | Y          | Integer       | [BOSH Registry](http://bosh.io/docs/bosh-components.html#registry) port
| actions.registry.username                 | Y          | String        | [BOSH Registry](http://bosh.io/docs/bosh-components.html#registry) username
| actions.registry.password                 | Y          | String        | [BOSH Registry](http://bosh.io/docs/bosh-components.html#registry) password
| actions.registry.tls.insecure_skip_verify | When https | Boolean       | Skip [BOSH Registry](http://bosh.io/docs/bosh-components.html#registry) server's certificate chain and host name verification
| actions.registry.tls.certfile             | When https | String        | [BOSH Registry](http://bosh.io/docs/bosh-components.html#registry) Client Certificate (PEM format) file location
| actions.registry.tls.keyfile              | When https | String        | [BOSH Registry](http://bosh.io/docs/bosh-components.html#registry) Client Key (PEM format) file location
| actions.registry.tls.cacertfile           | When https | String        | [BOSH Registry](http://bosh.io/docs/bosh-components.html#registry) Client Root CA certificates (PEM format) file location

### Run

Run the cpi using the previously created configuration file:

```
$ echo "{\"method\": \"method_name\", \"arguments\": []}" | cpi -configFile="/path/to/configuration_file.json"
```

## Features

### BOSH Network options

The BOSH OpenStack CPI supports these [BOSH Networks Types](http://bosh.io/docs/networks.html):

| Type    | Description
|:-------:|:-----------
| dynamic | To use dynamically assigned IPs by the OpenStack DHCP service
| manual  | To use manually assigned IPs
| vip     | To use previously allocated OpenStack Floating IPs

These options are specified under `cloud_properties` at the [networks](http://bosh.io/docs/networks.html) section of a BOSH deployment manifest and are only valid for `dynamic` and `manual` networks:

| Option          | Required | Type          | Description
|:----------------|:--------:|:------------- |:-----------
| network         | Y        | String        | The name of the OpenStack network to be used when creating servers (required when using OpenStack Neutron, optional otherwise)
| security_groups | N        | Array&lt;String&gt; | List of OpenStack security groups to be used when creating servers

### BOSH Resource pool options

These options are specified under `cloud_properties` at the [resource_pools](http://bosh.io/docs/deployment-basics.html#resource-pools) section of a BOSH deployment manifest:

| Option            | Required | Type   | Description
|:------------------|:--------:|:------ |:-----------
| flavor            | Y        | String | The name of the OpenStack flavor to be used when creating servers
| availability_zone | N        | String | The name of the OpenStack availability zone to be used when creating servers
| keypair           | N        | String | The name of the OpenStack keypair to be used when creating servers
| scheduler_hints   | N        | Hash   | List of OpenStack scheduler hints to be used when creating servers (see below for supported scheduler hints)

These are the list of supported `scheduler_hints` options:

| Option             | Required | Type          | Description
|:-------------------|:--------:|:------------- |:-----------
| group              | N        | String        | Server Group where the server will be placed
| different_host     | N        | Array&lt;String&gt; | Place the server on a compute node that does not host the given servers
| same_host          | N        | Array&lt;String&gt; | Place the server on a compute node that hosts the given servers
| query              | N        | String        | Conditional statement that results in compute nodes able to host the server
| target_cell        | N        | String        | Cell name where the server will be placed
| build_near_host_ip | N        | String        | Subnet of compute nodes to host the server

### BOSH Persistent Disks options

These options are specified under `cloud_properties` at the [disk_pools](http://bosh.io/docs/persistent-disks.html#persistent-disk-pool) section of a BOSH deployment manifest:

| Option            | Required | Type   | Description
|:------------------|:--------:|:------ |:-----------
| volume_type       | N        | String | The name of the OpenStack volume type to be used when creating volumes
| availability_zone | N        | String | The name of the OpenStack availability zone to be used when creating volumes

## Deployment Manifest Example

This is an example of how Google Compute Engine CPI specific properties are used in a BOSH deployment manifest:

```
---
name: example
director_uuid: 38ce80c3-e9e9-4aac-ba61-97c676631b91

...

networks:
  - name: private
    type: dynamic
    dns:
      - 8.8.8.8
      - 8.8.4.4
    cloud_properties:
      network: default
      security_groups:
        - bosh

  - name: public
    type: vip
    cloud_properties: {}
...

resource_pools:
  - name: vms
    network: private
    stemcell:
      name: bosh-openstack-kvm-ubuntu-trusty-go_agent
      version: latest
    cloud_properties:
      flavor: m1.medium
      availability_zone: az1
      keypair: bosh
      scheduler_hints:
        different_host:
          - a0cf03a5-d921-4877-bb5c-86d26cf818e1
          - 8c19174f-4220-44f0-824a-cd1eeef10287
...

disk_pools:
  - name: disks
    disk_size: 32_768
    cloud_properties:
      volume_type: ssd
      availability_zone: az1
...

```

## Contributing

In the spirit of [free software](http://www.fsf.org/licensing/essays/free-sw.html), **everyone** is encouraged to help improve this project.

Here are some ways *you* can contribute:

* by using alpha, beta, and prerelease versions
* by reporting bugs
* by suggesting new features
* by writing or editing documentation
* by writing specifications
* by writing code (**no patch is too small**: fix typos, add comments, clean up inconsistent whitespace)
* by refactoring code
* by closing [issues](https://github.com/frodenas/bosh-openstack-cpi/issues)
* by reviewing patches

### Submitting an Issue
We use the [GitHub issue tracker](https://github.com/frodenas/bosh-openstack-cpi/issues) to track bugs and features.
Before submitting a bug report or feature request, check to make sure it hasn't already been submitted. You can indicate
support for an existing issue by voting it up. When submitting a bug report, please include a
[Gist](http://gist.github.com/) that includes a stack trace and any details that may be necessary to reproduce the bug,
including your gem version, Ruby version, and operating system. Ideally, a bug report should include a pull request with
 failing specs.

### Submitting a Pull Request

1. Fork the project.
2. Create a topic branch.
3. Implement your feature or bug fix.
4. Commit and push your changes.
5. Submit a pull request.

## Copyright

Copyright (c) 2015 Ferran Rodenas. See [LICENSE](https://github.com/frodenas/bosh-openstack-cpi/blob/master/LICENSE) for details.
