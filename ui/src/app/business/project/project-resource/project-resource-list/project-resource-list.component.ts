import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelComponent} from '../../../../shared/class/BaseModelComponent';
import {ProjectResource} from '../project-resource';
import {ProjectResourceService} from '../project-resource.service';
import {ActivatedRoute} from '@angular/router';
import {Project} from '../../project';
import {ResourceTypes} from '../../../../constant/shared.const';

@Component({
    selector: 'app-project-resource-list',
    templateUrl: './project-resource-list.component.html',
    styleUrls: ['./project-resource-list.component.css']
})
export class ProjectResourceListComponent extends BaseModelComponent<ProjectResource> implements OnInit {

    currentProject: Project = new Project();
    resourceType: string;

    constructor(private projectResourceService: ProjectResourceService,
                private route: ActivatedRoute) {
        super(projectResourceService);
    }

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentProject = data.project;
            this.resourceType = ResourceTypes.Host;
            this.pageBy();
        });
    }

    onCreateBy() {
        this.createEvent.emit(this.resourceType);
    }

    pageBy() {
        this.projectResourceService.pageBy(this.page, this.size, this.currentProject.id, this.resourceType).subscribe(res => {
            this.items = res.items;
            this.loading = false;
        });
    }
}
