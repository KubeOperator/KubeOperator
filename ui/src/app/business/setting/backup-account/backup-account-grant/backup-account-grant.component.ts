import {Component, OnInit} from '@angular/core';
import {ProjectService} from '../../../project/project.service';
import {BackupAccount, Project} from '../backup-account';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {ProjectResourceService} from '../../../project/project-resource/project-resource.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../../layout/common-alert/alert';
import { ProjectResourceCreateRequest } from '../../../project/project-resource/project-resource';


@Component({
    selector: 'app-backup-account-grant',
    templateUrl: './backup-account-grant.component.html',
    styleUrls: ['./backup-account-grant.component.css']
})
export class BackupAccountGrantComponent implements OnInit {

    opened = false;
    projects: Project[] = [];
    backupAccounts: BackupAccount[] = [];
    isSubmitGoing = false;
    resourceType: string = 'BACKUP_ACCOUNT';

    constructor(
        private projectResourceService: ProjectResourceService,
        private modalAlertService: ModalAlertService,
        private commonAlertService: CommonAlertService,
        private translateService: TranslateService,
        private projectService: ProjectService) {
    }

    ngOnInit(): void {
    }

    listProjects() {
        this.projectService.list().subscribe(res => {
            this.projects = [];
            for (let pro of res.items) {
                this.projects.push({
                    id: pro.id, 
                    name: pro.name,
                    checked: false,
                })
            }
        })
    }

    open(items) {
        this.backupAccounts = items;
        this.listProjects();
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        this.isSubmitGoing = true;
        if (this.projects.length === 0) {
            return;
        }
        const items = [];
        for (const pro of this.projects) {
            if (pro.checked) {
                for (const backup of this.backupAccounts) {
                    const item = new ProjectResourceCreateRequest();
                    item.projectId = pro.id;
                    item.resourceType = this.resourceType;
                    item.resourceName = backup.name;
                    items.push(item);
                }
            }
        }

        this.projectResourceService.batch('create', items).subscribe(res => {
            this.isSubmitGoing = false;
            this.onCancel();
            this.commonAlertService.showAlert(this.translateService.instant('APP_GRANT_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.isSubmitGoing = false;
            this.modalAlertService.showAlert(error, AlertLevels.ERROR);
        });
    }
}
