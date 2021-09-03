# TWPT-server

TWPT-server is the code for the
[TW Power Tools](https://gerrit.avm99963.com/plugins/gitiles/infinitegforums/)
gRPC server and its frontend.

This repository is under active development, and is not ready yet for
production.

## Frontend
To serve the frontend for development, follow these instructions:

1. Run `docker-compose -f docker-compose.dev.yml up` to start the gRPC API,
database and Envoy proxy (which will translate gRPC-Web <-> gRPC).
1. Run `docker exec -i twpt-server_db_1 mysql -u root twpt < schema/common.sql`
and `docker exec -i twpt-server_db_1 mysql -u root twpt < schema/kill-switch.sql`.
1. Run `docker exec -it twpt-server_db_1 mysql -u root twpt` and enter the
following SQL sentence with your data to create the first authorized user:
`INSERT INTO KillSwitchAuthorizedUser (google_uid, email, access_level) VALUES ('', '{YOUR_EMAIL_ADDRESS}', 10);`.
1. Run `docker-compose -f docker-compose.dev.yml restart`.
1. Run `npm install`.
1. Run `make serve`.

## Backend
You can run the backend for development purposes by running `go run .` in the
`//cmd/serve` folder. You may also build its docker image to allow it to
interact with the frontend (since it needs Envy (the gRPC-Web proxy) to
translate between gRPC and HTTP requests).

## Deploying
Use the `docker-compose.yml` file instead of `docker-compose.dev.yml` and follow
steps 1-4 from the "frontend" section to spin up the backend.

For the frontend, run `make deploy` in the `//frontend` folder to deploy the
frontend to Firebase.
