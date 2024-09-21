#!/usr/bin/env bash
#
# Create https certs for local development
#

# Exit on error
set -e

CERTS_DIR="dockerfiles/nginx/certs"
CERT_NAME="docker-self-signed"
CERT_DAYS="100"
OPENSSL_CONFIG="dockerfiles/nginx/openssl.cnf"
COMMON_NAME="Docker development"
PASSWORD_PROMPT="Enter password: "

function is_macos {
  [ "$(uname)" == "Darwin" ]
}

function cert_exists {
  local cert_file="$CERTS_DIR/$CERT_NAME.crt"
  local key_file="$CERTS_DIR/$CERT_NAME.key"
  [ -f "$cert_file" ] && [ -f "$key_file" ]
}

function cert_expired {
  local cert_file="$CERTS_DIR/$CERT_NAME.crt"
  local current_date=$(date "+%s")
  local expiration_date=$(openssl x509 -enddate -noout -in "$cert_file" | cut -d= -f2)
  expiration_date=$(date -j -f "%b %d %T %Y %Z" "$expiration_date" "+%s")
  [ "$current_date" -gt "$expiration_date" ]
}

function create_cert {
  rm -f -- $CERTS_DIR/$CERT_NAME.crt $CERTS_DIR/$CERT_NAME.key
  openssl genrsa -out $CERTS_DIR/$CERT_NAME.key 2048
  openssl req \
    -new \
    -x509 \
    -sha256 \
    -config $OPENSSL_CONFIG \
    -key $CERTS_DIR/$CERT_NAME.key \
    -out $CERTS_DIR/$CERT_NAME.crt \
    -days $CERT_DAYS
}

function add_cert_to_keychain {
  if is_macos; then
    echo "Adding docker cert to keychain"
    sudo -p "$PASSWORD_PROMPT" security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain $CERTS_DIR/$CERT_NAME.crt
  else
    # For Linux
    echo "Adding docker cert to ca-certificates"
    sudo -p "$PASSWORD_PROMPT" cp $CERTS_DIR/$CERT_NAME.crt /usr/local/share/ca-certificates/$CERT_NAME.crt
    sudo -p "$PASSWORD_PROMPT" update-ca-certificates
  fi
}

function remove_cert_from_keychain {
  function cert_hashes {
    local cn="$1"
    sudo -p "$PASSWORD_PROMPT" security find-certificate -a -c "$cn" -Z | grep ^SHA-256 | cut -d' ' -f3
  }

  if is_macos; then
    echo "Removing docker cert from keychain"
    for SHA in $(cert_hashes "$COMMON_NAME"); do
      sudo -p "$PASSWORD_PROMPT" security delete-certificate -Z $SHA
    done
  else
    # For Linux
    echo "Removing docker cert from ca-certificates"
    sudo -p "$PASSWORD_PROMPT" rm -f /usr/local/share/ca-certificates/$CERT_NAME.crt
    sudo -p "$PASSWORD_PROMPT" update-ca-certificates
  fi
}

function install_linux_deps {
  is_macos && return

  function is_ca_certificates_installed {
    command -v update-ca-certificates >/dev/null
  }

  if ! is_ca_certificates_installed; then
    echo "Installing ca-certificates package"
    sudo -p "$PASSWORD_PROMPT" apt-get update -y
    sudo -p "$PASSWORD_PROMPT" apt-get install -y ca-certificates
  fi
}

function create_ssl_certificates {
  is_macos || install_linux_deps

  if ! cert_exists || cert_expired; then
    remove_cert_from_keychain
    create_cert
    add_cert_to_keychain
  fi
}

create_ssl_certificates
exit 0
