import {Component, OnInit} from '@angular/core';
import {Cluster} from '../cluster/cluster';
import {ActivatedRoute} from '@angular/router';
import {ClusterService} from '../cluster/cluster.service';
import {Package} from '../package/package';
import { AuthTemplate} from './class/auth';
import {PackageService} from '../package/package.service';
import {AuthService} from './service/auth.service';

@Component({
  selector: 'app-auth',
  templateUrl: './auth.component.html',
  styleUrls: ['./auth.component.css']
})
export class AuthComponent implements OnInit {

  currentCluster: Cluster;
  pkg: Package;
  authTemplate: AuthTemplate = null;
  authTemplates: AuthTemplate[];

  constructor(private route: ActivatedRoute, private authService: AuthService,
              private clusterService: ClusterService, private packageService: PackageService,
  ) {
  }

  ngOnInit() {
    this.loadData();
  }

  onCancel() {
    this.loadData();
  }

  loadData() {
    this.route.parent.data.subscribe(data => {
      this.currentCluster = data['cluster'];
      this.packageService.getPackage(this.currentCluster.package).subscribe(pkg => {
        this.pkg = pkg;
        this.authService.listAuthTemplate().subscribe(temps => {
          this.authTemplates = temps;
          this.authTemplates.forEach(t => {
            if (t.name === this.currentCluster.auth_template) {
              this.authTemplate = t;
              this.authService.fullAuth(this.authTemplate, this.currentCluster.name);
            }
          });
        });
      });
    });
  }


  // setDefaultValue() {
  //   this.auth.options.forEach(op => {
  //     if (op.type !== 'parent') {
  //       op.value = op.default;
  //     } else {
  //       op.value = op.default;
  //       op.children.forEach(cop => {
  //         op.value[cop.name] = cop.default;
  //       });
  //     }
  //     this.auth.vars.forEach(v => {
  //       v.value = v.default;
  //     });
  //   });
  // }
  //
  onSubmit() {
    this.authService.configAuth(this.authTemplate, this.currentCluster);
  }

}
