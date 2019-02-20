import {Pipe, PipeTransform} from '@angular/core';
import {Host} from './host';
import {Node} from '../node/node';
import {find} from 'rxjs/operators';

@Pipe({
  name: 'hostFilter'
})
export class HostFilterPipe implements PipeTransform {

  transform(value: Host[], currentNode: Node, nodeList: Node[]): any {
    const result: Host[] = [];
    value.forEach(host => {
      if (host.cluster === '无' && !this.finds(host, currentNode, nodeList)) {
        result.push(host);
      }
    });
    return result;
  }

  private finds(host: Host, currentNode: Node, node: Node[]): boolean {
    let flag = false;
    for (let i = 0; i < node.length; i++) {
      console.log(node[i].host + '     ' + host.id);
      console.log(currentNode === node[i]);
      if (node[i].host === host.id && currentNode !== node[i]) {
        flag = true;
        break;
      }
    }
    return flag;
  }
}
