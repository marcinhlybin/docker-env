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


## Installing

On MacOS use Homebrew to install:

```
brew install marcinhlybin/docker-env/docker-env
```

On Linux see building or download binary from releases.

## Building

Run `make` to build a binary to the current directory.

```
brew install go
make test
make install
```

Installs into `/usr/local/bin`. Sudo password required.

## Usage

```
NAME:
   docker-env - Docker environments manager

USAGE:
   docker-env [global options] command [command options]

VERSION:
   1.0.0

DESCRIPTION:
   All commands must run in the git repository directory of the project.
   If environment name is not specified current branch name is used.

COMMANDS:
   start, s, up                Start docker containers
   stop, ss, down              Stop docker containers
   restart, r, reboot          Restart docker containers
   remove, rm, delete          Remove docker containers
   ls, list, l, ll             List projects, 'll' to show containers.
   cleanup                     Removes all projects
   build, b                    Build docker images
   info, config, show          Show configuration
   terminal, term, shell, ssh  Run terminal
   code, open                  Open code editor
   version, v                  Show version
   logs, log                   Show container logs
   help, h                     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, -c value  config file path
   --debug, -d               enable debug mode (default: false)
   --quiet, -q               disable info messages (default: false)
   --quieter, --qq           disable docker output (default: false)
   --help, -h                show help
   --version, -v             show version string, alias for 'version --short' (default: false)
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

# Show container logs with follow option
docker-env logs -s nginx -f
```

## Docker compose

### Requirements

Container name in service definition **must be prefixed** with `$COMPOSE_PROJECT_NAME-`:

```
services:
  app:
    container_name: $COMPOSE_PROJECT_NAME-app
    image: 1234567890.dkr.ecr.eu-central-1.amazonaws.com/my-app
    ports:
    [...]
```

### Extending docker compose configuration

By default **docker-env** uses `docker-compose.yml` for service definition. The configuration can be extended with `docker-compose.override.yml`.

### Variables

Docker compose configuration can be manipulated with environmental variables stored in `env_files` defined in `./docker-env/config.yml` file. Usually it is local `.env` file which should be added to `.gitignore` as it may contain secrets like `GITHUB_TOKEN` or AWS credentials.

By setting `required_vars` option in `./docker-env/config.yml` you can check during runtime if variables can be found.

### Sidecar containers

By default **docker-env** uses `app` profile to start contianers. Profile `sidecar` is used to specify containers that are not started by default and can be started later individually using `docker-env start -s` command.

Profile names are defined in the configuration file:

```
compose_default_profile: app
compose_sidecar_profile: sidecar
```

## Configuration

Each repository should define its own configuration file located in `./docker-env/config.yml`. It is recommended to commit config file to the repository. It can be extended locally with `./docker-env/config.override.yml` (add it to `.gitignore`). For example it can be used to replace registry images and build own using `docker-env build` command.

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
show_executed_commands: true

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
aws_login: false
aws_region: eu-central-1
aws_repository: 1234567890.dkr.ecr.eu-central-1.amazonaws.com

# Command defaults
terminal_default_service: app
terminal_default_command: /bin/bash
vscode_default_service: app
vscode_default_dir: /app

# Scripts to run before and after
pre_start_hooks: 
  - .docker-env/pre-start.sh

post_start_hooks: 
  - .docker-env/post-start.sh

post_stop_hooks: 
  - .docker-env/post-stop.sh

```

## Hooks

Supported hooks are:
* pre-start
* post-start
* post-stop

Arguments passed to the hooks are `PROJECT_NAME` and `SERVICE_NAME` as positional arguments.

```
# In bash script use:
PROJECT_NAME="$1"
SERVICE_NAME="$2"
```

Sample hooks can be found in `.docker-env/` directory to:
* Generate SSL certificates
* Run ssh-agent
* Check ports availability


## Troubleshooting

### Error running docker command: Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?

Set `DOCKER_HOST` to docker socket in your home directory:

```
export DOCKER_HOST="unix:///$HOME/.docker/run/docker.sock"
```

or in Docker Desktop in Settings -> Advanced select `Allow the default Docker socket to be used (requires password)`

### Error getting AWS registry password: Unable to locate credentials. You can configure credentials by running "aws configure".

Export `AWS_PROFILE` variables matching the profile in `~/.aws/credentials`:

```
export AWS_PROFILE=default
```

## Release

```
git tag -a v1.0.0 -m "First release"
goreleaser release --clean
```