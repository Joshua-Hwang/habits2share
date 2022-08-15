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
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { useRef, useState } from "react";

export function HabitCreatorModal({
  opened,
  onClose,
  onSubmit,
}: {
  opened: boolean;
  onClose: () => void;
  onSubmit: (values: { name: string; frequency: number }) => Promise<void>;
}) {
  const form = useForm({
    initialValues: {
      name: "",
      frequency: 7,
    },
  });
  const handlers = useRef<NumberInputHandlers>();
  const [loading, setLoading] = useState(false);

  return (
    <Modal title="New habit" size="xs" opened={opened} onClose={onClose}>
      <form
        onSubmit={form.onSubmit(async (values) => {
          setLoading(true);
          await onSubmit(values);
          form.reset();
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
          handlersRef={handlers}
          {...form.getInputProps("frequency")}
        />
        <Space h="lg" />
        <Group position="center">
          <Button loading={loading} disabled={loading} type="submit">
            Create
          </Button>
        </Group>
      </form>
    </Modal>
  );
}
