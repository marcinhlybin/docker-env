# Docker environments manager

**docker-env** allows developers to create, manage and switch between isolated docker projects (or stacks) per git branch within a single repository. It simplifies working with docker-compose and supports hooks for customization and sidecar containers for optional services.

![demo](https://github.com/user-attachments/assets/52289faf-8d40-42dc-8670-b5260ccfedc6)

## Why

Managing Docker environments across multiple branches or projects can be cumbersome, especially when dealing with complex application stacks that include databases, caches, and other services.

The goal of docker-env is to streamline this process by automating the creation and management of docker-compose environments. Each environment is tied to a Git branch, allowing developers to easily switch between isolated stacks without affecting their work in progress. This ensures that data, services, and application states are preserved across branches and projects.

In short, docker-env abstracts common docker-compose tasks, allowing you to focus on development while it handles the heavy lifting.

## Key features

* Project (or stack) per branch. Create a new environment automatically when switching branches.
* Repository isolation: Prefix project by repository name, ensuring no conflicts between different repositories.
* Sidecar containers: Optional services, like admin tools or background jobs, can be started as needed without starting by default.
* Hooks: Customize pre-start, post-start, and post-stop behaviors with hooks.
* Multi-environment support: Easily manage multiple Docker Compose environments across projects.

## Hooks

Sample hooks can be found in `.docker-env/` directory.

Supported hooks are:
* pre-start
* post-start
* post-stop

## Usage

```
NAME:
   docker-env - Docker environments manager

USAGE:
   docker-env [global options] command [command options]

VERSION:
   1.0

DESCRIPTION:
   All commands must run in the git repository directory.
   If project name is not specified current branch name is used.

COMMANDS:
   start, s, up                Start docker containers
   stop, ss, down              Stop docker containers
   restart, r, reboot          Restart docker containers
   remove, rm, delete          Remove docker containers
   ls, list, l, ll             List projects. Use 'll' to show containers.
   cleanup                     Cleanup entire project
   build, b                    Build docker images
   info, config, show          Show configuration
   terminal, term, shell, ssh  Run terminal
   code, open                  Open code editor
   help, h                     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, -c value  config file path
   --debug, -d               enable debug mode (default: false)
   --help, -h                show help
   --version, -v             print the version
```

Start command usage:

```
NAME:
   docker-env start - Start docker containers

USAGE:
   docker-env start [command options]

DESCRIPTION:
   Start docker containers.
   If project name is not specified, current branch name is used.
   If project does not exist it will be created.

OPTIONS:
   --project value, -p value  set a project name
   --service value, -s value  start a single service
   --recreate, -r             recreate the containers (default: false)
   --update, -u               update the images and recreate the containers (default: false)
   --help, -h                 show help
```

## Sample commands

```
# Create new environment based on branch name
docker-env start

# Create new environment with custom name
docker-env start -p db-fix

# Restart environment
docker-env restart -p db-fix

# Restart a single container
docker-env restart -p db-fix -s app

# Recreate all containers
docker-env start -r

# Recreate a single container
docker-env start -s postgresql -r

# Update all images and recreate containers
docker-env start -u

# Update image of a single container and recreate
docker-env start -s app -u

# Cleanup environments and images
docker-env cleanup --with-images

# Run shell
docker-env shell

# Run command in a container
docker-env shell -s postgresql createdb -U postgres mydb
```

## Configuration

Each repository should define its own configuration file located in `./docker-env/config.yml`

```
### Docker-env configuration file

# Docker compose projects and their containers will be prefixed with this name
# Only alphanumeric characters and underscores are allowed, no hyphens
# Make sure to name service names in the docker-compose.yml file
# using "$COMPOSE_PROJECT_NAME-" prefix
compose_project_name: docker_env

# Docker compose configuration
compose_file: docker-compose.yml
compose_file_override: docker-compose.override.yml

# Profiles are used to distinguish between default startup services
# and services that are only started when explicitly requested by the user
# so called sidecar services
compose_default_profile: app
compose_sidecar_profile: sidecar

# Debug options
show_commands: true

# Env files to load environmental variables used in the docker-compose.yml file
# for substitution in the services section
env_files:
  - .env

# Check for following environment variables in env files
required_vars:
  - GITHUB_USER
  - GITHUB_TOKEN
  - AWS_ACCESS_KEY_ID
  - AWS_SECRET_ACCESS_KEY

# AWS registry
aws_login: true
aws_region: eu-central-1
aws_repository: 1234567890.dkr.ecr.eu-central-1.amazonaws.com

# Command defaults
terminal_default_service: app
terminal_default_command: /bin/bash
vscode_default_service: app
vscode_default_dir: /app

# Scripts to run before and after
pre_start_script: .docker-env/pre-start.sh
post_start_script: .docker-env/post-start.sh
post_stop_script: .docker-env/post-stop.sh
```

## Building

Run `make` to build a binary to the current directory.

```
make test
make install
```

Installs into `/usr/local/bin`. Sudo password required.
