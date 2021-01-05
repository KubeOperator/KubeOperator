import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {ProjectMember, ProjectMemberCreate} from '../project-member';
import {ProjectMemberService} from '../project-member.service';
import {NgForm} from '@angular/forms';
import {ActivatedRoute} from '@angular/router';
import {Project} from '../../project';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../../layout/common-alert/alert';

@Component({
    selector: 'app-project-member-create',
    templateUrl: './project-member-create.component.html',
    styleUrls: ['./project-member-create.component.css']
})
export class ProjectMemberCreateComponent extends BaseModelDirective<ProjectMember> implements OnInit {

    opened = false;
    item: ProjectMemberCreate = new ProjectMemberCreate();
    selectUsers: string[] = [];
    currentProject: Project = new Project();
    @Output() created = new EventEmitter();
    @ViewChild('memberForm') memberForm: NgForm;


    constructor(private projectMemberService: ProjectMemberService,
                private route: ActivatedRoute,
                private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {
        super(projectMemberService);
    }

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentProject = data.project;
        });
    }

    open() {
        this.opened = true;
        this.item = new ProjectMemberCreate();
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        this.item.projectName = this.currentProject.name;
        this.projectMemberService.create(this.item, this.currentProject.name).subscribe(res => {
            this.opened = false;
            this.created.emit();
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    leaveInput() {
        this.selectUsers = [];
    }

    handleValidation() {
        this.projectMemberService.getUsers(this.item.userName, this.currentProject.name).subscribe(res => {
            this.selectUsers = res.items;
        });
    }

    selectedName(name) {
        this.item.userName = name;
        this.selectUsers = [];
    }
}
