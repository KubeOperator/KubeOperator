import {Component, OnInit} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {ProjectMember, ProjectMemberCreate} from '../project-member';
import {ProjectMemberService} from '../project-member.service';
import {ActivatedRoute} from '@angular/router';
import {Project} from '../../project';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../../layout/common-alert/alert';
import {SessionUser} from "../../../../shared/auth/session-user";
import {SessionService} from "../../../../shared/auth/session.service";

@Component({
    selector: 'app-project-member-list',
    templateUrl: './project-member-list.component.html',
    styleUrls: ['./project-member-list.component.css']
})
export class ProjectMemberListComponent extends BaseModelDirective<ProjectMember> implements OnInit {


    currentProject: Project = new Project();
    batchItems: ProjectMemberCreate[] = [];
    currentMember: ProjectMember;
    user: SessionUser;

    constructor(private projectMemberService: ProjectMemberService,
                private route: ActivatedRoute,
                private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService,
                private sessionService: SessionService) {
        super(projectMemberService);
    }

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentProject = data.project;
            this.pageBy();
        });
        const p = this.sessionService.getCacheProfile();
        this.user = p.user;
        if (!this.user.isAdmin) {
            this.projectMemberService.getByUser(this.user.name, this.currentProject.name).subscribe(data => {
                this.currentMember = data;
            }, err => {
                this.commonAlertService.showAlert(err.error.msg, AlertLevels.ERROR);
            });
        }
    }

    pageBy() {
        this.projectMemberService.page(this.page, this.size, this.currentProject.name).subscribe(res => {
            this.items = res.items;
            this.loading = false;
        });
    }


    changeMembersRole(selected, role) {
        selected.forEach(item => {
            const create = new ProjectMemberCreate();
            create.projectName = this.currentProject.name;
            create.userName = item.userName;
            create.role = role;
            this.batchItems.push(create);
        });
        if (this.batchItems.length < 1) {
            return;
        }
        this.projectMemberService.batch('update', this.batchItems, this.currentProject.name).subscribe(res => {
            this.pageBy();
            this.commonAlertService.showAlert(this.translateService.instant('APP_UPDATE_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
