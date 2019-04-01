import {Node} from '../node/node';


export class Cluster {
  id: string;
  name: string;
  package: string;
  comment: string;
  template: string;
  date_created: string;
  current_task_id: string;
  state: string;
  // my field
  node: Node[];
  configs: ExtraConfig[];
}

export class ExtraConfig {
  key: string;
  value: any;
}
