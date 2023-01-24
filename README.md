# Habits2Share webserver
This is a webserver for a habit sharing app.

Currently it relies on a persistent file system.

Store personal files in the `gitignored/` folder.

## Running locally
You won't need environment variables for the `CLIENT_ID`s. Those can be
disabled locally by using the correct build constraints.

When running `scripts/start.sh` runs with the build constraint `dev`. This
means it will present a dev page for logins. No actual checking is done and you
don't need the `CLIENT_ID`s.

If you'd like to develop the login flow remove the `dev` from the `GOFLAGS`.

Configuration files and database files need to be created. To avoid adding them
accidentally to the repo run the `./scripts/init.sh` which will create
necessary folders and files that will be ignored by git.

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

Debugging with [Delve](https://github.com/go-delve/delve)

1. Run `./scripts/start.sh` to run the server. We currently don't have any dependent services.

## Generate test data
Test data can be generated with the scripts present in `scripts/test-generation`

The test generation relies on the existing test suite (because registering new accounts isn't possible).

## Building container
The `Dockerfile` doesn't automatically build the `frontend/` so the
`./scripts/build-container.sh` does that prior building the container.

## Testing
Run all unit tests with `go test ./...` or with `./scripts/test-unit-tests.sh`

Run integration tests with `./scripts/test-integration-tests.sh`. Ensure the service isn't running.

## Tips and tricks
`mockgen` is what I've used in the past to generate mocks. All interfaces are free game to mock.

[ ] Add steps here on how to mock

## TODO
* [ ] Document API
* [ ] All users are considered friends. Provide a friending system.

* [ ] Load testing to ensure data races don't pop up. Since only one instance of the app runs at a time the data races aren't as bad.
It's easier to solve these problems with an online database solution compared to a single file.
* [ ] Fix fragility of system on malformed data. Single file is parsed so any failures there brings the whole thing down.
production habits.json is 32K, for now we could do the swap method (maybe?)
* [ ] Same is true for the sessions.csv

* [ ] Integration tests could run on random port and data generators could take PORT
* [ ] Load testing framework would be a huge confidence boost
* [ ] Code coverage enforcement on deployment

* [ ] We should incorporate a logout button as well.
* [ ] Create endpoint for querying own information
* [ ] Use that info to show a login screen or not
* [ ] Excessive logging by default
* [ ] Ensure no sedentary code
