import {Component, OnInit} from '@angular/core';
import {Cluster} from '../cluster/cluster';
import {ActivatedRoute, Router} from '@angular/router';
import {ClusterService} from '../cluster/cluster.service';
import {OperaterService} from '../deploy/component/operater/operater.service';

@Component({
  selector: 'app-f5-big-ip',
  templateUrl: './f5-big-ip.component.html',
  styleUrls: ['./f5-big-ip.component.css']
})
export class F5BigIpComponent implements OnInit {

  currentCluster: Cluster;

  constructor(private route: ActivatedRoute, private clusterService: ClusterService,
              private router: Router, private operaterService: OperaterService) {
  }


  ngOnInit() {
    this.route.parent.data.subscribe(data => {
      this.currentCluster = data['cluster'];
      this.refreshCluster();
    });
  }

  refreshCluster() {
    this.clusterService.getCluster(this.currentCluster.name).subscribe(cluster => {
      this.currentCluster = cluster;
    });
  }

  onCommit() {
    this.clusterService.updateCluster(this.currentCluster).subscribe(() => {
      this.operaterService.executeOperate(this.currentCluster.name, 'bigip-config').subscribe(data => {
        this.redirect('deploy');
      });
    });
  }

  redirect(url: string) {
    if (url) {
      const linkUrl = ['kubeOperator', 'cluster', this.currentCluster.name, url];
      this.router.navigate(linkUrl);
    }
  }
}
