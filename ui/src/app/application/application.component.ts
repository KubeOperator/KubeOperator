import {Component, Input, OnInit} from '@angular/core';
import {Cluster, ClusterConfigs} from '../cluster/cluster';
import {App} from '../package/package';
import {PackageService} from '../package/package.service';
import {ActivatedRoute} from '@angular/router';
import {ClusterService} from '../cluster/cluster.service';

@Component({
  selector: 'app-application',
  templateUrl: './application.component.html',
  styleUrls: ['./application.component.css']
})
export class ApplicationComponent implements OnInit {

  apps: App[] = [];
  currentCluster: Cluster;

  constructor(private clusterService: ClusterService, private route: ActivatedRoute) {
  }

  ngOnInit() {
    this.route.parent.data.subscribe(data => {
      const clusterName = data['cluster'];
      this.clusterService.getCluster(clusterName).subscribe(c => {
        this.currentCluster = c;
        this.clusterService.getClusterConfigs().subscribe(d => {
          this.apps = d.apps;
        });
      });
    });
  }


  getAppUrl(app: App) {
    return this.currentCluster.apps[app.url_key];
  }

}
