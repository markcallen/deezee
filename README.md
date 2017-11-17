#dz

Docker task runner

## Get

go get github.com/markcallen/deezee

## Run

deezee run --file=example/whoami.yml

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
