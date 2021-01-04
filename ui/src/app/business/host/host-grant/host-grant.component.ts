import {Component, OnInit} from '@angular/core';
import {ProjectService} from '../../project/project.service';
import {Host, Project} from '../host';
import {ModalAlertService} from '../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {ProjectResourceService} from '../../project/project-resource/project-resource.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../layout/common-alert/alert';
import { ProjectResourceCreateRequest } from '../../project/project-resource/project-resource';

@Component({
    selector: 'app-host-grant',
    templateUrl: './host-grant.component.html',
    styleUrls: ['./host-grant.component.css']
})
export class HostGrantComponent implements OnInit {

    opened = false;
    projects: Project[] = [];
    projectIndex: number;
    hosts: Host[] = [];
    isSubmitGoing = false;
    resourceType: string = 'HOST';

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
            for (const pro of res.items) {
                this.projects.push({
                    id: pro.id, 
                    name: pro.name,
                })
            }
        })
    }

    open(items) {
        this.hosts = items;
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
            for (const ho of this.hosts) {
                let isExit = false;
                for (const re of res) {
                    if (ho.name === re.name) {
                        isExit = true;
                        break;
                    }
                }
                if (!isExit) {
                    this.isSubmitGoing = false;
                    this.modalAlertService.showAlert(this.translateService.instant('APP_HOST_BOUND'), AlertLevels.ERROR);
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
