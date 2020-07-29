import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelComponent} from '../../../../shared/class/BaseModelComponent';
import {ProjectResource} from '../project-resource';
import {ProjectResourceService} from '../project-resource.service';
import {ActivatedRoute} from '@angular/router';
import {Project} from '../../project';
import {ResourceTypes} from '../../../../constant/shared.const';
import {PlanService} from '../../../deploy-plan/plan/plan.service';
import {TranslateService} from '@ngx-translate/core';

@Component({
    selector: 'app-project-resource-list',
    templateUrl: './project-resource-list.component.html',
    styleUrls: ['./project-resource-list.component.css']
})
export class ProjectResourceListComponent extends BaseModelComponent<ProjectResource> implements OnInit {

    currentProject: Project = new Project();
    resourceType: string;
    @Output() deleteEvent = new EventEmitter<any>();

    constructor(private projectResourceService: ProjectResourceService,
                private route: ActivatedRoute,
                private translateService: TranslateService) {
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

    changeTab(resourceType) {
        this.resourceType = resourceType;
        this.pageBy();
    }

    pageBy() {
        this.projectResourceService.pageBy(this.page, this.size, this.currentProject.id, this.resourceType).subscribe(res => {
            this.items = res.items;
            this.loading = false;
        });
    }

    onDelete() {
        this.deleteEvent.emit({items: this.selected, resourceType: this.resourceType});
    }

    getDeployName(name: string) {
        switch (name) {
            case 'SINGLE':
                return this.translateService.instant('APP_PLAN_DEPLOY_TEMPLATE_SINGLE');
            case 'MULTIPLE':
                return this.translateService.instant('APP_PLAN_DEPLOY_TEMPLATE_MULTIPLE');
            default:
                return 'None';
        }
    }
}
