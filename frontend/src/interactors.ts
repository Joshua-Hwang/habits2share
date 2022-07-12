import dayjs, { Dayjs } from "dayjs";
import { Activity, Status, Habit } from "./models";

// TODO throw errors
// TODO use /habit/${habitId} endpoint which provides lots of information
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
    habitId: string,
    params?: {
      before?: string;
      after?: string;
      limit?: number;
    }
  ): Promise<{ Activities: Array<Activity>; HasMore: boolean }> {
    const { before, after, limit } = params ?? {};

    const queryParams = params
      ? new URLSearchParams({
          ...(before ? { before } : {}),
          ...(after ? { after } : {}),
          ...(limit ? { limit: limit.toFixed(0) } : {}), // TODO limit=0 is falsey but we don't expect limit=0 anyway
        }).toString()
      : "";
    const res = await fetch(
      `${this.baseUrl}/habit/${habitId}/activities${
        queryParams ? `?${queryParams}` : ""
      }`
    );

    const rawJson = (await res.json()) as {
      Activities: Array<{
        Id: string;
        HabitId: string;
        Logged: string;
        Status: string;
      }>;
      HasMore: boolean;
    };
    const parsedActivities = rawJson.Activities.map((activity) => ({
      ...activity,
      Logged: dayjs(activity.Logged),
    })) as Array<Activity>;

    return { Activities: parsedActivities, HasMore: rawJson.HasMore };
  }

  public async getScore(habitId: string): Promise<number> {
    const res = await fetch(`${this.baseUrl}/habit/${habitId}/score`);
    // TODO handle errors
    return parseInt(await res.text(), 10);
  }

  public async postActivity(
    habitId: string,
    logged: Dayjs,
    status: Status
  ): Promise<string> {
    const res = await fetch(`${this.baseUrl}/habit/${habitId}/activities`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        Logged: logged.format("YYYY-MM-DD"),
        Status: status,
      }),
    });
    return await res.text();
  }

  public async createHabit(name: string, frequency: number): Promise<string> {
    const res = await fetch(`${this.baseUrl}/my/habits`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ Name: name, Frequency: frequency }),
    });
    return await res.text();
  }

  public async updateHabit(
    habitId: string,
    name: string,
    frequency: number
  ): Promise<void> {
    await fetch(`${this.baseUrl}/habit/${habitId}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ Name: name, Frequency: frequency }),
    });
    // TODO throw if not okay
  }

  public async archiveHabit(habitId: string): Promise<void> {
    await fetch(`${this.baseUrl}/habit/${habitId}`, {
      method: "DELETE",
    });
  }

  public async shareHabit(habitId: string, userId: string): Promise<void> {
    await fetch(`${this.baseUrl}/user/${encodeURIComponent(userId)}/habit/${habitId}`, {
      method: "POST",
    });
  }

  public async unshareHabit(habitId: string, userId: string): Promise<void> {
    await fetch(`${this.baseUrl}/user/${encodeURIComponent(userId)}/habit/${habitId}`, {
      method: "DELETE",
    });
  }

  public async importHabits(csv: string): Promise<string[]> {
    const res = await fetch(`${this.baseUrl}/my/habits/upload`, {
      method: "POST",
      headers: {
        "Content-Type": "text/csv",
      },
      body: csv,
    });

    // TODO handle errors
    const habitIds = await res.text();
    return habitIds.split("\n");
  }
}
