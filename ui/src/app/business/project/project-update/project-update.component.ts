import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {BaseModelDirective} from '../../../shared/class/BaseModelDirective';
import {Project} from '../project';
import {ProjectService} from '../project.service';
import {ModalAlertService} from '../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {User} from '../../user/user';
import {NgForm} from '@angular/forms';
import {AlertLevels} from '../../../layout/common-alert/alert';

@Component({
    selector: 'app-project-update',
    templateUrl: './project-update.component.html',
    styleUrls: ['./project-update.component.css']
})
export class ProjectUpdateComponent extends BaseModelDirective<Project> implements OnInit {

    opened = false;
    item: Project = new Project();
    isSubmitGoing = false;
    @ViewChild('projectForm') projectForm: NgForm;
    @Output()
    update = new EventEmitter();
    oldName: string;

    constructor(private projectService: ProjectService, private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService, private translateService: TranslateService) {
        super(projectService);
    }

    ngOnInit(): void {
    }

    open(item) {
        this.opened = true;
        this.item = item;
        this.oldName = this.item.name;
    }

    onCancel() {
        this.opened = false;
        this.item = new Project();
        this.projectForm.resetForm(this.item);
        this.isSubmitGoing = false;
    }

    onSubmit() {
        this.isSubmitGoing = true;
        this.projectService.update(this.oldName, this.item).subscribe(data => {
            this.onCancel();
            this.update.emit();
            this.commonAlertService.showAlert(this.translateService.instant('APP_UPDATE_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.isSubmitGoing = false;
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
