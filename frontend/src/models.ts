// TODO generate this file from golang
// https://github.com/tkrajina/typescriptify-golang-structs

export type Habit = {
  Id: string;
  Owner: string;
  SharedWith: Record<string, {}>; // Due to nature of golang server
  Name: string;
  Frequency: number;
  Archived: boolean;
};

