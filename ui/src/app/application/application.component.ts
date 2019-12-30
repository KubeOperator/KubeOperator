import {Component, OnInit, ViewChild} from '@angular/core';
import {Cluster} from '../cluster/cluster';
import {App} from '../package/package';
import {ActivatedRoute} from '@angular/router';
import {ClusterService} from '../cluster/cluster.service';
import {NodeService} from '../node/node.service';
import * as clipboard from 'clipboard-polyfill';

@Component({
  selector: 'app-application',
  templateUrl: './application.component.html',
  styleUrls: ['./application.component.css']
})
export class ApplicationComponent implements OnInit {

  apps: App[] = [];
  currentCluster: Cluster;
  workers = [];
  workerIp = '';
  @ViewChild('alertModal', {static: true}) alertModal;
  openHost = false;


  constructor(private clusterService: ClusterService, private route: ActivatedRoute, private nodeService: NodeService) {
  }

  ngOnInit() {
    this.route.parent.data.subscribe(data => {
      const clusterName = data['cluster']['name'];
      this.clusterService.getCluster(clusterName).subscribe(c => {
        this.currentCluster = c;
        this.clusterService.getClusterConfigs().subscribe(d => {
          this.apps = d.apps;
        });
      });
      this.nodeService.listNodes(this.currentCluster.name).subscribe(data => {
        this.workers = data.filter((node) => {
          return node.roles.includes('worker');
        });
        if (this.workers.length > 0) {
          this.workerIp = this.workers[0].ip;
        }
      });
    });
  }


  getAppUrl(app: App) {
    return this.currentCluster.apps[app.url_key];
  }

  copy(text: string) {
    const success = clipboard.writeText(text);
    if (success) {
      const alert = this.alertModal;
      this.alertModal.showTip(false, '复制成功');
      this.sleep(1000).then(function (this) {
        alert.closeTip();
      });
    }
  }

  sleep(ms) {
    return new Promise(
      (resolve) => {
        setTimeout(resolve, ms);
      });
  }

}
