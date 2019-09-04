export class CloudZone {
  cluster: string;
  networks: [] = [];
  storages: CloudStorage[] = [];
  resourcePool: [] = [];
}

export class CloudStorage {
  name: string;
  type: string;
  multipleHostAccess: boolean;
}
