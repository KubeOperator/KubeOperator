import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {Cluster} from '../cluster';
import {TipService} from '../../tip/tip.service';
import {ClrWizard} from '@clr/angular';
import {Package, Template} from '../../package/package';
import {PackageService} from '../../package/package.service';
import {TipLevels} from '../../tip/tipLevels';
import {Node} from '../../node/node';
import {ClusterService} from '../cluster.service';
import {NodeService} from '../../node/node.service';

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
  packages: Package[] = [];
  templates: Template[] = [];
  nodes: Node[] = [];
  @Output() create = new EventEmitter<boolean>();
  loadingFlag = false;

  constructor(private tipService: TipService, private nodeService: NodeService, private clusterService: ClusterService,
              private packageService: PackageService) {
  }

  ngOnInit() {
    this.listPackages();
  }

  newCluster() {
    // 清空对象
    this.reset();
    this.createClusterOpened = true;
  }

  reset() {
    this.wizard.reset();
    this.cluster = new Cluster();
    this.template = null;
    this.templates = null;
    this.nodes = null;
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
    this.nodes = [];
    console.log(this.cluster.template);
    this.templates.forEach(tmp => {
      if (tmp.name === this.cluster.template) {
        tmp.roles.forEach(role => {
          if (!role.meta.hidden) {
            const name = role.name;
            console.log(role);
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
    this.clusterService.createCluster(this.cluster).subscribe(data => {
      this.createNodes(this.cluster.name);
      this.isSubmitGoing = false;
      this.createClusterOpened = false;
      this.create.emit(true);
    });
  }

  createNodes(clusterName) {
    this.isSubmitGoing = true;
    this.nodes.forEach(node => {
      this.nodeService.createNode(clusterName, node).subscribe();
    });
  }

  onCancel() {
    this.reset();
    this.createClusterOpened = false;
  }

}
