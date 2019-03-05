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
  template: Template;
  configs: Config[] = [];
  packages: Package[] = [];
  templates: Template[] = [];
  nodes: Node[] = [];
  hosts: Host[] = [];
  @Output() create = new EventEmitter<boolean>();
  loadingFlag = false;

  constructor(private tipService: TipService, private nodeService: NodeService, private clusterService: ClusterService,
              private packageService: PackageService, private relationService: RelationService, private hostService: HostService) {
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
  }

  packgeOnChange() {
    this.packages.forEach((pak) => {
      if (pak.name === this.cluster.package) {
        this.configs = pak.meta.configs;
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
    this.nodes = [];
    this.templates.forEach(tmp => {
      if (tmp.name === this.cluster.template) {
        tmp.roles.forEach(role => {
          if (!role.meta.hidden) {
            const roleNumber = role.meta.nodes_require[1];
            for (let i = 0; i < roleNumber; i++) {
              const node: Node = new Node();
              node.name = role.name + '-' + i;
              node.roles.push(role.name);
              this.nodes.push(node);
            }
          }
        });
      }
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

  fullNode() {
    this.nodes.forEach(node => {
      this.hosts.forEach(host => {
        if (node.host === host.id) {
          node.ip = host.ip;
          node.host_memory = host.memory;
          node.host_cpu_core = host.cpu_core;
          node.host_os = host.os;
          node.host_os_version = host.os_version;
        }
      });
    });
  }

  createNodes() {
    this.nodes.forEach(node => {
      this.nodeService.createNode(this.cluster.name, node).subscribe(data => {
        this.configCluster();
      });
    });
  }

  configCluster() {
    this.configs.forEach(config => {
      const extraConfig: ExtraConfig = new ExtraConfig();
      extraConfig.key = config.name;
      extraConfig.value = config.value;
      this.clusterService.configCluster(this.cluster.name, extraConfig).subscribe(() => {
        this.isSubmitGoing = false;
        this.createClusterOpened = false;
        this.create.emit(true);
      });
    });
  }


  getHostInfo(host: Host) {
    const template = '{N} [{C}核  {M}MB  {O}]';
    return template.replace('{C}', host.cpu_core.toString())
      .replace('{M}', host.memory.toString())
      .replace('{O}', host.os + host.os_version)
      .replace('{N}', host.name);
  }

  onCancel() {
    this.reset();
    this.createClusterOpened = false;
  }

}
