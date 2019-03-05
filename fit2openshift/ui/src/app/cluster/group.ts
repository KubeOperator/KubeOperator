import {Node} from '../node/node';

export class Group {
  name: string;
  nodes: Node[] = [];
  op: string;
  limit: number;
  node_sum = 0;
}
