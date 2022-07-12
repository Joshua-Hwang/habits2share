import { Group, Modal, Text, useMantineTheme } from "@mantine/core";
import { Dropzone } from "@mantine/dropzone";
import { FaUpload } from "react-icons/fa";

export function HabitImporterModal({
  opened,
  onClose,
  onSubmit,
}: {
  opened: boolean;
  onClose: () => void;
  onSubmit: (csv: string) => Promise<void>;
}) {
  const theme = useMantineTheme();

  // TODO make this limit an environment variable
  return (
    <Modal title="Import habits" size="xs" opened={opened} onClose={onClose}>
      <Dropzone
        multiple={false}
        maxSize={3 * 1024 ** 2}
        onDrop={async ([file, ..._]) => {
          const csv = await file.text();
          await onSubmit(csv);
          onClose();
        }}
      >
        {() => (
          <Group
            position="center"
            spacing="xl"
            style={{ minHeight: 220, pointerEvents: "none" }}
          >
            <FaUpload color={theme.colors.gray[3]} size={80} />

            <div>
              <Text size="xl" inline>
                Drag CSV files here or click to select files
              </Text>
              <Text size="sm" color="dimmed" inline mt={7}>
                This file usually comes from the HabitShare app. There is a 5mb
                limit. Talk to the hoster to change this limit.
              </Text>
            </div>
          </Group>
        )}
      </Dropzone>
    </Modal>
  );
}
