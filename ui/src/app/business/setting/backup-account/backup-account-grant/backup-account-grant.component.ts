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
    projectIndex: number;
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
        this.projectResourceService.listResources(this.resourceType, this.projects[this.projectIndex].name).subscribe(res => {
            const items = [];
            for (const ho of this.backupAccounts) {
                let isExit = false;
                for (const re of res) {
                    if (ho.name === re.name) {
                        isExit = true;
                        break;
                    }
                }
                if (!isExit) {
                    this.isSubmitGoing = false;
                    this.modalAlertService.showAlert(this.translateService.instant('APP_BACKUP_ACCOUNT_BOUND'), AlertLevels.ERROR);
                    return;
                } else {
                    const item = new ProjectResourceCreateRequest();
                    item.projectId = this.projects[this.projectIndex].id;
                    item.resourceType = this.resourceType;
                    item.resourceName = ho.name;
                    items.push(item);
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
        });
    }
}
