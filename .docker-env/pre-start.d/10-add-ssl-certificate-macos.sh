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

function is_macos {
  [ "$(uname)" == "Darwin" ]
}

function cert_file_not_found {
  [ ! -f "$CERT_FILE" ]
}

function cert_expired {
    local current_timestamp="$(date "+%s")"
    local expiration_date="$(openssl x509 -enddate -noout -in "$CERT_FILE" | cut -d= -f2)"
    local expiration_timestamp="$(date -j -f "%b %d %T %Y %Z" "$expiration_date" "+%s")"
    [ "$current_timestamp" -gt "$expiration_timestamp" ]
}

function cert_in_keychain {
  [ -n "$(get_cert_hashes_from_keychain "$COMMON_NAME")" ]
}

function get_cert_hashes_from_keychain {
  local common_name="$1"
  sudo -p "$PASSWORD_PROMPT" security find-certificate -a -c "$common_name" -Z | grep ^SHA-256 | cut -d' ' -f3
}

function add_cert_to_keychain {
    if ! cert_in_keychain; then
      echo "Adding SSL certificate to keychain"
      sudo -p "$PASSWORD_PROMPT" security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain $CERT_FILE
    fi
}

function remove_cert_from_keychain {
  echo "Getting SSL certificates from keychain"
  cert_hashes="$(get_cert_hashes_from_keychain "$COMMON_NAME")"
  for hash in $cert_hashes; do
    echo "Removing old SSL certificate from keychain"
    sudo -p "$PASSWORD_PROMPT" security delete-certificate -Z $hash
  done
}

function create_cert_file {
  if [ ! -f "$CREATE_SSL_CERTIFICATE_SCRIPT" ]; then
    echo "SSL certificate script not found: $CREATE_SSL_CERTIFICATE_SCRIPT"
    exit 1
  fi

  /usr/bin/env bash $CREATE_SSL_CERTIFICATE_SCRIPT
}

function main {
  # Ensure it runs on MacOS
  is_macos || return 0

  echo "Checking SSL certificates"
  if cert_file_not_found || cert_expired; then
    remove_cert_from_keychain
    create_cert_file
    add_cert_to_keychain
  fi
}

main
exit 0
