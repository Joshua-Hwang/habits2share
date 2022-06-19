import { Activity, Habit } from "./models";

export class Interactor {
  constructor(private baseUrl: string = "") {}

  public async getMyHabits(): Promise<Array<Habit>> {
    const res = await fetch(`${this.baseUrl}/my/habits`);
    return (await res.json()) as Array<Habit>;
  }

  public async getSharedHabits(): Promise<Array<Habit>> {
    const res = await fetch(`${this.baseUrl}/shared/habits`);
    return (await res.json()) as Array<Habit>;
  }

  public async getActivities(
    habitId: string
  ): Promise<{ Activities: Array<Activity>; HasMore: boolean }> {
    const res = await fetch(`${this.baseUrl}/habit/${habitId}/activities?after=2020-01-01`);
    console.log(res)
    return (await res.json()) as {
      Activities: Array<Activity>;
      HasMore: boolean;
    };
  }

  public async createHabit(name: string, frequency: number): Promise<string> {
    throw new Error("Not implemented yet");
  }
}
