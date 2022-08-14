import { useEffect, useReducer, useState } from "react";
import {
  ActionIcon,
  Box,
  Button,
  Card,
  Center,
  Divider,
  Grid,
  Group,
  Loader,
  MantineProvider,
  ScrollArea,
  SimpleGrid,
  Stack,
  Text,
} from "@mantine/core";
import { Interactor } from "./interactors";
import { Habit, Todo } from "./models";
import { HabitCard } from "./HabitCard";
import { HabitCreatorModal } from "./HabitCreator";
import { FaCheckCircle, FaFileImport, FaPlusCircle } from "react-icons/fa";
import { BinarySearch } from "./util/BinarySearch";
import { HabitImporterModal } from "./HabitImporter";
import { TodoCreatorModal } from "./TodoCreator";
import dayjs from "dayjs";

function HabitTab({
  interactor,
  hidden = false,
}: {
  interactor: Interactor;
  hidden?: boolean;
}) {
  const [newHabitModalOpened, setNewHabitModalOpened] = useState(false);
  const [importHabitModalOpened, setImportHabitModalOpened] = useState(false);

  const [loadingMyHabits, setLoadingMyHabits] = useState(true);
  const [myHabits, dispatchMyHabits] = useReducer(
    (
      arr: Array<Habit>,
      input:
        | { action: "overwrite"; arr: Array<Habit> }
        | { action: "insert"; value: Habit }
        | { action: "remove"; index: number }
        | { action: "replace"; index: number; value: Habit }
    ) => {
      switch (input.action) {
        case "overwrite": {
          return input.arr;
        }
        case "insert": {
          const index = BinarySearch(
            (i) => input.value.Name < arr[i].Name,
            arr.length
          );

          return [...arr.slice(0, index), input.value, ...arr.slice(index)];
        }
        case "remove": {
          return [...arr.slice(0, input.index), ...arr.slice(input.index + 1)];
        }
        case "replace": {
          // TODO could be made cleaner/faster in future but hardly a bottleneck
          const removed = [
            ...arr.slice(0, input.index),
            ...arr.slice(input.index + 1),
          ];
          const newIndex = BinarySearch(
            (i) => input.value.Name < removed[i].Name,
            removed.length
          );
          return [
            ...removed.slice(0, newIndex),
            input.value,
            ...removed.slice(newIndex),
          ];
        }
      }
    },
    [] as Array<Habit>
  );

  const [loadingSharedHabits, setLoadingSharedHabits] = useState(true);
  // shared habits is a simple state because the user can do little to modify it
  const [sharedHabits, setSharedHabits] = useState([] as Array<Habit>);

  useEffect(() => {
    interactor.getMyHabits().then((response) => {
      dispatchMyHabits({ action: "overwrite", arr: response });
      setLoadingMyHabits(false);
    });
    interactor.getSharedHabits().then((response) => {
      setSharedHabits(response);
      setLoadingSharedHabits(false);
    });
  }, [interactor]);

  return (
    <>
      <HabitImporterModal
        onSubmit={async (csv) => {
          // We could the habit ids and fetch each habit
          // We could also have the upload endpoint return all habits
          // instead it's easier to grab all the habits again
          await interactor.importHabits(csv);

          const arr = await interactor.getMyHabits();
          dispatchMyHabits({ action: "overwrite", arr });
        }}
        opened={importHabitModalOpened}
        onClose={() => setImportHabitModalOpened(false)}
      />
      <HabitCreatorModal
        onSubmit={async ({ name, frequency }) => {
          setLoadingMyHabits(true);

          const habitId = await interactor.createHabit(name, frequency);

          dispatchMyHabits({
            action: "insert",
            value: {
              Id: habitId,
              Name: name,
              Frequency: frequency,
              Owner: "self",
              SharedWith: {},
              Description: "",
              Archived: false,
            },
          });

          setLoadingMyHabits(false);
          setNewHabitModalOpened(false);
        }}
        opened={newHabitModalOpened}
        onClose={() => setNewHabitModalOpened(false)}
      />
      <ScrollArea hidden={hidden}>
        <Stack style={{ padding: "1em" }}>
          <Grid columns={6} grow>
            <Grid.Col span={1}>
              <Button
                fullWidth
                leftIcon={<FaFileImport />}
                color="green"
                onClick={() => setImportHabitModalOpened(true)}
              >
                Import CSV
              </Button>
            </Grid.Col>
            <Grid.Col span={5}>
              <Button
                fullWidth
                leftIcon={<FaPlusCircle />}
                color="green"
                onClick={() => setNewHabitModalOpened(true)}
              >
                Create a new habit
              </Button>
            </Grid.Col>
          </Grid>
          {loadingMyHabits ? (
            <Center>
              <Loader />
            </Center>
          ) : (
            myHabits.map((habit, index) => {
              return (
                <HabitCard
                  key={habit.Id}
                  habit={habit}
                  setHabit={async (habit, args) => {
                    // TODO handle errors
                    if (!args?.dontUpdateServer) {
                      await interactor.updateHabit(
                        habit.Id,
                        habit.Name,
                        habit.Frequency,
                        habit.Description
                      );
                    }
                    dispatchMyHabits({
                      action: "replace",
                      index,
                      value: habit,
                    });
                  }}
                  onArchive={async () => {
                    await interactor.archiveHabit(habit.Id);
                    dispatchMyHabits({
                      action: "remove",
                      index,
                    });
                  }}
                  interactor={interactor}
                />
              );
            })
          )}
          <Divider label="Shared habits below" labelPosition="center" />
          {loadingSharedHabits ? (
            <Center>
              <Loader />
            </Center>
          ) : (
            sharedHabits.map((habit) => {
              return (
                <HabitCard
                  disabled
                  key={habit.Id}
                  habit={habit}
                  interactor={interactor}
                  showOwner={true}
                />
              );
            })
          )}
        </Stack>
      </ScrollArea>
    </>
  );
}

function TodoTab({
  interactor,
  hidden = false,
}: {
  interactor: Interactor;
  hidden?: boolean;
}) {
  const [newTodoModalOpened, setNewTodoModalOpened] = useState(false);

  const [loadingMyTodos, setLoadingMyTodos] = useState(true);
  const [myTodos, dispatchMyTodos] = useReducer(
    (
      arr: Array<Todo>,
      input:
        | { action: "overwrite"; arr: Array<Todo> }
        | { action: "insert"; value: Todo }
        | { action: "remove"; index: number }
        | { action: "replace"; index: number; value: Todo }
    ) => {
      switch (input.action) {
        case "overwrite": {
          return input.arr;
        }
        case "insert": {
          const index = BinarySearch(
            (i) => input.value.Name < arr[i].Name,
            arr.length
          );

          return [...arr.slice(0, index), input.value, ...arr.slice(index)];
        }
        case "remove": {
          return [...arr.slice(0, input.index), ...arr.slice(input.index + 1)];
        }
        case "replace": {
          // TODO could be made cleaner/faster in future but hardly a bottleneck
          const removed = [
            ...arr.slice(0, input.index),
            ...arr.slice(input.index + 1),
          ];
          const newIndex = BinarySearch(
            (i) => input.value.Name < removed[i].Name,
            removed.length
          );
          return [
            ...removed.slice(0, newIndex),
            input.value,
            ...removed.slice(newIndex),
          ];
        }
      }
    },
    [] as Array<Todo>
  );

  useEffect(() => {
    interactor.getMyTodos().then((response) => {
      dispatchMyTodos({ action: "overwrite", arr: response });
      setLoadingMyTodos(false);
    });
  }, [interactor]);

  return (
    <>
      <TodoCreatorModal
        onSubmit={async ({ name, description, dueDate }) => {
          setLoadingMyTodos(true);

          const todoId = await interactor.createTodo(
            name,
            description,
            dueDate
          );

          dispatchMyTodos({
            action: "insert",
            value: {
              Id: todoId,
              Name: name,
              Description: description,
              Owner: "self",
              DueDate: dueDate,
              Completed: false,
            },
          });

          setLoadingMyTodos(false);
          setNewTodoModalOpened(false);
        }}
        opened={newTodoModalOpened}
        onClose={() => setNewTodoModalOpened(false)}
      />
      <ScrollArea hidden={hidden}>
        <Stack style={{ padding: "1em" }}>
          <Button
            fullWidth
            leftIcon={<FaPlusCircle />}
            color="green"
            onClick={() => setNewTodoModalOpened(true)}
          >
            Create a new todo
          </Button>
          {loadingMyTodos ? (
            <Center>
              <Loader />
            </Center>
          ) : (
            myTodos.map((todo, index) => {
              return (
                <Card key={todo.Id}>
                  <Group position="apart">
                    <Text>{todo.Name}</Text>
                    <ActionIcon
                      onClick={() => {
                        interactor.updateTodo(
                          todo.Id,
                          todo.Name,
                          todo.Description,
                          todo.DueDate,
                          true
                        );
                        dispatchMyTodos({ action: "remove", index });
                      }}
                    >
                      <FaCheckCircle />
                    </ActionIcon>
                  </Group>
                  <Group>
                    <Text
                      size="xs"
                      color={todo.DueDate.isBefore(dayjs()) ? "red" : "black"}
                    >
                      {
                        // TODO think of better format
                      }
                      Due date: {todo.DueDate.format()}
                    </Text>
                  </Group>
                  <Group>
                    <Divider sx={{ height: "auto" }} orientation="vertical" />
                    <Text style={{ whiteSpace: "pre-line" }} size="xs">
                      {todo.Description}
                    </Text>
                  </Group>
                </Card>
              );
            })
          )}
        </Stack>
      </ScrollArea>
    </>
  );
}

function App({ interactor }: { interactor: Interactor }) {
  const [currentTab, setCurrentTab] = useState("habits" as "habits" | "todos");
  return (
    <MantineProvider
      theme={{
        colorScheme: "light",
        // Override any other properties from default theme
        //fontFamily: "Open Sans, sans serif",
        spacing: { xs: 2, sm: 5, md: 7, lg: 10, xl: 20 },
      }}
    >
      <SimpleGrid
        sx={(theme) => ({
          gridTemplateRows: "1fr 3em",
          background: theme.colors.gray[3],
          height: "100vh",
        })}
        cols={1}
        spacing={0}
      >
        <HabitTab hidden={currentTab !== "habits"} interactor={interactor} />
        <TodoTab hidden={currentTab !== "todos"} interactor={interactor} />
        <Group position="apart" sx={{ paddingInline: "1rem" }}>
          <Button onClick={() => setCurrentTab("habits")}>Habits</Button>
          <Button onClick={() => setCurrentTab("todos")}>Checklist</Button>
        </Group>
      </SimpleGrid>
    </MantineProvider>
  );
}

export default App;
