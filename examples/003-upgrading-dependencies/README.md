# Location

This Dockerised service reads NMEA GPS data from a serial device and posts it
as a protocol buffer to a Redis server with a given key.

## Base setup

The host system must have a `/dev/gps` device. A udev rules file is provided in
`rules.d/99-gps.conf` that can be symlinked or copied to
`/etc/udev/rules.d/99-gps.conf` when using a BU-353-S4 receiver; it should be
easily adapted for any USB serial device by changing the vendor and product
IDs.

## Building

### Dependency management

This project uses [Glide](https://glide.sh/) to manage dependencies. Due to the
use of private repositories for Æcerbot dependencies, this will need to be run
on the building machine before building the Docker image.

Once Go and Glide are installed, this should be as simple as:

```bash
GOPATH=$(pwd) glide install
```

### Docker

Building the Docker image should be straightforward once the dependencies are
installed:

```bash
docker-compose build
```

Note that in most cases you'll want to actually run this from within the
`brain` project, but building it locally won't cost any significant time due to
Docker's caching.

## Running

By default, the location service will try to connect to `redis:6379`, use the
`position` key when updating, and use `/dev/gps` as the serial device. All of
these are configurable via command line options.

## Development

The `docker-compose.yml` file sets this project up in development mode by
setting up a volume mount for the current directory. Starting from the host, this should get it up and running in development mode:

```bash
GOPATH=$(pwd) glide install
docker-compose build
docker-compose run --rm location bash
go install location
./bin/location
```

You can also run on the host if Go is installed (which it probably is):

```bash
GOPATH=$(pwd) glide install
GOPATH=$(pwd) go install location
./bin/location
```
