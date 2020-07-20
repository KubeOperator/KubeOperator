import {Component, OnInit, ViewChild} from '@angular/core';
import {ProjectListComponent} from './project-list/project-list.component';
import {ProjectCreateComponent} from './project-create/project-create.component';
import {ProjectDeleteComponent} from './project-delete/project-delete.component';
import {ProjectUpdateComponent} from './project-update/project-update.component';

@Component({
    selector: 'app-project',
    templateUrl: './project.component.html',
    styleUrls: ['./project.component.css']
})
export class ProjectComponent implements OnInit {

    @ViewChild(ProjectListComponent, {static: true})
    list: ProjectListComponent;

    @ViewChild(ProjectCreateComponent, {static: true})
    create: ProjectCreateComponent;

    @ViewChild(ProjectDeleteComponent, {static: true})
    delete: ProjectDeleteComponent;

    @ViewChild(ProjectUpdateComponent, {static: true})
    update: ProjectUpdateComponent;

    constructor() {
    }

    ngOnInit(): void {
    }

    openCreate() {
        this.create.open();
    }

    openDelete(items) {
        this.delete.open(items);
    }

    refresh() {
        this.list.reset();
        this.list.refresh();
    }

    openUpdate(item) {
        this.update.open(item);
    }
}
