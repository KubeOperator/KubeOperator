export class CloudZone {
  cluster: string;
  networks: [] = [];
  storages: CloudStorage[] = [];
  resourcePool: [] = [];
  securityGroups: [] = [];
  networkList: Network[] = [];
  floatingNetworkList: Network[] = [];
}

export class CloudStorage {
  name: string;
  type: string;
  multipleHostAccess: boolean;
}

export class Network {
  name: string;
  id: string;
  subnetList: Subnet[] = [];
}

export class Subnet {
  name: string;
  id: string;
}
