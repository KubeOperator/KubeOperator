export class StorageNode {
  name: string;
  username: string;
  password: string;
  port: string;
  ip: string;
}

export class StorageGroup {
  name: string;
  nodes: StorageNode[] = [];
}
