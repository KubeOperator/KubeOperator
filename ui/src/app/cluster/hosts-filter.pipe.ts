import {Pipe, PipeTransform} from '@angular/core';
import {Host} from '../host/host';
import {Node} from '../node/node';
import {findNode} from '@angular/compiler';

@Pipe({
  name: 'hostsFilter',
  pure: false

})
export class HostsFilterPipe implements PipeTransform {

  transform(hosts: Host[], nodes: Node[], node: Node): Host[] {
    if (!(hosts && node && nodes)) {
      return [];
    }
    const result: Host[] = [];
    hosts.forEach(host => {
      const flag = nodes.filter((n) => {
        return host.id === n.host && n !== node;
      });
      if (flag.length === 0) {
        result.push(host);
      }
    });
    return result;
  }
}

