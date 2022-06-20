import dayjs, { Dayjs } from "dayjs";
import { Activity, Status, Habit } from "./models";

// TODO throw errors
export class Interactor {
  constructor(private baseUrl: string = "") {}

  public async getMyHabits(): Promise<Array<Habit>> {
    await new Promise((resolve) => setTimeout(resolve, 1500));
    const res = await fetch(`${this.baseUrl}/my/habits`);
    return (await res.json()) as Array<Habit>;
  }

  public async getSharedHabits(): Promise<Array<Habit>> {
    await new Promise((resolve) => setTimeout(resolve, 2000));
    const res = await fetch(`${this.baseUrl}/shared/habits`);
    return (await res.json()) as Array<Habit>;
  }

  public async getActivities(
    habitId: string
  ): Promise<{ Activities: Array<Activity>; HasMore: boolean }> {
    await new Promise((resolve) => setTimeout(resolve, 1234));

    const res = await fetch(`${this.baseUrl}/habit/${habitId}/activities`);
    const rawJson = (await res.json()) as {
      Activities: Array<{Id: string, HabitId: string, Logged: string, Status: string}>;
      HasMore: boolean;
    };
    const parsedActivities = rawJson.Activities.map((activity) => ({
      ...activity,
      Logged: dayjs(activity.Logged),
    })) as Array<Activity>;

    return { Activities: parsedActivities, HasMore: rawJson.HasMore };
  }

  public async postActivity(
    habitId: string,
    logged: Dayjs,
    status: Status,
  ): Promise<string> {
    const res = await fetch(`${this.baseUrl}/habit/${habitId}/activities`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({Logged: logged.format("YYYY-MM-DD"), Status: status}),
    });
    return await res.text();
  }

  public async createHabit(name: string, frequency: number): Promise<string> {
    throw new Error("Not implemented yet");
  }
}
