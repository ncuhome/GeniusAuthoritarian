namespace Admin {
  type LoginRecordAdminView = {
    id: number;
    createdAt: number;
    validBefore: number;
    uid: number;
    // unstable
    destroyed: boolean;
  };

  type DataFetchRange = "week" | "month" | "year";
}
