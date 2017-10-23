To build:

    go get -u github.com/sh3rp/plan/cmd/plan-server

To start the server, run:

    plan-server -port <port> -dbdir <db_directory> -pass <password>

This will start a plan webserver on port specified by <port> and load
the database directory from <db_directory>.  To add new plans, make a
POST/PUT call to http://<host>:<port>/now with plan JSON and add the
header "x-plan-auth" with the specified <password>.
