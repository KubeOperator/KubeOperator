import {Component, Input, OnInit} from '@angular/core';
import {Cluster} from '../cluster/cluster';
import {App} from '../package/package';
import {PackageService} from '../package/package.service';
import {ActivatedRoute} from '@angular/router';

@Component({
  selector: 'app-application',
  templateUrl: './application.component.html',
  styleUrls: ['./application.component.css']
})
export class ApplicationComponent implements OnInit {

  apps: App[] = [];
  currentCluster: Cluster;

  constructor(private packageService: PackageService, private route: ActivatedRoute) {
  }

  ngOnInit() {
    this.route.parent.data.subscribe(data => {
      this.currentCluster = data['cluster'];
    });

    this.packageService.getPackage(this.currentCluster.package).subscribe(data => {
      this.apps = data.meta.apps;
    });
  }


  getAppUrl(app: App) {
    return this.currentCluster.apps[app.url_key];
  }

}
