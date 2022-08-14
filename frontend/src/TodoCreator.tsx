import { Button, Group, Modal, Space, TextInput } from "@mantine/core";
import { DatePicker } from "@mantine/dates";
import { useForm } from "@mantine/form";
import dayjs from "dayjs";
import { useState } from "react";

export function TodoCreatorModal({
  opened,
  onClose,
  onSubmit,
}: {
  opened: boolean;
  onClose: () => void;
  onSubmit: (values: {
    name: string;
    dueDate: dayjs.Dayjs;
    description: string;
  }) => Promise<void>;
}) {
  const form = useForm({
    initialValues: {
      name: "",
      dueDate: new Date(),
      description: "",
    },
  });
  const [loading, setLoading] = useState(false);

  return (
    <Modal title="New todo" size="xs" opened={opened} onClose={onClose}>
      <form
        onSubmit={form.onSubmit(async (values) => {
          setLoading(true);
          await onSubmit({ ...values, dueDate: dayjs(values.dueDate) });
          setLoading(false);
        })}
      >
        <TextInput
          label="Todo name"
          required
          data-autofocus
          {...form.getInputProps("name")}
        />
        <Space h="lg" />
        <DatePicker
          label="Due date"
          required
          {...form.getInputProps("dueDate")}
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
