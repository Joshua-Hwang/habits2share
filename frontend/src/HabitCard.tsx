import { Button, Card, Group, Loader, Text } from "@mantine/core";
import { useEffect, useReducer, useState } from "react";
import { Interactor } from "./interactors";
import { Activity, Habit, Status } from "./models";
import dayjs, { Dayjs } from "dayjs";

export function SevenDayDisplay({
  activities,
  onChange,
  numDays = 7,
}: {
  activities: Array<Activity>;
  onChange?: (date: Dayjs, newStatus: Status) => void;
  numDays?: number;
}) {
  const today = dayjs().startOf("day");
  const acts = activities
    .map((activity) => ({
      ...activity,
      Logged: dayjs(activity.Logged),
    }))
    .reverse(); // reverse is performed in-place hence it's put after mapping
  const allDays = [] as JSX.Element[];

  for (let d = 0, i = 0; d < numDays; d++) {
    const day = today.subtract(d, "days");
    let status = "NOT_DONE" as Status;
    if (i < acts.length && day.isSame(acts[i].Logged, "day")) {
      status = acts[i].Status;
      i++;
    }

    allDays.push(
      <Button
        compact
        size="xs"
        radius="xl"
        style={{ width: "3em" }}
        key={day.format("YYYY-MM-DD")}
        value={day.format("YYYY-MM-DD")}
        color={
          status === "SUCCESS"
            ? "green"
            : status === "MINIMUM"
            ? "blue"
            : "gray"
        }
        onClick={() =>
          onChange?.(
            day,
            status === "NOT_DONE"
              ? "SUCCESS"
              : status === "SUCCESS"
              ? "MINIMUM"
              : "NOT_DONE"
          )
        }
      >
        {day.format("dd")}
      </Button>
    );
  }

  return <Group spacing="xs">{allDays.reverse()}</Group>;
}

export function HabitCard({
  habit,
  interactor,
  showOwner,
}: {
  habit: Habit;
  interactor: Interactor;
  showOwner?: boolean;
}) {
  const [loadingDates, setLoadingDates] = useState(true);

  const [activities, dispatch] = useReducer(
    (
      arr: Array<Activity>,
      input:
        | {
            action: "overwrite";
            array: Array<Activity>;
          }
        | {
            action: "replace";
            index: number;
            value: Activity;
          }
        | {
            action: "insert";
            index: number; // element at i will be pushed to i+1
            value: Activity;
          }
        | {
            action: "append";
            value: Activity;
          }
    ) => {
      // NOTE: react performs a comparison to check if an update is necessary
      // thus in-place algorithms aren't okay
      switch (input.action) {
        case "overwrite":
          return input.array;
        case "replace":
          const x = arr.map((activity, i) =>
            i === input.index ? input.value : activity
          );
          return x;
        case "insert":
          return [
            ...arr.slice(0, input.index),
            input.value,
            ...arr.slice(input.index),
          ];
        default:
          return [...arr, input.value];
      }
    },
    [] as Array<Activity>
  );

  useEffect(() => {
    interactor.getActivities(habit.Id).then((response) => {
      // NOTE activities are ordered oldest to latest
      dispatch({ action: "overwrite", array: response.Activities });
      setLoadingDates(false);
    });
  }, [habit, interactor]);

  return (
    <Card key={habit.Id}>
      <Text>{habit.Name}</Text>
      {showOwner && <Text>{habit.Owner}</Text>}
      {loadingDates ? (
        <Loader />
      ) : (
        <SevenDayDisplay
          numDays={7}
          activities={activities}
          onChange={(day, status) => {
            interactor.postActivity(habit.Id, day, status).then((_) => {
              // make a notification
              console.log("activity posted", day, status);
            });

            // TODO algorithm could definitely be improved but it's only 7 entries
            let i = 0;
            for (; i < activities.length; i++) {
              // if the day is equal then change the existing entry then return
              if (day.isSame(activities[i].Logged, "day")) {
                dispatch({
                  action: "replace",
                  index: i,
                  value: { ...activities[i], Status: status },
                });
                return;
              }
              // if the day is after activities[i] then we need to splice to insert then return
              if (day.isBefore(activities[i].Logged)) {
                dispatch({
                  action: "insert",
                  index: i,
                  value: {
                    Id: "local",
                    HabitId: habit.Id,
                    Logged: day,
                    Status: status,
                  },
                });
                return;
              }
            }
            // if we're here then the day is after all activities so push to end and return
            dispatch({
              action: "append",
              value: {
                Id: "local",
                HabitId: habit.Id,
                Logged: day,
                Status: status,
              },
            });
          }}
        />
      )}
    </Card>
  );
}
