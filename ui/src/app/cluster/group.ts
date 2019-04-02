import {Node} from '../node/node';
import {NodeVars} from '../package/package';

export class Group {
  name: string;
  nodes: Node[] = [];
  op: string;
  limit: number;
  node_sum = 0;
  node_vars: NodeVars[];
}
