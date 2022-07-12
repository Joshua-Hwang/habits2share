The intention of this package is to import the csv provided by HabitShare to
produce Habits2Share habits with corresponding records.

When we import new habits we won't attempt to merge with existing habits for
reduced complexity.

The CSV relies on the name as the defining part of the habit. We will have a
mapping from name to habit id. If no mapping exists we create one and continue.

For each entry we register the date of the activity. There are four possible
entires

* `empty -> NOT_DONE`
* `fail -> NOT_DONE`
* `skip -> NOT_DONE`
* `success -> SUCCESS`

We will assume all completions were a success as an arbitrary choice (also lets
people keep their existing streaks).
