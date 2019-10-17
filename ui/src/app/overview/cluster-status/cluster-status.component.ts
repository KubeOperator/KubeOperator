import {Component, Input, OnInit, ViewChild} from '@angular/core';
import {Cluster} from '../../cluster/cluster';
import {NodeService} from '../../node/node.service';
import {Node} from '../../node/node';
import {ScaleComponent} from '../scale/scale.component';
import {Router} from '@angular/router';
import {OperaterService} from '../../deploy/component/operater/operater.service';

@Component({
  selector: 'app-cluster-status',
  templateUrl: './cluster-status.component.html',
  styleUrls: ['./cluster-status.component.css']
})
export class ClusterStatusComponent implements OnInit {

  @Input() currentCluster: Cluster;
  workers: Node[] = [];
  @ViewChild(ScaleComponent, {static: true}) scale: ScaleComponent;

  constructor(private nodeService: NodeService,
              private router: Router, private operaterService: OperaterService) {
  }

  ngOnInit() {
    this.nodeService.listNodes(this.currentCluster.name).subscribe(data => {
      this.workers = data.filter((node) => {
        return node.roles.includes('worker');
      });
    });
  }

  handleScale() {
    const params = {'num': this.scale.worker_size};
    this.operaterService.executeOperate(this.currentCluster.name, 'scale', params).subscribe(() => {
      this.redirect('deploy');
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

}
