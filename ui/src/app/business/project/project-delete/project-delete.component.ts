import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelComponent} from '../../../shared/class/BaseModelComponent';
import {Project} from '../project';
import {ProjectService} from '../project.service';
import {AlertLevels} from '../../../layout/common-alert/alert';
import {ModalAlertService} from '../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';

@Component({
    selector: 'app-project-delete',
    templateUrl: './project-delete.component.html',
    styleUrls: ['./project-delete.component.css']
})
export class ProjectDeleteComponent extends BaseModelComponent<Project> implements OnInit {

    opened = false;
    items: Project[] = [];
    @Output() deleted = new EventEmitter();

    constructor(private projectService: ProjectService,
                private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {
        super(projectService);
    }

    ngOnInit(): void {
    }

    open(items) {
        this.opened = true;
        this.items = items;
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        this.projectService.batch('delete', this.items).subscribe(data => {
            this.deleted.emit();
            this.opened = false;
            this.commonAlertService.showAlert(this.translateService.instant('APP_DELETE_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.modalAlertService.showAlert(error, AlertLevels.ERROR);
        });
    }
}
