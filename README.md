As part of revisiting old code. Big changes will be made here.

[O] Improve and simplify the `.gitignore` to inspire more confidence in what is being shared in this public repo
[O] Actually fix the gross authentication logic
[O] Develop system to copy production data into local
[O] Generate scripts that generate particular types of data
[O] Use these scripts inside integration tests
[O] script to run all tests or document how to run tests. Things to inspire confidence when running this app.
[/] curl commands are annoying and finicky. Scripts or otherwise for command in just the CLI
[ ] Integration tests could run on random port and data generators could take PORT
[ ] Load testing framework would be a huge confidence boost
[ ] Code coverage enforcement on deployment

[O] Others who attempt to run the script will find the same problems with `.gitignored/` and `.env` because these are ignored from the repo.
Create a script for these.
[ ] Connect locally running frontend to prod and connect locally running backend to frontend (service overrides?)
It's probably fine right now given how lightweight the app is (and if the data generation scripts are good enough).
Are service overrides an anti-pattern? Things become big ball of mud if backend services talk to each other and don't know about service override
[ ] We should incorporate a logout button as well.
[ ] absolute path packages freak me out. What happens if we move the repo or someone forks it?
[X] Changes to init.sh need to be encourage and thought about. How can developers remember to consider the init.sh?
[ ] Excessive logging by default
[ ] Ensure no sedentary code

# Thinking through problems

Currently I'm using the .env file for pointing the service to the right file paths.
I can see a time when an insane number of configuration options become available.
Environment variables no longer become an elegant option. Chances of typos and
malformed input shoot through the roof.
A: That time is likely not going to happen in the near future. I'll kick this
problem down the road.

## How will authentication work?
The app is dependent on sessions to distinguish users.
It feels very orthogonal but not as orthogonal as I'd like.

If it was truly orthogonal I could just remove the authentication step during
testing. Maybe even have an authentication sidecar in-front of our service.

## What even is the actual problem?
Current session management is clunky and requires I manually add the cookie
during testing in the browser or I use my real client ID.
Perhaps using the real client ID is fine but for setup that is clunky and
becomes another point of struggling to get the service working.

What it even does or how it behaves isn't clear. For something so simple that shouldn't be the case.

The authenticator will work as a middleware on all endpoints. If it doesn't find the session cookie
it will return 403 and show a link to the login page. The only exception is the log in page of course.
Accessing the static resources will also be behind that, though it doesn't need to be.
This prevents redirection which I'm not a fan of for REST endpoints that a CLI might interact with.

Do we need a login flow for the frontend? Currently log in is done the old-school way.

What if I create a "session manager" frontend that provides buttons for each user listed in `user.json`.
Clicking the button sets the cookie.
Log out by clearing cookies.

## How will we inject the right authentication layer?
At runtime? Should we have a script that modifies the code?

The script to generate code sounds like an awful idea given how much that code might change
and chance of deploying to environments is pretty high. Not confidence inspiring.

Runtime configuration sounds like an unnecessary security hole.

Is it possible to build with specific dependencies?
GoLang doesn't like macros. I wonder what their solution to something like this is.
Build tags?

Two options with build tags. Multiple `main` and go builds them or single
`main` and functions are considered the real `init_auth()`.

### TokenParser
`TokenParser` is only used for the `PostLogin`. If we're willing to turn off
the production login page then `TokenParser` isn't necessary.

I can see a situation where a new developer thinks they need to create Google
OAuth tokens in order to begin developing.

How can we communicate they don't need to do that (unless working on the login flow)

Feels like the situation Rich Hickey was thinking of with his talk Maybe Not.

If we create two versions of `Server` for dev and for prod but that could
explode in complexity if there's another part of the program we'd like to
turn off.

Certainly the wrong system for the future but for now where the only thing
worth turning off is the login. Maybe it's worth it.

That's not a good way to think about it. People like copying. Thinking of a
better system is unlikely.

If a proper solution exists then I don't expect it to take much more time than
the current option.

Okay, whatever. I'll just write something in the `CONTRIBUTING` that explains
the environment variables aren't necessary.

## How can developers remember to consider the init.sh?
Add a comment to all the places. There aren't that many

## Using the build tags effectively

## I'm not liking dependency injection
Constantly need to check the dependency is actually created and is there.
What happens if the dependency isn't there and we just continue as if it were?

Okay moved to another dependency injection and liking the compile time stuff.
Problem is it's explicit for the parmeters we don't need. Might just be a fact of Go though.

## Operations?
Currently we return 5XX on unexpected errors. These 5XXs don't alert me though.
Also adding the `ok` logic everywhere is getting kind of annoying.

We could change from `ok` to `panic`s which seem more appropriate in this situation.
Will have to learn how to handle panics `defer recover`.

What will happen in this panic? Log something? Seems appropriate.
How about alerting via something like opsgenie? Seems like a lot of work given
I'm not going to use that. If I'm just logging something it's equivalent
to the logging already being done.

I think I'm going to just let the panic happen.

I'm also unhappy that we can only determine this at runtime. Seems BS for a compiled language.
I wonder if there's a better way?
Use what they did in https://github.com/Khan/typed-context/tree/main/07-server-interface
Change the http-helper to accept three parameters instead of the normal two.
Make the dependency injector a more integrated part of the helper.

They didn't have/need observability and monitoring when software was distributed.
Why suddenly now with the service?

flyctl specific commands only work when I'm in the `secrets/` folder. Maybe
that's an okay tradeoff for a gentler introduction

## How to run the service in multiple ways?
Probably going to write several scripts.

First will probably be running without the proper login screen.

### How to control the login and turn it off
Given the above it would be good to control the login based on compilation.

## Should I use the init function?
It's an integral part of the language so I will assume the average Go developer should know about it and look for it.

A clueless developer would be mightily confused when debugging. Will need to
include a section in the CONTRIBUTING or README about the init function.

Runtime investigation will also reveal it if I add a log message there.

Within the `init` function for dependency I take a number of environment variables.
It's no longer clear where all the environment variables are or if a variable
collision might occur.

Solution might be to create a global config struct which gets populated on first retrieval.
This will mean environment variables are not retrieved immediately but a bit later.
This should be fine as it's very difficult to change the environment variables of a running process.

## How do we stop people from using the global config?
We need to encourage all accesses via the deps. This should be easy if we have
a culture of unit testing.
Messing with globals is painful compared to the dependency parameter.

Force unit test coverage

## GlobalConfig I don't know what's secret and what isn't
The answer is nothing is secret but I can see fresh eyes freaking out about that.

Either create an empty `SecretsConfig` or an inner struct of the GlobalConfig.
If we go with these options I'd prefer the `SecretsConfig` to have complete logical separation.

I could just comment.

Coding a solution before it's needed could be a bad idea if my assumptions
about the future are wrong. I will just add a comment.

## Duplication
I'm finding I have to recreate basic parts of the http library in order to pass
a new parameter in.
Seems like a lot of duplicate work.
Is this a problem with the language or am I not thinking about this correctly?

The piece of code I'm most sad to give up is the multiplexer logic. I'm
beginning to give up that there's a way to maintain the multiplexer logic in
this regard.

Will need to be careful because all this should be unit testable

### Approaches
Create special root handler which builds all dependencies given what is present in the request.
A custom handler parses the URL and runs the correct handler.

This is a terrible idea because I have to now do routing. What if I want to use
something like `httprouter`?

Okay, what if we register each method of our root handler in main? The struct
is imbued with all dependencies and the methods will have access then.

Some dependencies are known at request and must also be different per request.
Modifying the struct has huge risks of leaking into other requests.

Hmmm. App contains bare minimum shared dependencies (mutexes, names of files via config)
All dependencies are built during request from these minimal shared dependencies.
Though it may be slow it might be the only possible solution.

I'm concerned it will be pretty bad for complexity and unit testing. It's
unknown the dependents for each endpoint except that it's a subset of the
struct.

What happens when you want to add an existing dependency somewhere else? Won't
that be a painful experience? Modifications to unit tests and potentially
changing how `buildDependencies` is going to work.

If use a single struct for all routes list of dependencies gets huge and unit tests become difficult.
It's probably our best bet tbh.
Google's Wire will help with dependencies. We create providers/constructors and
we'll manually marry everything together.
How will this system work for new developers? What would it take to create a
new service and hook it up to the existing dependencies? How will they register
this new service? Is there a chance they make a mistake?
They might but I don't think it will compile. Might be the best we got.

An interesting approach would be to force all services to always be request
scoped. The mutex can be moved to a goroutine that's handling the data race in
the background. I have thought about this before and it was decided that the it
was over-engineering the problem and probably slow.

## The use of "function builders"
I'm noticing frequent use of functions that produce new functions. I'm not sure how appropriate
or flexible such an approach is.

I'll try messing with things to see if an alternative approach is better.

## Does it make sense to make unit tests for main?
Main is designed to change very often. Unfortunately I've put a lot into main.

If I'm only willing to make unit tests for packages then the story becomes

Write in main -> want to validate code -> move code to library -> write unit tests

There's a lot of effort to make code testable this way. One of the benefits of unit testing is that
the code can demonstrably be isolated from the rest of the program.

Converting into a package also proves this point.

Another point is refactoring too early. I'd like to test code very early on.

All these things suggest I should be writing tests for different parts of main.

## Build constraints
I think this is a problem with Go. It doesn't let us explore all combinations
of build constraints at once in my IDE (not sure about VSCode or others).

You have to enable and disable them yourself so it's subject to human efforts.

## Fetching my prod data
I'm realising I don't publish my deployment scripts because the service is attempting to be
platform agnostic.

Fetching data from the platform definitionally isn't.
I'm not sure how much stock I want to put into this aspect now. Probably still
worth doing for troubleshooting sake.

Should focus more on generating interesting test cases.

## Generating test data
Should this be done via file or should it be done through the service itself?

If we do it via file
* it's much easier to script up.
* But the test generation only works for the database as a file.
* Changes to the file format require constant changes in test generation.
* Interesting tests can be made like migration or data recovery and created_at
  * Arguably it shouldn't cause problems. habit_share_file unit tests should handle interesting cases.

If we do it via API calls
* we have to change data generation when APIs change (That should technically be less often)
* Manual validation of the format needs to be done by the developer.
  * "Why can't we check via `GET`s? Subtleties store format might be important for the developer.
* Data generation is limited to what is possible with the APIs

-------------------------------------
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
* [ ] Remove need for Google Cloud OAuth. Need to rethink the API.
* [X] Sharing habits
* [ ] Document API
* [ ] All users are considered friends. Provide a friending system.
* [X] Import csv from HabitShare app
* [X] Create todo section

* [ ] Load testing to ensure data races don't pop up. Since only one instance of the app runs at a time the data races aren't as bad.
It's easier to solve these problems with an online database solution compared to a single file.
* [ ] Fix fragility of system on malformed data. Single file is parsed so any failures there brings the whole thing down.
production habits.json is 32K, for now we could do the swap method (maybe?)
* [ ] Same is true for the sessions.csv
