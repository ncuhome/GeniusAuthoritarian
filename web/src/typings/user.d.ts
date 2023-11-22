declare namespace User {
  type Group = {
    id: number;
    name: string;
  };

  type Profile = {
    user: {
      id: number;
      name: string;
      avatar_url: string;
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

    type Verified = {
      token: string;
      callback: string;
    };
  }

  namespace Mfa {
    type New = {
      url: string;
      qrcode: string;
    };

    type Status = {
      mfa: boolean;
    };
  }

  namespace SSH {
    type KeyMode = "pem" | "ssh";

    type Keys = {
      username: string;
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

  namespace Passkey {
    type Cred = {
      id: number;
      created_at: number;
      last_used_at: number;
      name: string;
    };
  }

  namespace U2F {
      type Methods = "phone" | "mfa" | "passkey" | "";

    type Status = {
      prefer: Methods;
      phone: boolean;
      mfa: boolean;
      passkey: boolean;
    };

    type Result = {
      token: string;
      valid_before: number;
    };
  }
}
