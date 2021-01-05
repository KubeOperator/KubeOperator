import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {ProjectResource} from '../project-resource';
import {ProjectResourceService} from '../project-resource.service';
import {ActivatedRoute} from '@angular/router';
import {Project} from '../../project';
import {ResourceTypes} from '../../../../constant/shared.const';
import {TranslateService} from '@ngx-translate/core';
import {SessionUser} from "../../../../shared/auth/session-user";
import {ProjectMember} from "../../project-member/project-member";
import {AlertLevels} from "../../../../layout/common-alert/alert";
import {SessionService} from "../../../../shared/auth/session.service";
import {ProjectMemberService} from "../../project-member/project-member.service";
import {CommonAlertService} from "../../../../layout/common-alert/common-alert.service";

@Component({
    selector: 'app-project-resource-list',
    templateUrl: './project-resource-list.component.html',
    styleUrls: ['./project-resource-list.component.css']
})
export class ProjectResourceListComponent extends BaseModelDirective<ProjectResource> implements OnInit {

    currentProject: Project = new Project();
    resourceType: string;
    @Output() deleteEvent = new EventEmitter<any>();

    constructor(private projectResourceService: ProjectResourceService,
                private route: ActivatedRoute,
                private translateService: TranslateService,
                private sessionService: SessionService,
                private projectMemberService: ProjectMemberService,
                private commonAlertService: CommonAlertService) {
        super(projectResourceService);
    }

    user: SessionUser;
    currentMember: ProjectMember;

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentProject = data.project;
            this.resourceType = ResourceTypes.Host;
            this.pageBy();
            const p = this.sessionService.getCacheProfile();
            this.user = p.user;
            if (!this.user.isAdmin) {
                this.projectMemberService.getByUser(this.user.name, this.currentProject.name).subscribe(res => {
                    this.currentMember = res;
                }, err => {
                    this.commonAlertService.showAlert(err.error.msg, AlertLevels.ERROR);
                });
            }
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
        this.loading = true;
        this.projectResourceService.pageBy(this.page, this.size, this.currentProject.name, this.resourceType).subscribe(res => {
            this.items = res.items;
            this.loading = false;
        }, error => {
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
