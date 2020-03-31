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
import * as globals from '../../globals';
import {CephService} from '../../ceph/ceph.service';
import {SessionService} from '../../shared/session.service';
import {ItemService} from '../../item/item.service';
import {SessionUser} from '../../shared/session-user';

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
  global_domain: string;
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
  @ViewChild('workerForm', {static: true}) workerForm: NgForm;
  @ViewChild('alertModal', {static: true}) alertModal;
  isNameValid = true;
  nameTooltipText = '只允许小写英文字母! 请勿包含特殊符号！';
  checkOnGoing = false;
  Manual = 'MANUAL';
  Automatic = 'AUTOMATIC';
  clusterNameChecker: Subject<string> = new Subject<string>();
  name_pattern = globals.cluster_name_pattern;
  domain_pattern = globals.domain_pattern;
  name_pattern_tip = globals.cluster_name_pattern_tip;
  itemName: string;
  items = [];

  @Output() create = new EventEmitter<boolean>();

  constructor(private alertService: CommonAlertService, private nodeService: NodeService, private clusterService: ClusterService
    , private packageService: PackageService, private relationService: RelationService,
              private hostService: HostService, private deviceCheckService: DeviceCheckService,
              private settingService: SettingService, private planService: PlanService, private storageService: StorageService,
              private cephService: CephService, private sessionService: SessionService, private itemService: ItemService) {
  }

  ngOnInit() {
    this.clusterNameChecker.pipe(debounceTime(500)).subscribe(() => {
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
              this.alertModal.showTip(true, this.nameTooltipText);
            }, error1 => {
              this.checkOnGoing = false;
              this.alertModal.closeTip();
            });
          }
        }
      }
    });

    const profile = this.sessionService.getCacheProfile();
    const user = profile.user;

    this.itemService.listItem().subscribe(res => {
      this.items = res;
      if (!user.is_superuser) {
        this.items = this.sessionService.getManageItems(this.items);
      }
    });
  }

  ngOnDestroy(): void {
    this.clusterNameChecker.unsubscribe();
  }

  reset() {
    this.wizard.reset();
    this.basicForm.resetForm();
    this.storageForm.resetForm();
    this.networkForm.resetForm();
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
    this.settingService.getSettings().subscribe(data => {
      this.global_domain = data['domain_suffix'];
      this.cluster.cluster_doamin_suffix = this.global_domain;
    });
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
      this.storageService.list(this.cluster.persistent_storage, this.itemName).subscribe(data => {
        this.storageList = data;
      });
    }
    if (this.cluster.persistent_storage === 'external-ceph') {
      this.storageService.list('ceph', this.itemName).subscribe(data => {
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

  onChangeItem(itemName) {
    this.itemName = itemName;
    this.getItemResources();
  }

  newCluster() {
    this.reset();
    this.createClusterOpened = true;
    this.listPackages();
    this.loadClusterConfig();
  }

  getItemResources() {
    this.getAllHost();
    this.listPlans();
  }

  getAllHost() {
    this.hostService.byItem(this.itemName).subscribe(data => {
      console.log(data);
      this.hosts = data.filter(host => {
        return !host.cluster;
      });
    });
  }


  listPackages() {
    this.packageService.listPackage().subscribe(data => {
      this.packages = data;
      console.log(data);
    }, error => {
      this.alertService.showAlert('加载离线包错误!: \n' + error, AlertLevels.ERROR);
    });
  }

  listPlans() {
    this.planService.listItemPlan(this.itemName).subscribe(data => {
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
    if (this.cluster.worker_size < 1) {
      this.cluster.worker_size = 1;
    }
  }

  templateOnChange() {
    this.templates.forEach(template => {
      if (template.name === this.cluster.template) {
        this.template = template;
        this.configs.concat(this.template.private_config);
      }
    });
  }

  workerSizeOnChange() {
    if (this.cluster.worker_size < 1) {
      this.cluster.worker_size = 1;
    }
    if (this.cluster.worker_size > this.hosts.length) {
      this.cluster.worker_size = this.hosts.length;
    }

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
            if (group.name === 'worker') {
              group.limit = this.cluster.worker_size;
            }
            for (let i = group.node_sum; i < group.limit; i++) {
              this.addNode(group, false);
            }
            this.groups.push(group);
          }
        });
      }
    });
  }

  getPackageVars() {
    let vars = null;
    this.packages.forEach(p => {
      if (p.name === this.cluster.package) {
        vars = p.meta.vars;
      }
    });
    return vars;
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
    node.name = group.name + no + '.' + this.cluster.name + '.' + this.cluster.cluster_doamin_suffix;
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
    this.cluster.item_name = this.itemName;
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
    for (const node of this.nodes) {
      if (!node.host) {
        return false;
      }
    }
    return true;
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
