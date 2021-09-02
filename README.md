# TWPT-server

TWPT-server is the code for the
[TW Power Tools](https://gerrit.avm99963.com/plugins/gitiles/infinitegforums/)
gRPC server and its frontend.

This repository is under active development, and is not ready yet for
production.

## Frontend
To serve the frontend, follow these instructions:

1. Run `docker-compose -f docker-compose.dev.yml up` to start the gRPC API,
database and Envoy proxy (which will translate gRPC-Web <-> gRPC).
1. Run `mysql -u root twpt < schema/*.sql`.
1. Run `docker-compose -f docker-compose.dev.yml restart`.
1. Run `npm install`.
1. Run `make serve`.
