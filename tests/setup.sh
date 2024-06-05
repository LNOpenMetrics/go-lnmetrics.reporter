#!/bin/bash
set -x

latest_release_tag=$(grep '^[[:space:]]*version:' coffee.yaml | awk '{print $2}' | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')

# Get the directory of the script
script_dir=$(dirname -- "$(readlink -f -- "$0")")

name="go-lnmetrics.reporter"
repo="LNOpenMetrics"


get_platform_file_end() {
    machine=$(uname -m)
    kernel=$(uname -s)

    case $kernel in
        Darwin)
            echo 'darwin-amd64'
            ;;
        Linux)
            case $machine in
                x86_64)
                    echo 'linux-amd64'
                    ;;
                armv7l)
                    echo 'linux-arm'
                    ;;
                aarch64)
                    echo 'linux-arm64'
                    ;;
                *)
                    echo "No self-compiled binary found and unsupported release-architecture: $machine" >&2
                    exit 1
                    ;;
            esac
            ;;
        *)
            echo "No self-compiled binary found and unsupported OS: $kernel" >&2
            exit 1
            ;;
    esac
}
platform_file_end=$(get_platform_file_end)
archive_file="go-lnmetrics-$platform_file_end"

github_url="https://github.com/$repo/$name/releases/download/v$latest_release_tag/$archive_file"


# Download the archive using curl
if ! curl -L "$github_url" -o "$script_dir/go-lnmetrics"; then
    echo "Error downloading the file from $github_url" >&2
    exit 1
fi

chmod u+x "$script_dir/go-lnmetrics"
