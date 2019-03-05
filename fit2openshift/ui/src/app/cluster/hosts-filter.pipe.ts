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
      if (host.cluster === '无') {
        const others: Node[] = [];
        nodes.forEach(n => {
          if (n !== node) {
            others.push(n);
          }
        });

        let flag = false;
        for (let i = 0; i < others.length; i++) {
          if (host.id === others[i].host) {
            flag = true;
          }
        }
        if (!flag) {
          result.push(host);
        }
      }
    });
    return result;
  }
}

