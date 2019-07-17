import {Component, Input, OnInit} from '@angular/core';
import {Cluster} from '../../cluster/cluster';
import {PackageService} from '../../package/package.service';
import {App} from '../../package/package';

@Component({
  selector: 'app-apps',
  templateUrl: './apps.component.html',
  styleUrls: ['./apps.component.css']
})
export class AppsComponent implements OnInit {

  @Input() currentCluster: Cluster;
  private apps: App[] = [];

  constructor(private packageService: PackageService) {
  }

  ngOnInit() {
    this.packageService.getPackage(this.currentCluster.package).subscribe(data => {
      this.apps = data.meta.apps;
    });
  }


  getAppUrl(app: App) {
    return this.currentCluster.apps[app.url_key];
  }

}
