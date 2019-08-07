import {Node} from '../node/node';
import {Execution} from '../deploy/component/operater/execution';

export class Cluster {
  id: string;
  name: string;
  package: string;
  comment: string;
  template: string;
  deploy_type: string;
  plan: string;
  worker_size = 3;
  persistent_storage: string;
  date_created: string;
  auth_template: string;
  node: Node[];
  configs: ExtraConfig[];
  current_execution: Execution;
  status: string;
  resource: string;
  nodes: string[] = [];
  network_plugin: string;
  apps: {};
}

export class ExtraConfig {
  key: string;
  value: any;
}

export class Operation {
  name: string;
  comment: string;
  icon: string;
  event: string;
  redirect: string;
  display_on: string[];
}
