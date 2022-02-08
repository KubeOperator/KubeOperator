import {Component, OnInit} from '@angular/core';
import {Project} from '../project';
import {ProjectService} from '../project.service';
import {PermissionService} from '../../../shared/auth/permission.service';
import {BaseModelDirective} from '../../../shared/class/BaseModelDirective';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {AlertLevels} from '../../../layout/common-alert/alert';
import {TranslateService} from '@ngx-translate/core';
import {SessionService} from "../../../shared/auth/session.service";
import {SessionUser} from "../../../shared/auth/session-user";
import {ProjectMemberService} from "../project-member/project-member.service";
import {ProjectMember} from "../project-member/project-member";
import {Router} from '@angular/router';
import {CommonRoutes} from '../../../constant/route';

@Component({
    selector: 'app-project-list',
    templateUrl: './project-list.component.html',
    styleUrls: ['./project-list.component.css']
})
export class ProjectListComponent extends BaseModelDirective<Project> implements OnInit {


    constructor(private projectService: ProjectService, private router: Router,
                private permissionService: PermissionService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService,
                private sessionService: SessionService) {
        super(projectService);
    }

    user: SessionUser;

    ngOnInit(): void {
        super.ngOnInit();
        this.sessionService.getProfile().subscribe(res => {
            this.user = res.user;
        }, error => {
            this.sessionService.clear();
            this.router.navigateByUrl(CommonRoutes.LOGIN).then();
        })
    }

    onUpdate(item: Project) {
        this.permissionService.authOperate('PROJECT.UPDATE', item.name).then(result => {
            if (result) {
                super.onUpdate(item);
            } else {
                this.commonAlertService.showAlert(this.translateService.instant('APP_NO_AUTH'), AlertLevels.ERROR);
            }
        });
    }

    onDelete() {
        let result = true;
        let clusterName = '';
        for (const item of this.selected) {
            const auth = this.permissionService.authOp('PROJECT.DELETE', item.name);
            if (!auth) {
                result = false;
                clusterName = clusterName + item.name + ',';
            }
        }
        if (result) {
            super.onDelete();
        } else {
            this.commonAlertService.showAlert(this.translateService.instant('APP_CLUSTER') + clusterName + this.translateService.instant('APP_NO_AUTH'), AlertLevels.ERROR);
        }
    }
}
