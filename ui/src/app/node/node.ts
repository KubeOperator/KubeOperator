import {Condition, Host} from '../host/host';

export class Node {
  id: string;
  name: string;
  ip: string;
  vars: {} = {};
  roles: any[] = [];
  host: string;
  host_memory: number;
  host_cpu_core: number;
  host_os: string;
  host_os_version: string;
  delete = true;
  volumes: string[] = [];
  status: string;
  conditions: Condition[] = [];
  info: {} = {};
}
