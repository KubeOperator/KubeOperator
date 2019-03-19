import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {FormControl, FormGroup} from '@angular/forms';
import {NodeService} from '../node.service';
import {Cluster} from '../../cluster/cluster';
import {Node} from '../node';
import {TipService} from '../../tip/tip.service';
import {TipLevels} from '../../tip/tipLevels';
import {HostService} from '../../host/host.service';
import {Host} from '../../host/host';
import {map} from 'rxjs/operators';
import {RoleService} from '../role.service';
import {Role} from '../role';

@Component({
  selector: 'app-node-create',
  templateUrl: './node-create.component.html',
  styleUrls: ['./node-create.component.css']
})
export class NodeCreateComponent implements OnInit {

  @Input() currentCluster: Cluster;
  node: Node = new Node();
  hosts: Host[] = [];
  roles: Role[] = [];
  staticBackdrop = true;
  closable = false;
  createNodeOpened: boolean;
  isSubmitGoing = false;
  @Output() create = new EventEmitter<boolean>();


  constructor(private hostService: HostService, private roleService: RoleService,
              private nodeService: NodeService, private tipService: TipService) {
  }

  ngOnInit() {
    this.listHosts();
    this.listRoles();
  }

  onSubmit() {
  }

  listHosts() {
    this.hostService.listHosts().pipe(map(data => {
      const hosts: Host[] = [];
      data.forEach(host => {
        if (host.cluster === '无') {
          hosts.push(host);
        }
      });
      return hosts;
    })).subscribe(data => {
      this.hosts = data;
    });
  }

  listRoles() {
    this.roleService.listRoles(this.currentCluster.name).pipe(map(da => {
      const roles: Role[] = [];
      da.forEach(role => {
        if (!role.meta['hidden']) {
          roles.push(role);
        }
      });
      return roles;
    })).subscribe(data => {
      this.roles = data;
      console.log(data);
    });
  }


  newNode() {
    this.node = new Node();
    this.createNodeOpened = true;
  }


  onCancel() {
    this.createNodeOpened = false;
  }
}
