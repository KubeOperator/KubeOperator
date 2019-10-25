import {Node} from '../node/node';
import {Execution} from '../deploy/component/operater/execution';
import {App, Config, Template} from '../package/package';

export class Cluster {
  id: string;
  name: string;
  package: string;
  comment: string;
  template: string;
  deploy_type: string;
  cloud_provider: string;
  plan: string;
  worker_size = 1;
  persistent_storage: string;
  date_created: string;
  node: Node[];
  current_execution: Execution;
  status: string;
  resource: string;
  nodes: string[] = [];
  network_plugin: string;
  apps: {};
  region: string;
  zone: string;
  zones: string[] = [];
  meta: {} = {};
  configs: {} = {};

  constructor() {
    this.configs['SERVICE_CIDR'] = '10.68.0.0/16';
    this.configs['CLUSTER_CIDR'] = '172.20.0.0/16';
  }
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

export class ClusterConfigs {
  templates: Template[];
  storages: Storage[];
  networks: Network[];
  apps: App[];
}

export class Network {
  name: string;
  configs: Config[] = [];
}

export class Storage {
  name: string;
  deploy_type: string[] = [];
  provider: string[] = [];
  configs: Config[] = [];
}
