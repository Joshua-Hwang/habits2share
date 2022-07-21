// TODO generate this file from golang
// https://github.com/tkrajina/typescriptify-golang-structs

import { Dayjs } from "dayjs";

export type Habit = {
  Id: string;
  Owner: string;
  SharedWith: Record<string, {}>; // Due to nature of golang server
  Name: string;
  Frequency: number;
  Description: string;
  Archived: boolean;
};

export type Status = "SUCCESS" | "MINIMUM" | "NOT_DONE";

export type Activity = {
  Id: string;
  HabitId: string;
  Logged: Dayjs;
  Status: Status;
};
