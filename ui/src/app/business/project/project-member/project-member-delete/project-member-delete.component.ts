import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {ProjectMember, ProjectMemberCreate} from '../project-member';
import {ProjectMemberService} from '../project-member.service';
import {ActivatedRoute} from '@angular/router';
import {Project} from '../../project';
import {AlertLevels} from '../../../../layout/common-alert/alert';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';

@Component({
    selector: 'app-project-member-delete',
    templateUrl: './project-member-delete.component.html',
    styleUrls: ['./project-member-delete.component.css']
})
export class ProjectMemberDeleteComponent extends BaseModelDirective<ProjectMember> implements OnInit {


    opened = false;
    isSubmitGoing = false;
    currentProject: Project = new Project();
    items: ProjectMember[] = [];
    @Output() delete = new EventEmitter();
    batchItems: ProjectMemberCreate[] = [];

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

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        this.items.forEach(item => {
            const create = new ProjectMemberCreate();
            create.projectName = this.currentProject.name;
            create.userName = item.userName;
            this.batchItems.push(create);
        });
        if (this.batchItems.length < 1) {
            return;
        }
        this.projectMemberService.batch('delete', this.batchItems).subscribe(res => {
            this.delete.emit();
            this.opened = false;
            this.commonAlertService.showAlert(this.translateService.instant('APP_DELETE_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    open(items) {
        this.opened = true;
        this.items = items;
    }

}
