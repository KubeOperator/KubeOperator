import {Component, OnInit, ViewChild} from '@angular/core';
import {ProjectResourceListComponent} from './project-resource-list/project-resource-list.component';
import {ProjectResourceCreateComponent} from './project-resource-create/project-resource-create.component';
import {ProjectResourceDeleteComponent} from './project-resource-delete/project-resource-delete.component';

@Component({
    selector: 'app-project-resource',
    templateUrl: './project-resource.component.html',
    styleUrls: ['./project-resource.component.css']
})
export class ProjectResourceComponent implements OnInit {

    @ViewChild(ProjectResourceListComponent, {static: true})
    list: ProjectResourceListComponent;

    @ViewChild(ProjectResourceCreateComponent, {static: true})
    create: ProjectResourceCreateComponent;

    @ViewChild(ProjectResourceDeleteComponent, {static: true})
    delete: ProjectResourceDeleteComponent;

    constructor() {
    }

    ngOnInit(): void {
    }

    refresh() {
        this.list.reset();
        this.list.pageBy();
    }

    openCreate(resourceType) {
        this.create.open(resourceType);
    }

    openDelete(deleteItem) {
        this.delete.open(deleteItem);
    }
}
