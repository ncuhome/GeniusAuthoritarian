declare namespace App {
  type LoginInfo = {
    name: string;
    host: string;
  };

  type Info = {
    id: number;
    name: string;
    callback: string;
    permitAllGroup: boolean;
    views: number;
  };

  type Owned = Info & {
    appCode: string;
  };

  type Detailed = Owned & {
    groups: User.Group[];
  };

  type New = Detailed & {
    appSecret: string;
  };

  type Accessible = {
    permitAll: Info[];
    accessible: {
      group: User.Group;
      app: Info[];
    }[];
  };
}
