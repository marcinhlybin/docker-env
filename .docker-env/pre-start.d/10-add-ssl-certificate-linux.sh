#!/usr/bin/env bash
#
# Create https certs for local development
#

# Exit on error
set -e

CERTS_DIR="dockerfiles/nginx/certs"
CERT_NAME="docker-development"
COMMON_NAME="Docker development"
CREATE_SSL_CERTIFICATE_SCRIPT="dockerfiles/nginx/create-ssl-certificate.sh"
PASSWORD_PROMPT="Enter password: "

# Certificate paths
CERT_FILE="$CERTS_DIR/$CERT_NAME.crt"
KEY_FILE="$CERTS_DIR/$CERT_NAME.key"

function is_linux {
  [ "$(uname)" == "Linux" ]
}

function cert_file_not_found {
  [ ! -f "$CERT_FILE" ]
}

function cert_expired {
    local current_timestamp="$(date "+%s")"
    local expiration_date="$(openssl x509 -enddate -noout -in "$CERT_FILE" | cut -d= -f2)"
    expiration_timestamp="$(date -d "$expiration_date" "+%s")"
    [ "$current_timestamp" -gt "$expiration_timestamp" ]
}

function cert_in_ca_certificates {
  [ -f "/usr/local/share/ca-certificates/$CERT_NAME.crt" ]
}

function add_cert_to_ca_certificates {
    if ! cert_in_ca_certificates; then
        echo "Adding certificate to ca-certificates"
        sudo -p "$PASSWORD_PROMPT" cp $CERT_FILE /usr/local/share/ca-certificates/$CERT_NAME.crt
        sudo -p "$PASSWORD_PROMPT" update-ca-certificates >/dev/null 2>&1
    fi
}

function remove_cert_from_ca_certificates {
    if cert_in_ca_certificates; then
      echo "Removing certificate from ca-certificates"
      sudo -p "$PASSWORD_PROMPT" rm -f /usr/local/share/ca-certificates/$CERT_NAME.crt
      sudo -p "$PASSWORD_PROMPT" update-ca-certificates >/dev/null 2>&1
    fi
}

function install_dependencies {
  function is_ca_certificates_installed {
    command -v update-ca-certificates >/dev/null
  }

  if ! is_ca_certificates_installed; then
    echo "Installing ca-certificates package"
    sudo -p "$PASSWORD_PROMPT" apt-get update -y
    sudo -p "$PASSWORD_PROMPT" apt-get install -y ca-certificates
  fi
}

function create_cert_file {
  if [ ! -f "$CREATE_SSL_CERTIFICATE_SCRIPT" ]; then
    echo "SSL certificate script not found: $CREATE_SSL_CERTIFICATE_SCRIPT"
    exit 1
  fi

  /usr/bin/env bash $CREATE_SSL_CERTIFICATE_SCRIPT
}

function main {
    # Ensure it runs on Linux
    is_linux || return 0

    echo "Checking SSL certificates"
    if cert_file_not_found || cert_expired; then
        install_dependencies
        remove_cert_from_ca_certificates
        create_cert_file
        add_cert_to_ca_certificates
    fi
}

main
exit 0
