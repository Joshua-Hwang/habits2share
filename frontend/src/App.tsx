import React, { useEffect, useState } from "react";
import "./App.css";
import {
  Button,
  Card,
  Center,
  Divider,
  Loader,
  MantineProvider,
  ScrollArea,
  SimpleGrid,
  Stack,
  Text,
} from "@mantine/core";
import { Interactor } from "./interactors";
import { Habit } from "./models";
import { HabitCard } from "./HabitCard";

function App({ interactor }: { interactor: Interactor }) {
  const [loadingMyHabits, setLoadingMyHabits] = useState(true);
  let [myHabits, setMyHabits] = useState([] as Array<Habit>);

  const [loadingSharedHabits, setLoadingSharedHabits] = useState(true);
  let [sharedHabits, setSharedHabits] = useState([] as Array<Habit>);

  useEffect(() => {
    interactor.getMyHabits().then((response) => {
      setMyHabits(response);
      setLoadingMyHabits(false);
    });
    interactor.getSharedHabits().then((response) => {
      setSharedHabits(response);
      setLoadingSharedHabits(false);
    });
  }, [interactor]);

  return (
    <MantineProvider
      theme={{
        // Override any other properties from default theme
        //fontFamily: "Open Sans, sans serif",
        spacing: { xs: 2, sm: 5, md: 7, lg: 10, xl: 20 },
      }}
    >
      <SimpleGrid className="App" cols={1} spacing={0}>
        <ScrollArea style={{ height: "90vh" }}>
          <Stack style={{ paddingInline: "1em" }}>
            {loadingMyHabits ? (
              <Center>
                <Loader />
              </Center>
            ) : (
              myHabits.map((habit) => {
                return (
                  <HabitCard
                    key={habit.Id}
                    habit={habit}
                    interactor={interactor}
                  />
                );
              })
            )}
            <Button>Create a new habit</Button>
            <Divider label="Shared habits below" labelPosition="center" />
            {loadingSharedHabits ? (
              <Center>
                <Loader />
              </Center>
            ) : (
              sharedHabits.map((habit) => {
                return (
                  <HabitCard
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
        <div style={{ height: "10vh" }}>Nav area</div>
      </SimpleGrid>
    </MantineProvider>
  );
}

export default App;
