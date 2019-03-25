import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {Cluster, ExtraConfig} from '../cluster';
import {TipService} from '../../tip/tip.service';
import {ClrWizard} from '@clr/angular';
import {Config, Package, Template} from '../../package/package';
import {PackageService} from '../../package/package.service';
import {TipLevels} from '../../tip/tipLevels';
import {ClusterService} from '../cluster.service';
import {NodeService} from '../../node/node.service';
import {RelationService} from '../relation.service';
import {Host} from '../../host/host';
import {Node} from '../../node/node';
import {HostService} from '../../host/host.service';
import {Group} from '../group';
import {CheckResult, DeviceCheckService} from '../device-check.service';
import {config} from 'rxjs';

export const CHECK_STATE_PENDING = 'pending';
export const CHECK_STATE_SUCCESS = 'success';
export const CHECK_STATE_FAIL = 'fail';

@Component({
  selector: 'app-cluster-create',
  templateUrl: './cluster-create.component.html',
  styleUrls: ['./cluster-create.component.css']
})


export class ClusterCreateComponent implements OnInit {


  @ViewChild('wizard') wizard: ClrWizard;
  createClusterOpened: boolean;
  isSubmitGoing = false;
  cluster: Cluster = new Cluster();
  template: Template = new Template();
  configs: Config[] = [];
  packages: Package[] = [];
  templates: Template[] = [];
  nodes: Node[] = [];
  hosts: Host[] = [];
  groups: Group[] = [];
  checkCpuState = CHECK_STATE_PENDING;
  checkMemoryState = CHECK_STATE_PENDING;
  checkOsState = CHECK_STATE_PENDING;
  checkCpuResult: CheckResult = new CheckResult();
  checkMemoryResult: CheckResult = new CheckResult();
  checkOsResult: CheckResult = new CheckResult();
  suffix = '.f2o';


  @Output() create = new EventEmitter<boolean>();
  loadingFlag = false;

  constructor(private tipService: TipService, private nodeService: NodeService, private clusterService: ClusterService
    , private packageService: PackageService, private relationService: RelationService,
              private hostService: HostService, private deviceCheckService: DeviceCheckService) {
  }

  ngOnInit() {
    this.listPackages();
    this.getAllHost();
  }

  newCluster() {
    // 清空对象
    this.reset();
    this.createClusterOpened = true;
  }

  getAllHost() {
    this.hostService.listHosts().subscribe(data => {
      this.hosts = data;

    }, error => {
      console.log(error);
    });
  }

  reset() {
    this.wizard.reset();
    this.cluster = new Cluster();
    this.template = null;
    this.templates = null;
    this.nodes = null;
    this.configs = null;
    this.groups = null;
    this.resetCheckState();
  }


  packgeOnChange() {
    this.packages.forEach((pak) => {
      if (pak.name === this.cluster.package) {
        this.templates = pak.meta.templates;
      }
    });
  }

  listPackages() {
    this.packageService.listPackage().subscribe(data => {
      this.packages = data;
    }, error => {
      this.tipService.showTip('加载离线包错误!: \n' + error, TipLevels.ERROR);
    });
  }

  templateOnChange() {
    this.templates.forEach(template => {
      if (template.name === this.cluster.template) {
        this.template = template;
        this.configs = template.private_config;
        this.configs.forEach(c => {
          c.value = c.default;
        });
      }
    });
    this.nodes = [];
    this.groups = [];
    this.templates.forEach(tmp => {
      if (tmp.name === this.cluster.template) {
        tmp.roles.forEach(role => {
          if (!role.meta.hidden) {
            const group: Group = new Group();
            group.node_vars = role.meta.node_vars;
            group.name = role.name;
            group.op = role.meta.requires.nodes_require[0];
            group.limit = role.meta.requires.nodes_require[1];
            for (let i = group.node_sum; i < group.limit; i++) {
              this.addNode(group, false);
            }
            this.groups.push(group);
          }
        });
      }
    });
  }

  onHostChange(node: Node) {
    if (node.host) {
      node.volumes = [];
      this.hosts.forEach(host => {
        if (host.id === node.host) {
          host.info.volumes.forEach(volume => {
            node.volumes.push(volume.name);
          });
        }
      });
    }

  }

  deleteNode(group: Group, node: Node) {
    let indexG;
    let indexN;
    for (let i = 0; i < group.nodes.length; i++) {
      if (node.name === group.nodes[i].name) {
        indexG = i;
      }
    }
    for (let i = 0; i < this.nodes.length; i++) {
      if (node.name === this.nodes[i].name) {
        indexN = i;
      }
    }
    group.nodes.splice(indexG, 1);
    this.nodes.splice(indexN, 1);
    group.node_sum--;

  }

  addNode(group: Group, canDelete?: boolean) {
    const node: Node = new Node();
    if (canDelete !== undefined && canDelete !== null) {
      node.delete = canDelete;
    }
    node.name = group.name + '-' + group.node_sum + '.' + this.cluster.name + this.suffix;
    group.node_sum++;
    node.roles.push(group.name);
    group.nodes.push(node);
    this.nodes.push(node);
  }

  fullNode() {
    this.resetCheckState();
    this.deviceCheck();
    this.nodes.forEach(node => {
      this.hosts.forEach(host => {
        if (node.host === host.id) {
          node.ip = host.ip;
          node.host_memory = host.info.memory;
          node.host_cpu_core = host.info.cpu_core;
          node.host_os = host.info.os;
          node.host_os_version = host.info.os_version;
        }
      });
    });
  }

  onSubmit() {
    if (this.isSubmitGoing) {
      return;
    }
    this.isSubmitGoing = true;
    this.clusterService.createCluster(this.cluster).subscribe(data => {
      this.cluster = data;
      this.createNodes();
    });
  }


  createNodes() {
    const promises: Promise<{}>[] = [];
    this.nodes.forEach(node => {
      promises.push(this.nodeService.createNode(this.cluster.name, node).toPromise());
    });

    Promise.all(promises).then(() => {
      this.configCluster();
    });
  }

  canNodeNext(): boolean {
    let result = false;
    if (this.nodes) {
      this.nodes.some(node => {
        if (!node.host) {
          result = true;
          return true;
        }
      });
    }
    return result;
  }

  configCluster() {
    const promises: Promise<{}>[] = [];
    this.configs.forEach(c => {
      const extraConfig: ExtraConfig = new ExtraConfig();
      extraConfig.key = c.name;
      extraConfig.value = c.value;
      promises.push(this.clusterService.configCluster(this.cluster.name, extraConfig).toPromise());
      Promise.all(promises).then(() => {
        this.isSubmitGoing = false;
        this.createClusterOpened = false;
        this.create.emit(true);
      });
    });
  }


  getHostInfo(host: Host) {
    const template = '{N} [{C}核  {M}MB  {O}]';
    return template.replace('{C}', host.info.cpu_core.toString())
      .replace('{M}', host.info.memory.toString())
      .replace('{O}', host.info.os + host.info.os_version)
      .replace('{N}', host.name);
  }

  getHostById(hostId: string): Host {
    let h: Host;
    this.hosts.forEach(host => {
      if (host.id === hostId) {
        h = host;
      }
    });
    return h;
  }

  canConfigNext() {
    let result = true;
    if (this.configs) {
      this.configs.some(c => {
        if (c.value != null && c.value !== '') {
          result = false;
          return true;
        }
      });
    }
    return result;
  }


  replaceVolumeName(name: string, template: string): string {
    return template.replace('$value', name);
  }

  deviceCheck() {
    setTimeout(() => {
      this.checkCpu();
    }, 2000);
    setTimeout(() => {
      this.checkMemory();
    }, 4000);
    setTimeout(() => {
      this.checkOS();
    }, 6000);
  }

  checkCpu() {
    this.checkCpuResult = this.deviceCheckService.checkCpu(this.nodes, this.hosts, this.template);
    if (this.checkCpuResult.passed.length === this.nodes.length) {
      this.checkCpuState = CHECK_STATE_SUCCESS;
    } else {
      this.checkCpuState = CHECK_STATE_FAIL;
    }
  }

  checkMemory() {
    this.checkMemoryResult = this.deviceCheckService.checkMemory(this.nodes, this.hosts, this.template);
    if (this.checkMemoryResult.passed.length === this.nodes.length) {
      this.checkMemoryState = CHECK_STATE_SUCCESS;
    } else {
      this.checkMemoryState = CHECK_STATE_FAIL;
    }
  }

  checkOS() {
    this.checkOsResult = this.deviceCheckService.checkOs(this.nodes, this.hosts, this.template);
    if (this.checkOsResult.passed.length === this.nodes.length) {
      this.checkOsState = CHECK_STATE_SUCCESS;
    } else {
      this.checkOsState = CHECK_STATE_FAIL;
    }
  }

  resetCheckState() {
    this.checkCpuState = CHECK_STATE_PENDING;
    this.checkMemoryState = CHECK_STATE_PENDING;
    this.checkOsState = CHECK_STATE_PENDING;
  }

  canCheckNext() {
    if (this.checkOsState === CHECK_STATE_SUCCESS && this.checkMemoryState === CHECK_STATE_SUCCESS &&
      this.checkCpuState === CHECK_STATE_SUCCESS) {
      return true;
    }
    return false;
  }

  onCancel() {
    this.reset();
    this.createClusterOpened = false;
  }

}
