namespace Admin {
  type AppDataView = {
    id: number;
    name: string;
  };

  type LoginRecordDataView = {
    id: number;
    createdAt: number;
    validBefore: number;
    uid: number;
    // unstable
    destroyed: boolean;
  };

  type DataFetchRange = "week" | "month" | "year";

  type LoginDataView = {
    apps: AppDataView[];
    records: LoginRecordDataView[];
  };
}
