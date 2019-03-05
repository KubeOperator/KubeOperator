import {Host} from '../host/host';

export class Node {
  // 'id', 'ip','name', 'vars', 'roles', 'host', 'host_memory', 'host_cpu_core', 'host_os', 'host_os_version'
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
}
