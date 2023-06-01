declare namespace User {
  type Group = {
    id: number;
    name: string;
  };

  type LoginResult = {
    token: string;
    groups: string[];
  };
}
