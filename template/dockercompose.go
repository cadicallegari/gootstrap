package template

const DockerCompose = `# used for dev purpose only
version: '3'

services:

  {{.Project}}:
    image: ${docker_image}

`
