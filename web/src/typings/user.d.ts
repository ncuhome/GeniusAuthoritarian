declare namespace User {
  type Group = {
    id: number;
    name: string;
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

  namespace Login {
    type Result = {
      token: string;
      groups: string[];
    };

    type ThirdParty = {
      token: string;
      mfa: boolean;
      callback?: string;
    };

    type Mfa = {
      token: string;
      callback: string;
    };
  }

  namespace Mfa {
    type New = {
      url: string;
      qrcode: string;
    };
  }

  namespace SSH {
    type Keys = {
      pem: {
        public: string;
        private: string;
      };
      ssh: {
        public: string;
        private: string;
      };
    };
  }
}
