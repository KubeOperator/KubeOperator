import {Component, OnInit, Output} from '@angular/core';
import {BaseModelComponent} from '../../../shared/class/BaseModelComponent';
import {Project} from '../project';
import {ProjectService} from '../project.service';

@Component({
    selector: 'app-project-list',
    templateUrl: './project-list.component.html',
    styleUrls: ['./project-list.component.css']
})
export class ProjectListComponent extends BaseModelComponent<Project> implements OnInit {


    constructor(private projectService: ProjectService) {
        super(projectService);
    }

    ngOnInit(): void {
        super.ngOnInit();
    }

}
