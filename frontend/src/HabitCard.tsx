import {
  ActionIcon,
  Button,
  Card,
  Code,
  Divider,
  Group,
  Loader,
  Space,
  Text,
} from "@mantine/core";
import { useEffect, useReducer, useState } from "react";
import { Interactor } from "./interactors";
import { Activity, Habit, Status } from "./models";
import dayjs, { Dayjs } from "dayjs";
import { FaEdit, FaShareSquare } from "react-icons/fa";
import { HabitEditorModal } from "./HabitEditor";
import { HabitSharerModal } from "./HabitSharer";

export function SevenDayDisplay({
  activities,
  onChange,
  numDays = 7,
  disabled,
  start = "today",
}: {
  activities: Array<Activity>;
  onChange?: (date: Dayjs, newStatus: Status) => void;
  numDays?: number;
  disabled?: boolean;
  start?: "today" | "monday";
}) {
  // firstDay to be displayed (Sunday of that week)
  // startOf("week") considers Sunday to be the start
  // subtract 1 will round Sunday to the previous Sunday instead
  const firstDay =
    start === "today"
      ? dayjs().startOf("day")
      : dayjs().subtract(1, "day").startOf("week").add(7, "days");

  // activities first index is earliest
  activities = activities.slice();
  const allDays = [] as JSX.Element[];

  for (let d = 0, i = activities.length - 1; d < numDays; d++) {
    const day = firstDay.subtract(d, "days");
    let status = "NOT_DONE" as Status;
    if (i >= 0 && day.isSame(activities[i].Logged, "day")) {
      status = activities[i].Status;
      i--;
    }

    allDays.push(
      <Button
        compact
        size="xs"
        radius="xl"
        style={{ width: "3em" }}
        key={day.format("YYYY-MM-DD")}
        value={day.format("YYYY-MM-DD")}
        variant="default"
        sx={(theme) => {
          return {
            backgroundColor:
              status === "SUCCESS"
                ? theme.colors.green[5]
                : status === "MINIMUM"
                ? theme.colors.blue[3]
                : theme.colors.gray[3],
            ":hover": {
              backgroundColor:
                status === "SUCCESS"
                  ? theme.colors.green[5]
                  : status === "MINIMUM"
                  ? theme.colors.blue[3]
                  : theme.colors.gray[3],
            },
          };
        }}
        onClick={() => {
          if (disabled) {
            return;
          }

          return onChange?.(
            day,
            status === "NOT_DONE"
              ? "SUCCESS"
              : status === "SUCCESS"
              ? "MINIMUM"
              : "NOT_DONE"
          );
        }}
      >
        {day.format("dd")}
      </Button>
    );
  }

  return <Group spacing="xs">{allDays.reverse()}</Group>;
}

export function HabitCard({
  habit,
  setHabit,
  onArchive,
  interactor,
  showOwner,
  disabled,
}: {
  habit: Habit;
  setHabit?: (
    habit: Habit,
    args?: { dontUpdateServer: boolean }
  ) => Promise<void>;
  onArchive?: () => Promise<void>;
  interactor: Interactor;
  showOwner?: boolean;
  disabled?: boolean;
}) {
  const DEBUG = false;
  const [editModalOpened, setEditModalOpened] = useState(false);
  const [shareModalOpened, setShareModalOpened] = useState(false);

  const [loadingDates, setLoadingDates] = useState(true);

  const [score, setScore] = useState(-1); // -1 is sentinel value
  const [activitiesThisWeek, setActivitiesThisWeek] = useState(-1);

  const [activities, activitiesDispatch] = useReducer(
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
    // TODO this could be handled with a single call to /habit/:habitId
    interactor
      .getActivities(habit.Id, {
        after: dayjs().subtract(7, "days").format("YYYY-MM-DD"),
        limit: 7,
      })
      .then(({ Activities }) => {
        // NOTE activities are ordered ascending on unix time
        activitiesDispatch({ action: "overwrite", array: Activities });

        // start of week is technically Sunday but I think it's Monday
        const startOfWeek = dayjs().startOf("week").add(1, "day");
        let activitiesThisWeek = Activities.length;

        for (let i = Activities.length - 1; i >= 0; i--) {
          if (Activities[i].Logged.isBefore(startOfWeek)) {
            activitiesThisWeek -= i + 1;
            break;
          }
        }
        setActivitiesThisWeek(activitiesThisWeek);
        setLoadingDates(false);
      }); // TODO handle catch

    interactor.getScore(habit.Id).then((score) => {
      setScore(score);
    });
  }, [habit.Id, interactor]);

  // number done this week
  // TODO there is a habit editor modal for each card. This could be a performance hit
  // Create one global habit editor modal and each card updates that
  return (
    <Card key={habit.Id}>
      <HabitSharerModal
        habit={habit}
        onShare={async (userId) => {
          await setHabit?.(
            { ...habit, SharedWith: { ...habit.SharedWith, [userId]: {} } },
            { dontUpdateServer: true } // TODO improve name from dontUpdateServer
          );
          await interactor.shareHabit(habit.Id, userId);
        }}
        onUnshare={async (userId) => {
          const { [userId]: _userId, ...newSharedWith } = habit.SharedWith;
          await setHabit?.(
            { ...habit, SharedWith: newSharedWith },
            { dontUpdateServer: true }
          );
          await interactor.unshareHabit(habit.Id, userId);
        }}
        opened={shareModalOpened}
        onClose={() => setShareModalOpened(false)}
      />
      <HabitEditorModal
        habit={habit}
        interactor={interactor}
        onSubmit={async ({ name, frequency, description }) => {
          await setHabit?.({
            ...habit,
            Name: name,
            Frequency: frequency,
            Description: description,
          });
          setEditModalOpened(false);
        }}
        opened={editModalOpened}
        onClose={() => setEditModalOpened(false)}
        onArchive={async () => {
          await onArchive?.();
          setEditModalOpened(false);
        }}
      />
      <Group position="apart">
        <Text>{habit.Name}</Text>
        {setHabit && (
          <Group>
            <ActionIcon
              onClick={() => {
                setShareModalOpened(true);
              }}
            >
              <FaShareSquare />
            </ActionIcon>
            <ActionIcon
              onClick={() => {
                setEditModalOpened(true);
              }}
            >
              <FaEdit />
            </ActionIcon>
          </Group>
        )}
      </Group>
      <Group>
        {showOwner && (
          <>
            <Text size="xs">owner: {habit.Owner}</Text>
            <Divider sx={{ height: "auto" }} orientation="vertical" />
          </>
        )}
        <Text size="xs">
          This Week: {activitiesThisWeek === -1 ? "..." : activitiesThisWeek}/
          {habit.Frequency}
        </Text>
        <Divider style={{ height: "auto" }} orientation="vertical" />
        <Text size="xs">Score: {score === -1 ? "..." : score}</Text>
        {Object.keys(habit.SharedWith).length > 0 && (
          <>
            <Divider style={{ height: "auto" }} orientation="vertical" />
            <Text size="xs">
              Shared with: {Object.keys(habit.SharedWith).length}
            </Text>
          </>
        )}
        {DEBUG && (
          <>
            <Divider sx={{ height: "auto" }} orientation="vertical" />
            <Text size="xs">
              <Code>{habit.Id}</Code>
            </Text>
          </>
        )}
      </Group>
      <Group>
        <Divider sx={{ height: "auto" }} orientation="vertical" />
        <Text style={{ whiteSpace: "pre-line" }} size="xs">
          {habit.Description}
        </Text>
      </Group>
      <Space h="md" />
      {loadingDates ? (
        <Loader />
      ) : (
        <SevenDayDisplay
          start="monday"
          disabled={disabled}
          numDays={7}
          activities={activities}
          onChange={(day, status) => {
            setScore(-1); // reset as we retrieve new score
            interactor.postActivity(habit.Id, day, status).then((_) => {
              interactor.getScore(habit.Id).then(setScore);
              // TODO make a notification
              console.log("activity posted", day, status);
            });

            // TODO algorithm could definitely be improved but it's only 7 entries
            for (let i = 0; i < activities.length; i++) {
              // if the day is equal then change the existing entry then return
              if (day.isSame(activities[i].Logged, "day")) {
                if (activities[i].Status !== status) {
                  activitiesDispatch({
                    action: "replace",
                    index: i,
                    value: { ...activities[i], Status: status },
                  });
                  if (status === "NOT_DONE") {
                    setActivitiesThisWeek(activitiesThisWeek - 1);
                  }
                  if (activities[i].Status === "NOT_DONE") {
                    // changing from NOT_DONE to something useful
                    setActivitiesThisWeek(activitiesThisWeek + 1);
                  }
                }
                return;
              }
              // if the day is after activities[i] then we need to splice to insert then return
              if (day.isBefore(activities[i].Logged)) {
                activitiesDispatch({
                  action: "insert",
                  index: i,
                  value: {
                    Id: "local",
                    HabitId: habit.Id,
                    Logged: day,
                    Status: status,
                  },
                });
                if (status !== "NOT_DONE") {
                  setActivitiesThisWeek(activitiesThisWeek + 1);
                }
                return;
              }
            }
            // if we're here then the day is after all activities so push to end and return
            activitiesDispatch({
              action: "append",
              value: {
                Id: "local",
                HabitId: habit.Id,
                Logged: day,
                Status: status,
              },
            });
            if (status !== "NOT_DONE") {
              setActivitiesThisWeek(activitiesThisWeek + 1);
            }
          }}
        />
      )}
    </Card>
  );
}
