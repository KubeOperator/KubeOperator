import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {ProjectResourceDeleteRequest} from '../project-resource';
import {ProjectResourceService} from '../project-resource.service';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../../layout/common-alert/alert';
import {ActivatedRoute} from '@angular/router';
import {Project} from '../../project';

@Component({
    selector: 'app-project-resource-delete',
    templateUrl: './project-resource-delete.component.html',
    styleUrls: ['./project-resource-delete.component.css']
})
export class ProjectResourceDeleteComponent extends BaseModelDirective<any> implements OnInit {

    opened = false;
    isSubmitGoing = false;
    currentProject = new Project();
    resourceType: string;
    @Output() deleted = new EventEmitter();

    constructor(private projectResourceService: ProjectResourceService,
                private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService,
                private route: ActivatedRoute) {
        super(projectResourceService);
    }

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentProject = data.project;
        });
    }

    open(deleteItem) {
        this.items = deleteItem.items;
        this.resourceType = deleteItem.resourceType;
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        const deleteItems = [];
        for (const item of this.items) {
            const projectResource = new ProjectResourceDeleteRequest();
            projectResource.projectId = this.currentProject.id;
            projectResource.resourceName = item.name;
            projectResource.resourceType = this.resourceType;
            deleteItems.push(projectResource);
        }

        this.projectResourceService.batch('delete', deleteItems).subscribe(res => {
            this.onCancel();
            this.deleted.emit();
            this.commonAlertService.showAlert(this.translateService.instant('APP_DELETE_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.onCancel();
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
