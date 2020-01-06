export class Host {
  id: string;
  name: string;
  ip: string;
  port: number;
  username: string;
  password: string;
  credential: string;
  cluster: string;
  comment: string;
  memory: number;
  os: string;
  os_version: string;
  cpu_core: number;
  region: string;
  zone: string;
  status: string;
  volumes: Volume[];
  has_gpu: boolean;
  gpus: GPU[] = [];
  conditions: Condition[] = [];
}

export class GPU {
  id: string;
  name: string;
}

export class Volume {
  id: string;
  name: string;
  size: string;
  blank: boolean;
}

export class Condition {
  status: boolean;
  message: string;
  reason: string;
  last_time: string;
}
