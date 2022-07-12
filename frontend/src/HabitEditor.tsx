import {
  Group,
  Modal,
  NumberInput,
  TextInput,
  Button,
  Space,
  Checkbox,
  Center,
  Text,
  Box,
  Paper,
  useMantineTheme,
} from "@mantine/core";
import { Calendar } from "@mantine/dates";
import { useForm } from "@mantine/form";
import dayjs from "dayjs";
import { useEffect, useReducer, useState } from "react";
import { Interactor } from "./interactors";
import { Activity, Habit, Status } from "./models";

const DATE_FORMAT = "YYYY-MM-DD";

function ActivityViewer({
  habit,
  interactor,
}: {
  habit: Habit;
  interactor: Interactor;
}) {
  // TODO handle loading better
  const [month, onMonthChange] = useState(dayjs().startOf("month"));
  const [activities, dispatch] = useReducer(
    (
      map: Record<string, Activity | undefined>,
      input: {
        action: "insert";
        key: string; // element at i will be pushed to i+1
        value: Activity;
      }
    ) => {
      // NOTE: react performs a comparison to check if an update is necessary
      // thus in-place algorithms aren't okay
      switch (input.action) {
        case "insert":
          return { ...map, [input.key]: input.value };
      }
    },
    {} as Record<string, Activity>
  );

  useEffect(() => {
    // TODO catch errors
    interactor
      .getActivities(habit.Id, {
        after: month.subtract(1, "month").format(DATE_FORMAT),
        before: month.add(2, "month").format(DATE_FORMAT),
        limit: 100,
      })
      .then(({ Activities: activitiesList }) => {
        activitiesList.forEach((activity) => {
          dispatch({
            action: "insert",
            key: activity.Logged.format(DATE_FORMAT),
            value: activity,
          });
        });
      });
  }, [month]);

  const theme = useMantineTheme();
  return (
    <Calendar
      onChange={(date) => {
        if (!date) {
          return;
        }
        const d = dayjs(date);
        const dateString = d.format(DATE_FORMAT);
        const activity = activities[dateString];

        const status: Status = !activity
          ? "SUCCESS"
          : activity.Status === "SUCCESS"
          ? "MINIMUM"
          : "NOT_DONE";

        dispatch({
          action: "insert",
          key: dateString,
          value: {
            Id: "local",
            HabitId: habit.Id,
            Logged: d,
            Status: status,
          },
        });
        interactor.postActivity(habit.Id, d, status);
      }}
      month={month.toDate()}
      onMonthChange={(month) => onMonthChange(dayjs(month))}
      dayStyle={(date, { outside }) => {
        const dateString = dayjs(date).format(DATE_FORMAT);
        const activity = activities[dateString];
        return {
          color: theme.colors.gray[8],
          ...(!activity || activity.Status === "NOT_DONE"
            ? {}
            : {
                backgroundColor:
                  activity.Status === "SUCCESS"
                    ? theme.colors.green[outside ? 1 : 5]
                    : theme.colors.blue[outside ? 1 : 3],
              }),
        };
      }}
    />
  );
}

function HabitEditorForm({
  name,
  frequency,
  onSubmit,
  onArchive,
}: {
  name: string;
  frequency: number;
  onSubmit: (values: { name: string; frequency: number }) => Promise<void>;
  onArchive: () => Promise<void>;
}) {
  const [loading, setLoading] = useState(false);
  const [archiveConfirmed, setArchiveConfirmed] = useState(false);

  const form = useForm({
    initialValues: {
      name,
      frequency,
    },
  });

  return (
    <form
      onSubmit={form.onSubmit((values) => {
        setLoading(true);
        onSubmit(values);
        setLoading(false);
      })}
    >
      <TextInput
        label="Habit name"
        required
        data-autofocus
        {...form.getInputProps("name")}
      />
      <Space h="lg" />
      <NumberInput
        label="Frequency"
        required
        min={1}
        max={7}
        {...form.getInputProps("frequency")}
      />
      <Space h="lg" />
      <Checkbox
        label="Confirm archive"
        checked={archiveConfirmed}
        onChange={(event) => setArchiveConfirmed(event.currentTarget.checked)}
      />
      <Space h="lg" />
      <Group position="center">
        <Button type="submit">Save</Button>
        {archiveConfirmed && (
          <Button
            color="red"
            variant="outline"
            loading={loading}
            disabled={loading || !archiveConfirmed}
            onClick={onArchive}
          >
            Archive
          </Button>
        )}
      </Group>
    </form>
  );
}

export function HabitEditorModal({
  habit,
  interactor,
  opened,
  onClose,
  onSubmit,
  onArchive,
}: {
  habit: Habit;
  interactor: Interactor;
  opened: boolean;
  onClose: () => void;
  onSubmit: (values: { name: string; frequency: number }) => Promise<void>;
  onArchive: () => Promise<void>;
}) {
  return (
    <Modal
      title={`Edit habit: "${habit.Name}"`}
      size="xs"
      opened={opened}
      onClose={onClose}
      styles={(theme) => ({
        modal: {
          backgroundColor: theme.colors.gray[2],
        },
      })}
    >
      <Paper style={{ width: "100%" }} withBorder p="md">
        <Center>
          <ActivityViewer habit={habit} interactor={interactor} />
        </Center>
      </Paper>
      <Space h="sm" />
      <Paper style={{ width: "100%" }} withBorder p="md">
        <HabitEditorForm
          name={habit.Name}
          frequency={habit.Frequency}
          onSubmit={onSubmit}
          onArchive={onArchive}
        />
      </Paper>
    </Modal>
  );
}
