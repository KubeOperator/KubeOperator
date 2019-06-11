import {Node} from '../node/node';
import {Execution} from '../deploy/component/operater/execution';

export class Cluster {
  id: string;
  name: string;
  package: string;
  comment: string;
  template: string;
  storage: string;
  date_created: string;
  current_task_id: string;
  auth_template: string;
  node: Node[];
  configs: ExtraConfig[];
  current_execution: Execution;
  status: string;
  resource: string;
  resource_version: string;
}

export class ExtraConfig {
  key: string;
  value: any;
}
