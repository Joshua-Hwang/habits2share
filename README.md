# Habits2Share webserver
This is a webserver for a habit sharing app.

Currently it relies on a persistent file system.

## Running locally
You need Golang to run the server and you need `npm` to build the frontend.

In order to test the login flow you will need to generate a Google Cloud OAuth
token. You may skip the login flow by following the steps below. In future this
should not be necessary.

Create a file called `accounts.json` which is an array of `{"Id":"user
id","Email": " email here"}` objects. If you don't need to test the login flow
you can create a `sessions.csv` file of the form
```csv
insert some uuid,insert userId here,some date of the form 2006-01-02T15:04:05Z07:00 (RFC3339)
```

When making `curl` requests add the following cookie
```bash
curl --cookie __Secure-SESSIONID='uuid mentioned before' localhost:8080...
```

1. Run `./scripts/start.sh` to run the server. We currently don't have any dependencies.

## Building container
The `Dockerfile` doesn't automatically build the `frontend/` so the
`./scripts/build-container.sh` does that prior building the container.

## TODO
* [ ] Remove need for Google Cloud OAuth
* [X] Sharing habits
* [ ] Document API
* [ ] All users are considered friends. Provide a friending system.
* [X] Import csv from HabitShare app
* [ ] Create todo section

