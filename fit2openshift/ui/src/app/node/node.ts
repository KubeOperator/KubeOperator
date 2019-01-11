export class Node {
  id: string;
  name: string;
  ip: string;
  port: number;
  username: string;
  password: string;
  vars: string;
  comment: string;
  roles: any[] = [];
  status: string;
}
