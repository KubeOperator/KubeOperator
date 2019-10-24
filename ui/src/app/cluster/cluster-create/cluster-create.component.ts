import {Component, EventEmitter, OnDestroy, OnInit, Output, ViewChild} from '@angular/core';
import {Cluster, ExtraConfig} from '../cluster';
import {ClrWizard} from '@clr/angular';
import {Config, Network, Package, Template} from '../../package/package';
import {PackageService} from '../../package/package.service';
import {ClusterService} from '../cluster.service';
import {NodeService} from '../../node/node.service';
import {RelationService} from '../relation.service';
import {Host, Volume} from '../../host/host';
import {Node} from '../../node/node';
import {HostService} from '../../host/host.service';
import {Group} from '../group';
import {CheckResult, DeviceCheckService} from '../device-check.service';
import {Subject} from 'rxjs';
import {NgForm} from '@angular/forms';
import {debounceTime} from 'rxjs/operators';
import {SettingService} from '../../setting/setting.service';
import {PlanService} from '../../plan/plan.service';
import {Plan} from '../../plan/plan';
import {CommonAlertService} from '../../base/header/common-alert.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';
import {Storage} from '../cluster';
import {Storage as StorageItem} from '../cluster';
import {StorageService} from '../storage.service';

export const CHECK_STATE_PENDING = 'pending';
export const CHECK_STATE_SUCCESS = 'success';
export const CHECK_STATE_FAIL = 'fail';

@Component({
  selector: 'app-cluster-create',
  templateUrl: './cluster-create.component.html',
  styleUrls: ['./cluster-create.component.css']
})


export class ClusterCreateComponent implements OnInit, OnDestroy {


  @ViewChild('wizard', {static: true}) wizard: ClrWizard;
  createClusterOpened: boolean;
  isSubmitGoing = false;
  cluster: Cluster = new Cluster();
  template: Template = new Template();
  configs: Config[] = [];
  package: Package;
  packages: Package[] = [];
  templates: Template[] = [];
  networks: Network[] = [];
  network: Network = null;
  storages: Storage[] = [];
  storageList: StorageItem[] = [];
  storage: Storage = null;
  nodes: Node[] = [];
  hosts: Host[] = [];
  groups: Group[] = [];
  plans: Plan[] = [];
  plan: Plan;
  checkCpuState = CHECK_STATE_PENDING;
  checkMemoryState = CHECK_STATE_PENDING;
  checkOsState = CHECK_STATE_PENDING;
  checkCpuResult: CheckResult = new CheckResult();
  checkMemoryResult: CheckResult = new CheckResult();
  checkOsResult: CheckResult = new CheckResult();
  suffix = 'f2o';
  @ViewChild('basicFrom', {static: true}) basicForm: NgForm;
  @ViewChild('storageForm', {static: true}) storageForm: NgForm;
  @ViewChild('networkForm', {static: true}) networkForm: NgForm;
  @ViewChild('nodeForm', {static: false}) nodeForm: NgForm;
  @ViewChild('configForm', {static: true}) configForm: NgForm;
  isNameValid = true;
  nameTooltipText = '只允许小写英文字母! 请勿包含特殊符号！';
  checkOnGoing = false;
  Manual = 'MANUAL';
  Automatic = 'AUTOMATIC';
  clusterNameChecker: Subject<string> = new Subject<string>();

  @Output() create = new EventEmitter<boolean>();

  constructor(private alertService: CommonAlertService, private nodeService: NodeService, private clusterService: ClusterService
    , private packageService: PackageService, private relationService: RelationService,
              private hostService: HostService, private deviceCheckService: DeviceCheckService,
              private settingService: SettingService, private planService: PlanService, private storageService: StorageService) {
  }

  ngOnInit() {
    this.clusterNameChecker.pipe(debounceTime(3000)).subscribe(() => {
      const cluster_name = this.basicForm.controls['cluster_name'];
      if (cluster_name) {
        this.isNameValid = cluster_name.valid;
        if (this.isNameValid) {
          if (!this.checkOnGoing) {
            this.checkOnGoing = true;
            this.clusterService.getCluster(this.cluster.name).subscribe(data => {
              this.checkOnGoing = false;
              this.nameTooltipText = '集群名称 ' + this.cluster.name + '已存在！';
              this.isNameValid = false;
            }, error1 => {
              this.checkOnGoing = false;
            });
          }
        }
      }
    });
    this.settingService.getSetting('domain_suffix').subscribe(data => {
      this.suffix = '.' + data.value;
    });
  }

  ngOnDestroy(): void {
    this.clusterNameChecker.unsubscribe();
  }

  public get isBasicFormValid(): boolean {
    return this.basicForm && this.basicForm.valid && this.isNameValid && !this.checkOnGoing && this.cluster.package !== '';
  }

  handleValidation(): void {
    const cont = this.basicForm.controls['cluster_name'];
    if (cont) {
      this.clusterNameChecker.next(cont.value);
    }
  }


  onNetworkChange() {
    this.networks.forEach(network => {
      if (this.cluster.network_plugin === network.name) {
        this.network = network;
      }
    });
  }

  onStorageChange() {
    if (this.cluster.persistent_storage === 'nfs') {
      this.storageService.list(this.cluster.persistent_storage).subscribe(data => {
        this.storageList = data;
      });
    }
    this.storages.forEach(storage => {
      if (this.cluster.persistent_storage === storage.name) {
        this.storage = storage;
      }
    });
  }

  loadClusterConfig() {
    this.clusterService.getClusterConfigs().subscribe(data => {
      this.templates = data.templates;
      this.storages = data.storages;
      this.networks = data.networks;
    });
  }

  loadStorage() {
    if (this.cluster.deploy_type) {
      this.storages = this.storages.filter(data => {
        return data.deploy_type.includes(this.cluster.deploy_type);
      });
      if (this.cluster.deploy_type === 'AUTOMATIC') {
        this.storages = this.storages.filter(data => {
          return data.provider.includes(this.plan.provider);
        });
      }
    }
  }

  newCluster() {
    this.reset();
    this.createClusterOpened = true;
    this.listPackages();
    this.getAllHost();
    this.listPlans();
    this.loadClusterConfig();
  }


  getAllHost() {
    this.hostService.listHosts().subscribe(data => {
      console.log(this.hosts);
      this.hosts = data;
    }, error => {
      console.log(error);
    });
  }

  reset() {
    this.wizard.reset();
    this.basicForm.resetForm();
    this.cluster = new Cluster();
    this.cluster.template = '';
    this.template = null;
    this.templates = [];
    this.nodes = [];
    this.configs = [];
    this.groups = null;
    this.storage = null;
    this.network = null;
    this.networks = null;
    this.resetCheckState();
  }


  listPackages() {
    this.packageService.listPackage().subscribe(data => {
      this.packages = data;
    }, error => {
      this.alertService.showAlert('加载离线包错误!: \n' + error, AlertLevels.ERROR);
    });
  }

  listPlans() {
    this.planService.listPlan().subscribe(data => {
      this.plans = data;
    });
  }

  planOnChange() {
    this.plans.forEach(plan => {
      if (this.cluster.plan === plan.name) {
        this.plan = plan;
        this.templates.forEach(template => {
          if (template.deploy_type === plan.deploy_template) {
            this.template = template;
            this.cluster.template = template.name;
          }
        });
      }
    });
  }

  onWorkerSizeChange() {
    if (this.cluster.worker_size < 3) {
      this.cluster.worker_size = 3;
    }
  }

  templateOnChange() {
    this.templates.forEach(template => {
      if (template.name === this.cluster.template) {
        this.template = template;
        this.configs.concat(this.template.private_config);
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
          host.volumes.forEach(volume => {
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
    const no = group.node_sum + 1;
    node.name = group.name + no + '.' + this.cluster.name + this.suffix;
    group.node_sum++;
    node.roles.push(group.name);
    group.nodes.push(node);
    this.nodes.push(node);
    console.log(this.nodes);
  }

  fullNode() {
    this.resetCheckState();
    this.deviceCheck();
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

  onSubmit() {
    if (this.isSubmitGoing) {
      return;
    }
    this.isSubmitGoing = true;
    this.clusterService.createCluster(this.cluster).subscribe(data => {
      this.cluster = data;
      if (this.nodes) {
        this.createNodes();
      }
    });
  }

  createNodes() {
    const promises: Promise<{}>[] = [];
    this.nodes.forEach(node => {
      promises.push(this.nodeService.createNode(this.cluster.name, node).toPromise());
    });

    Promise.all(promises).then(() => {
      this.finishForm();
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

  finishForm() {
    this.isSubmitGoing = false;
    this.createClusterOpened = false;
    this.create.emit(true);
  }

  getHostInfo(host: Host) {
    const template = '{N} [{C}核  {M}MB  {O}]';
    return template.replace('{C}', host.cpu_core.toString())
      .replace('{M}', host.memory.toString())
      .replace('{O}', host.os + host.os_version)
      .replace('{N}', host.name);
  }

  getVolumeInfo(volume: Volume) {
    const template = '{N}  {S}';
    return template.replace('{N}', volume.name).replace('{S}', volume.size);
  }

  canNetworkNext() {
    let result = true;
    if (!this.cluster.network_plugin) {
      result = false;
    }
    if (this.network) {
      this.network.configs.some(cfg => {
        if (!cfg.value) {
          result = false;
          return true;
        }
      });
    }
    return result;
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
    return true;
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
