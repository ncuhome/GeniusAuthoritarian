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
  };

  type Accessible = {
    permitAll: Info[];
    accessible: {
      group: User.Group;
      app: Info[];
    }[];
  };
}
