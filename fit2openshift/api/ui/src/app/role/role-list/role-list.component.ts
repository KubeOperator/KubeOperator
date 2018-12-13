import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from "@angular/router";
import { RoleService } from "../role.service";
import { Role } from "../role";
import {Project} from "../../project/project";

@Component({
  selector: 'app-role-list',
  templateUrl: './role-list.component.html',
  styleUrls: ['./role-list.component.css']
})
export class RoleListComponent implements OnInit {
  currentProject: Project;
  projectId: string;
  roles: Role[];

  constructor(
    private service: RoleService,
    private route: ActivatedRoute) {

    this.currentProject = this.route.snapshot.parent.data['project'];
    this.projectId = this.route.snapshot.parent.params['project'];
  }

  ngOnInit() {
    console.log(this.route.paramMap);
    this.getRoles(this.projectId);
  }

  getRoles(projectId: string){
    this.service.getRoles(projectId)
      .subscribe(roles => this.roles = roles)
  }
}
