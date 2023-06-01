declare namespace User {
  type Group = {
    id: number;
    name: string;
  };

  type LoginResult = {
    token: string;
    groups: string[];
  };

  type Profile = {
    user: {
      id: number;
      name: string;
      phone: string;
      groups: Group[];
    };
    loginRecord: Array<{
      id: number;
      createdAt: number;
      target: string;
      ip: string;
    }>;
  };
}
