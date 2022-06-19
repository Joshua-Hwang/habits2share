import { Card, Text } from "@mantine/core";
import { useEffect, useState } from "react";
import { Interactor } from "./interactors";
import { Activity, Habit } from "./models";

export function SevenDayDisplay({
  activities,
}: {
  activities: Array<Activity>;
}) {
  return <Text>{JSON.stringify(activities)}</Text>;
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
  const [activities, setActivities] = useState([] as Array<Activity>);

  useEffect(() => {
    interactor.getActivities(habit.Id).then((response) => {
      console.log(response);
      setActivities(response.Activities);
    });
  }, [habit, interactor]);

  return (
    <Card id={habit.Id}>
      <Text>{habit.Name}</Text>
      {showOwner && <Text>{habit.Owner}</Text>}
      <SevenDayDisplay activities={activities} />
    </Card>
  );
}
