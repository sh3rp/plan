## Quick Start

To build:

    go get -u github.com/sh3rp/plan/cmd/plan-server

To start the server, run:

    plan-server -pass <password>

This will start a plan webserver on port 8080 and load the database
directory from $HOME/.plan.  If a database doesn't exist, the server
will ask you for some information about yourself to include in all
requests for plans.

To add new plans, make a POST/PUT call to http://<host>:<port>/now 
with plan JSON as the body and add the header "x-plan-auth" with the 
specified <password>.  For example:

    curl -X POST -H "x-plan-auth: <password>" -d '{"body":"This is my plan!"}' http://localhost:8080/now

## API Endpoints

| Endpoint | Method | Description |
| -------- | ------ | --------------------------------------------------------------- |
| **/now** | GET | Retrieves the lastest posted plan |
| **/now** | POST | Posts a new plan; this can be the plan body, links, tags, or any combination of the three |
| **/all** | GET | Retrieves all posted plans, in order of descending timeline |
| **/plan/<ID>** | GET | Retrieves a posted plan by ID |
| **/info** | POST | Posts an update to the planner's information |
