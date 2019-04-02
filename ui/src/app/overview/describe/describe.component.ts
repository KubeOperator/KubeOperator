import {Component, Input, OnInit} from '@angular/core';
import {Cluster} from '../../cluster/cluster';
import {PackageService} from '../../package/package.service';
import {Portal, Template} from '../../package/package';
import {RoleService} from '../../node/role.service';
import {ClusterRoleService} from '../../cluster/cluster-role.service';

@Component({
  selector: 'app-describe',
  templateUrl: './describe.component.html',
  styleUrls: ['./describe.component.css']
})
export class DescribeComponent implements OnInit {

  @Input() currentCluster: Cluster;
  private portals: Portal[] = [];

  constructor(private packageService: PackageService, private roleService: ClusterRoleService) {
  }

  ngOnInit() {
    this.packageService.getPackage(this.currentCluster.package).subscribe(data => {
      const template: Template = data.meta.templates.filter((temp => {
        if (temp.name === this.currentCluster.template) {
          return true;
        }
      }))[0];
      this.roleService.getClusterRole(this.currentCluster.name, 'OSEv3').subscribe(role => {
        for (const key in role.vars) {
          if (key) {
            template.portals.forEach(p => {
              if (p.redirect.includes(key)) {
                p.redirect = p.redirect.replace('$' + key, role.vars[key]);
              }
            });
          }
        }
        this.portals = template.portals;
      });
    });
  }
}
