
export class ClusterHealth {
  data: HealthData[];
  success: Boolean;
  rate: number;
}

export class HealthData {
  job: string;
  data: Data[];
  rate: number;
}

export class Data {
  key: string;
  value: any;
}
