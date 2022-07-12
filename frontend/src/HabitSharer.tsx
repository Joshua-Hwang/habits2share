import {
  ActionIcon,
  Button,
  Card,
  Group,
  Modal,
  ScrollArea,
  Stack,
  Text,
  TextInput,
} from "@mantine/core";
import { useState } from "react";
import { FaEyeSlash } from "react-icons/fa";
import { Habit } from "./models";

function ShareWithNewUser({
  onShare,
}: {
  onShare: (userId: string) => Promise<void>;
}) {
  const [newUserId, setNewUserId] = useState("");
  return (
    <Group noWrap>
      <TextInput
        placeholder="Share with"
        value={newUserId}
        onChange={(event) => setNewUserId(event.currentTarget.value)}
        style={{ width: "100%" }}
      />
      <Button
        onClick={async () => {
          await onShare?.(newUserId);
          setNewUserId('');
        }}
      >
        Press me
      </Button>
    </Group>
  );
}

export function HabitSharerModal({
  habit,
  onShare,
  onUnshare,
  opened,
  onClose,
}: {
  habit: Habit;
  onShare: (userId: string) => Promise<void>;
  onUnshare: (userId: string) => Promise<void>;
  opened: boolean;
  onClose: () => void;
}) {
  const userIds = [];
  for (let userId in habit.SharedWith) {
    userIds.push(userId);
  }
  return (
    <Modal
      title={`Share habit: "${habit.Name}"`}
      opened={opened}
      onClose={onClose}
    >
      <Stack>
        <Text>Shared with:</Text>
        <ScrollArea style={{ height: "33vh" }}>
          <Stack>
            {userIds.map((userId) => (
              <Card withBorder>
                <Group position="apart">
                  <Text>{userId}</Text>
                  <ActionIcon onClick={async () => await onUnshare?.(userId)}>
                    <FaEyeSlash />
                  </ActionIcon>
                </Group>
              </Card>
            ))}
          </Stack>
        </ScrollArea>
        <ShareWithNewUser onShare={onShare} />
      </Stack>
    </Modal>
  );
}
