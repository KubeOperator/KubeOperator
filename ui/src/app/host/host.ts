export class Host {
  id: string;
  name: string;
  ip: string;
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
  volumes: Volume[];
}

export class Volume {
  id: string;
  name: string;
  size: string;
  blank: boolean;
}
