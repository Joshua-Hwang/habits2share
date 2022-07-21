import dayjs from "dayjs";
import { Interactor } from "./interactors";
import { Habit, Activity } from "./models";

type PublicPart<T> = { [K in keyof T]: T[K] };

export class MockInteractor /*implements PublicPart<Interactor>*/ {
  public async getMyHabits(): Promise<Array<Habit>> {
    await new Promise((resolve) => setTimeout(resolve, 2000));
    return [
      {
        Id: "1234",
        Owner: "hello",
        SharedWith: { wilson: {} },
        Name: "adfdsafafds",
        Frequency: 3,
        Description: "",
        Archived: false,
      },
      {
        Id: "1235",
        Owner: "hello",
        SharedWith: { wilson: {} },
        Name: "what",
        Frequency: 3,
        Description: "",
        Archived: false,
      },
      {
        Id: "1236",
        Owner: "hello",
        SharedWith: { wilson: {} },
        Name: "new",
        Frequency: 3,
        Description: "",
        Archived: false,
      },
      {
        Id: "1233",
        Owner: "hello",
        SharedWith: { wilson: {} },
        Name: "noted",
        Frequency: 3,
        Description: "",
        Archived: false,
      },
      {
        Id: "1239",
        Owner: "hello",
        SharedWith: { wilson: {} },
        Name: "let's go",
        Frequency: 3,
        Description: "",
        Archived: false,
      },
    ];
  }

  public async getSharedHabits(): Promise<Array<Habit>> {
    await new Promise((resolve) => setTimeout(resolve, 1000));
    return [
      {
        Id: "999",
        Owner: "george",
        SharedWith: { wilson: {} },
        Name: "big",
        Frequency: 3,
        Description: "",
        Archived: false,
      },
      {
        Id: "3245",
        Owner: "marcus",
        SharedWith: { wilson: {} },
        Name: "names",
        Frequency: 3,
        Description: "",
        Archived: false,
      },
      {
        Id: "4829",
        Owner: "stanley",
        SharedWith: { wilson: {} },
        Name: "matter",
        Frequency: 3,
        Description: "",
        Archived: false,
      },
    ];
  }

  public async getActivities(
    habitId: string
  ): Promise<{ Activities: Array<Activity>; HasMore: boolean }> {
    await new Promise((resolve) => setTimeout(resolve, 1000));
    if (habitId === "1234") {
      return {
        Activities: [
          {
            Id: "dfa",
            HabitId: "1234",
            Logged: dayjs("2022-06-20T00:00:00Z"),
            Status: "SUCCESS",
          },
        ],
        HasMore: false,
      };
    }
    return {
      Activities: [],
      HasMore: false,
    };
  }

  public async createHabit(name: string, frequency: number): Promise<string> {
    return "gferea";
  }
}
