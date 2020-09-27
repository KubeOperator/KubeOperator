import {Component, OnInit, Output} from '@angular/core';
import {Project} from '../project';
import {ProjectService} from '../project.service';
import {PermissionService} from '../../../shared/auth/permission.service';
import {BaseModelDirective} from "../../../shared/class/BaseModelDirective";

@Component({
    selector: 'app-project-list',
    templateUrl: './project-list.component.html',
    styleUrls: ['./project-list.component.css']
})
export class ProjectListComponent extends BaseModelDirective<Project> implements OnInit {


    constructor(private projectService: ProjectService,
                private permissionService: PermissionService) {
        super(projectService);
    }

    ngOnInit(): void {
        super.ngOnInit();
    }

    getProjectRole(projectName: string) {
        return this.permissionService.getProjectRole(projectName);
    }
}
