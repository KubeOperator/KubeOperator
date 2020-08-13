import {Component, OnInit, ViewChild} from '@angular/core';
import {ProjectMemberListComponent} from './project-member-list/project-member-list.component';
import {ProjectMemberCreateComponent} from './project-member-create/project-member-create.component';
import {ProjectMemberDeleteComponent} from './project-member-delete/project-member-delete.component';

@Component({
    selector: 'app-project-member',
    templateUrl: './project-member.component.html',
    styleUrls: ['./project-member.component.css']
})
export class ProjectMemberComponent implements OnInit {

    @ViewChild(ProjectMemberListComponent, {static: true})
    list: ProjectMemberListComponent;

    @ViewChild(ProjectMemberCreateComponent, {static: true})
    create: ProjectMemberCreateComponent;

    @ViewChild(ProjectMemberDeleteComponent, {static: true})
    delete: ProjectMemberDeleteComponent;

    constructor() {
    }

    ngOnInit(): void {
    }

    refresh() {
        this.list.reset();
        this.list.pageBy();
    }

    openCreate() {
        this.create.open();
    }

    openDelete(items) {
        this.delete.open(items);
    }
}
