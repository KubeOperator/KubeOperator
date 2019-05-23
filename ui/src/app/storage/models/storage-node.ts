export class StorageNode {
  name: string;
  username: string;
  password: string;
  port: string;
  ip: string;
  status = 'running';
}

export class StorageGroup {
  name: string;
  nodes: StorageNode[] = [];
}
