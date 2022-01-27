import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {Project, ProjectCreateRequest} from '../project';
import {ProjectService} from '../project.service';
import {BaseModelDirective} from '../../../shared/class/BaseModelDirective';
import {ModalAlertService} from '../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {NgForm} from '@angular/forms';
import {AlertLevels} from '../../../layout/common-alert/alert';
import {NamePattern} from '../../../constant/pattern';
import {SessionService} from '../../../shared/auth/session.service';

@Component({
    selector: 'app-project-create',
    templateUrl: './project-create.component.html',
    styleUrls: ['./project-create.component.css']
})
export class ProjectCreateComponent extends BaseModelDirective<Project> implements OnInit {

    namePattern = NamePattern;
    opened = false;
    item: ProjectCreateRequest = new ProjectCreateRequest();
    isSubmitGoing = false;
    @Output() created = new EventEmitter();
    @ViewChild('projectForm') hostForm: NgForm;
    user;

    constructor(private projectService: ProjectService,
                private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService,
                private sessionService: SessionService) {
        super(projectService);
    }

    ngOnInit(): void {
        this.sessionService.getProfile().subscribe(res => {
            if (res != null) {
                this.user = res.user;
            }
        })
    }

    open() {
        this.opened = true;
        this.item = new ProjectCreateRequest();
    }

    onCancel() {
        this.opened = false;
        this.isSubmitGoing = false;
        this.hostForm.resetForm(this.item);
    }

    onSubmit() {
        this.item.userName = this.user.name;
        this.projectService.create(this.item).subscribe(res => {
            this.onCancel();
            this.created.emit();
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
            // window.location.reload();
        }, error => {
            this.isSubmitGoing = false;
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
