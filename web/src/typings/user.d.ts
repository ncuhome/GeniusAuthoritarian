declare namespace User {
  type Group = {
    id: number;
    name: string;
  };

  type ThirdPartyLoginResult = {
    token: string;
    mfa: boolean;
    callback?: string;
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
      mfa: boolean;
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
