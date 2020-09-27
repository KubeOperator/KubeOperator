import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {ProjectResource, ProjectResourceCheck, ProjectResourceCreateRequest} from '../project-resource';
import {ProjectResourceService} from '../project-resource.service';
import {ActivatedRoute} from '@angular/router';
import {ResourceTypes} from '../../../../constant/shared.const';
import {Project} from '../../project';
import {NgForm} from '@angular/forms';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../../layout/common-alert/alert';

@Component({
    selector: 'app-project-resource-create',
    templateUrl: './project-resource-create.component.html',
    styleUrls: ['./project-resource-create.component.css']
})
export class ProjectResourceCreateComponent extends BaseModelDirective<ProjectResource> implements OnInit {

    resourceType: string;
    opened: boolean;
    isSubmitGoing = false;
    item: ProjectResourceCreateRequest = new ProjectResourceCreateRequest();
    currentProject: Project = new Project();
    resources: ProjectResourceCheck[] = [];
    @ViewChild('resourceForm') resourceForm: NgForm;
    @Output() created = new EventEmitter();

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

    open(resourceType) {
        this.resourceType = resourceType;
        this.projectResourceService.listResources(this.resourceType, this.currentProject.name).subscribe(res => {
            if (res.length === 0) {
                this.commonAlertService.showAlert(this.translateService.instant('APP_PROJECT_RESOURCE'), AlertLevels.ERROR);
                return;
            }
            for (const re of res) {
                const resource = new ProjectResourceCheck();
                resource.checked = false;
                resource.data = re;
                this.resources.push(resource);
            }
            this.item = new ProjectResourceCreateRequest();
            this.opened = true;
        });
    }

    onCancel() {
        this.opened = false;
        this.resources = [];
        this.item = new ProjectResourceCreateRequest();
        this.resourceForm.resetForm(this.resources);
    }

    onSubmit() {

        if (this.resources.length === 0) {
            return;
        }
        const items = [];
        for (const re of this.resources) {
            if (re.checked) {
                const item = new ProjectResourceCreateRequest();
                item.projectId = this.currentProject.id;
                item.resourceType = this.resourceType;
                item.resourceName = re.data['name'];
                items.push(item);
            }
        }
        this.isSubmitGoing = true;
        this.projectResourceService.batch('create', items).subscribe(res => {
            this.isSubmitGoing = false;
            this.onCancel();
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
            this.created.emit();
        }, error => {
            this.isSubmitGoing = false;
            this.modalAlertService.showAlert(error, AlertLevels.ERROR);
        });
    }

}
