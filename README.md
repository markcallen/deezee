# deezee

Docker task runner.  Run tasks using docker images and pass information between
them.  Replace shell scripts and no need to install local programs.

## Get

````
go get github.com/markcallen/deezee
````

## Run

````
deezee run --file=example/whoami.yml
````

or in Docker

````
docker run --rm -v ${PWD}/examples:/deezee -v /var/run/docker.sock:/var/run/docker.sock markcallen/deezee run --file=/deezee/whoami.yml
````

## File format

````
version: "1"

globals:
  PUBLIC_IP: null

processes:
- name: whoami
  tasks:
  - name: getip
    image: byrnedo/alpine-curl
    command: -Ls http://ipinfo.io/ip
    result: PUBLIC_IP
  - name: whois
    image: airtonzanon/whois:latest
    command: ${PUBLIC_IP}
````

### Version

````
version: "1"
`````

The version of the file.  Currently only 1 is supported.

### Globals

````
globals:
  PUBLIC_IP: null
````

Global variables that need to be declared in order to be set in a result and used
in commands, environment or volumes.  Global variables can be used in the
following format ${<global variable>}, for a result in JSON the values can be
accessed using ${<global variable>.key1.key2}.

Environment variables can be used with the format $env{<VAR>}.  For example to pass $PWD
into the process use $env{PWD}.

### Processes

````
processes:
- name: whoami
````

There can be multiple processes containing one or many tasks, each one can have
a unique name.

### Tasks

````
- name: getip
  image: byrnedo/alpine-curl
  command: -Ls http://ipinfo.io/ip
  result: PUBLIC_IP
````

Each task needs a name, docker image and command.

Any output can be assigned to a global variable using the result.

Volumes can be mounted using:

````
volumes:
  - $env{PWD}/packer:/terraform
````

Environment variables can be set in the container using:

````
environment:
  - AWS_ACCESS_KEY_ID=$env{AWS_ACCESS_KEY_ID}
````

## Development

Get all the packages for development

````
dep ensure
````

## License & Authors
- Author:: Mark Allen (mark@markcallen.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
