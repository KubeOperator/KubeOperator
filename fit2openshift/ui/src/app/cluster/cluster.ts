import {Node} from '../node/node';


export class Cluster {
  id: string;
  name: string;
  package: string;
  comment: string;
  template: string;
  date_created: string;
  // my field
  node: Node[];
  configs: ExtraConfig[];
}

export class ExtraConfig {
  key: string;
  value: any;
}
