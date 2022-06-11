# Habits2Share webserver
This is a webserver for a habit sharing app.

Currently it relies on a persistent file system.

## Running locally
For now you will need to generate a Google Cloud OAuth token. In future this should not be necessary.

Create a file called `accounts.json` which is an array of `{"Email": "insert email here"}` objects.

1. Run `./scripts/start.sh` to run the server. We currently don't have any dependencies.

## TODO
* [ ] Remove need for Google Cloud OAuth
* [ ] Sharing habits
