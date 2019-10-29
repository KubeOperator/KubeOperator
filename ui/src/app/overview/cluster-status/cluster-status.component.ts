import {Component, Input, OnInit, ViewChild} from '@angular/core';
import {Cluster} from '../../cluster/cluster';
import {NodeService} from '../../node/node.service';
import {Node} from '../../node/node';
import {ScaleComponent} from '../scale/scale.component';
import {Router} from '@angular/router';
import {OperaterService} from '../../deploy/component/operater/operater.service';
import {ClusterHealthService} from '../../cluster-health/cluster-health.service';
import {ClusterHealth} from '../../cluster-health/cluster-health';


@Component({
  selector: 'app-cluster-status',
  templateUrl: './cluster-status.component.html',
  styleUrls: ['./cluster-status.component.css']
})
export class ClusterStatusComponent implements OnInit {

  @Input() currentCluster: Cluster;
  workers: Node[] = [];
  @ViewChild(ScaleComponent, {static: true}) scale: ScaleComponent;
  clusterHealth: ClusterHealth = new ClusterHealth();


  constructor(private nodeService: NodeService, private clusterHealthService: ClusterHealthService,
              private router: Router, private operaterService: OperaterService) {
  }

  ngOnInit() {
    this.nodeService.listNodes(this.currentCluster.name).subscribe(data => {
      this.workers = data.filter((node) => {
        return node.roles.includes('worker');
      });
    });
    this.getClusterStatus();
  }

  handleScale() {
    const params = {'num': this.scale.worker_size};
    this.operaterService.executeOperate(this.currentCluster.name, 'scale', params).subscribe(() => {
      this.redirect('deploy');
    }, error => {
      this.scale.opened = false;
    });
  }

  redirect(url: string) {
    if (url) {
      const linkUrl = ['kubeOperator', 'cluster', this.currentCluster.name, url];
      this.router.navigate(linkUrl);
    }
  }

  onScale() {
    this.scale.worker_size = this.workers.length;
    this.scale.opened = true;
  }

  toHealth() {
    this.redirect('health');
  }

  getClusterStatus() {
    this.clusterHealth.data = [];
    if (this.currentCluster.status === 'READY') {
      return;
    }
    this.clusterHealthService.listClusterHealth(this.currentCluster.name).subscribe(res => {
      this.clusterHealth = res;
    }, error1 => {
      this.clusterHealth.data = [];
    });
  }

  getServiceStatus(type) {
    if (this.clusterHealth == null || this.clusterHealth.data.length === 0 ) {
      return '';
    }
    let status = 'UNKNOWN';
    for (const ch of this.clusterHealth.data) {
      if (ch.job === type) {
        if (ch.rate === 100) {
          status =  'RUNNING';
        } else {
          status =  'WARNING';
        }
      }
    }
    return status;
  }
}
