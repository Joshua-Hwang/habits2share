import {
  ActionIcon,
  Group,
  Modal,
  NumberInput,
  NumberInputHandlers,
  Stack,
  TextInput,
  Text,
  Button,
  Space,
  Checkbox,
  Center,
  Indicator,
  Box,
} from "@mantine/core";
import { Calendar } from "@mantine/dates";
import { useForm } from "@mantine/form";
import { useState } from "react";

function ActivityViewer() {
  return <Calendar />;
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
      frequency
    },
    // validate: {
    //   frequency: 
    // }
  });

  return (
    <form onSubmit={form.onSubmit((values) => onSubmit(values))}>
      <TextInput
        label="Habit name"
        required
        data-autofocus
        {...form.getInputProps('name')}
      />
      <Space h="lg" />
      <NumberInput
        label="Frequency"
        required
        min={1}
        max={7}
        {...form.getInputProps('frequency')}
      />
      <Space h="lg" />
      <Checkbox
        label="Confirm archive"
        checked={archiveConfirmed}
        onChange={(event) => setArchiveConfirmed(event.currentTarget.checked)}
      />
      <Space h="lg" />
      <Group position="center">
        <Button type="submit">
          Save
        </Button>
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
  name,
  frequency,
  opened,
  onClose,
  onSubmit,
  onArchive,
}: {
  name: string;
  frequency: number;
  opened: boolean;
  onClose: () => void;
  onSubmit: (values: { name: string; frequency: number }) => Promise<void>;
  onArchive: () => Promise<void>;
}) {
  return (
    <Modal
      title={`Edit habit: "${name}"`}
      size="xs"
      opened={opened}
      onClose={onClose}
    >
      <ActivityViewer />
      <HabitEditorForm
        name={name}
        frequency={frequency}
        onSubmit={onSubmit}
        onArchive={onArchive}
      />
    </Modal>
  );
}
