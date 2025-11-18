#!/usr/bin/env bash
set -e

URL_FMT="https://download.db-ip.com/free/dbip-country-lite-%s.csv.gz"
ARCHIVE_FMT="/tmp/dbip-country-lite-%s.csv.gz"

function show_usage() {
  return
}

function unknown_command() {
  echo "Unknown command $1."

  show_usage
}

function fetch_unpack_csv() {
  local month="$1"

  if [ -z "$month" ]; then
    echo "Month is not specified"
    exit 1
  fi

  local url
  url=$(awk "{printf \"$URL_FMT\", \$1}" <<<"$month")

  local archive
  archive=$(awk "{printf \"$ARCHIVE_FMT\", \$1}" <<< "$month")

  local csv
  csv=$(awk '{sub(/\.gz$/, "", $0); print}' <<< "$archive" )

  echo "Downloading: $url to $archive"
  curl "$url" --output "$archive"

  gzip -d < "$archive" > "$csv"
  rm "$archive"
}

function fetch_and_generate() {
  fetch_unpack_csv "$1"
}

function run() {
  local main_cmd=${1:-"help"}
  shift

  echo -n "Running command "
  echo "$main_cmd"
  echo ""

  case $main_cmd in
  run:dev)
    run_dev "$@"
    ;;
  code:fetch-and-generate)
    fetch_and_generate "$@"
    ;;
  cert:install)
    cert_install "$@"
    ;;
  key:generate_ec)
    key_gen_ec "$@"
    ;;
  *)
    unknown_command "$main_cmd"
    ;;
  esac
}

run "$@"
