
export class ClusterHealth {
  data: HealthData[];
  status: string;
}

export class HealthData {
  type: string;
  data: Data[];
}

export class Data {
  key: string;
  value: any;
}
