import {Injectable} from '@angular/core';
import {HostService} from '../host/host.service';
import {Node} from '../node/node';
import {Host} from '../host/host';
import {Os, Template} from '../package/package';

@Injectable({
  providedIn: 'root'
})

export class CheckResult {
  passed: string[] = [];
  failed: string[] = [];
}

export class DeviceCheckService {

  constructor() {
  }

  checkOs(nodes: Node[], hosts: Host[], template: Template): CheckResult {
    const checkResult = new CheckResult();
    nodes.forEach(node => {
      const host = this.getHostById(node.host, hosts);
      template.roles.forEach(role => {
        if (node.roles.includes(role.name)) {
          role.meta.allow_os.forEach(o => {
            if (host.os === o.name && o.version.includes(host.os_version.substr(0, 3))) {
              checkResult.passed.push(node.name);
            }
          });
        }
      });
    });
    nodes.forEach(node => {
      if (!checkResult.passed.includes(node.name)) {
        checkResult.failed.push(node.name);
      }
    });
    return checkResult;
  }

  checkCpu(nodes: Node[], hosts: Host[], template: Template): CheckResult {
    const checkResult = new CheckResult();
    nodes.forEach(node => {
      const host = this.getHostById(node.host, hosts);
      template.roles.forEach(role => {
        if (node.roles.includes(role.name)) {
          const cpuCore = role.meta.requires.device_require[0].minimal;
          if (host.cpu_core >= cpuCore) {
            checkResult.passed.push(node.name);
          }
        }
      });
    });
    nodes.forEach(node => {
      if (!checkResult.passed.includes(node.name)) {
        checkResult.failed.push(node.name);
      }
    });
    return checkResult;
  }


  checkMemory(nodes: Node[], hosts: Host[], template: Template): CheckResult {
    const checkResult = new CheckResult();
    nodes.forEach(node => {
      const host = this.getHostById(node.host, hosts);
      template.roles.forEach(role => {
        if (node.roles.includes(role.name)) {
          const memory = role.meta.requires.device_require[1].minimal;
          if ((host.memory / memory) > 800) {
            checkResult.passed.push(node.name);
          }
        }
      });
    });
    nodes.forEach(node => {
      if (!checkResult.passed.includes(node.name)) {
        checkResult.failed.push(node.name);
      }
    });
    return checkResult;
  }

  getHostById(hostId: string, hosts: Host[]): Host {
    let result: Host = null;
    hosts.forEach(h => {
      if (h.id === hostId) {
        result = h;
      }
    });
    return result;
  }
}
