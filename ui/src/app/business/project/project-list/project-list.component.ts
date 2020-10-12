import {Component, OnInit} from '@angular/core';
import {Project} from '../project';
import {ProjectService} from '../project.service';
import {PermissionService} from '../../../shared/auth/permission.service';
import {BaseModelDirective} from '../../../shared/class/BaseModelDirective';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {AlertLevels} from '../../../layout/common-alert/alert';
import {TranslateService} from '@ngx-translate/core';

@Component({
    selector: 'app-project-list',
    templateUrl: './project-list.component.html',
    styleUrls: ['./project-list.component.css']
})
export class ProjectListComponent extends BaseModelDirective<Project> implements OnInit {


    constructor(private projectService: ProjectService,
                private permissionService: PermissionService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {
        super(projectService);
    }

    ngOnInit(): void {
        super.ngOnInit();
    }

    getProjectRole(projectName: string) {
        return this.permissionService.getProjectRole(projectName);
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
