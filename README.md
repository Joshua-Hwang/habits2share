# Habits2Share webserver
This is a webserver for a habit sharing app.

Currently it relies on a persistent file system.

## Running locally
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
curl --cookie __Host-SESSIONID='uuid mentioned before' localhost:8080...
```

1. Run `./scripts/start.sh` to run the server. We currently don't have any dependencies.

## TODO
* [ ] Remove need for Google Cloud OAuth
* [ ] Sharing habits
* [ ] Document API
