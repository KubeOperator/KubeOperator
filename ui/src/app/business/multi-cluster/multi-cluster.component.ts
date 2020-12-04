import {Component, OnInit, ViewChild} from '@angular/core';
import {MultiClusterRepositoryListComponent} from "./multi-cluster-repository-list/multi-cluster-repository-list.component";
import {MultiClusterRepositoryCreateComponent} from "./multi-cluster-repository-create/multi-cluster-repository-create.component";
import {MultiClusterRepositoryDeleteComponent} from "./multi-cluster-repository-delete/multi-cluster-repository-delete.component";
import {MultiClusterRepository} from "./multi-cluster-repository";
import {MultiClusterRepositoryErrorDetailComponent} from "./multi-cluster-repository-error-detail/multi-cluster-repository-error-detail.component";

@Component({
    selector: 'app-multi-cluster',
    templateUrl: './multi-cluster.component.html',
    styleUrls: ['./multi-cluster.component.css']
})
export class MultiClusterComponent implements OnInit {

    constructor() {
    }

    @ViewChild(MultiClusterRepositoryListComponent, {static: true})
    list: MultiClusterRepositoryListComponent;

    @ViewChild(MultiClusterRepositoryCreateComponent, {static: true})
    create: MultiClusterRepositoryCreateComponent;

    @ViewChild(MultiClusterRepositoryDeleteComponent, {static: true})
    delete: MultiClusterRepositoryDeleteComponent;

    @ViewChild(MultiClusterRepositoryErrorDetailComponent, {static: true})
    errDetail: MultiClusterRepositoryErrorDetailComponent;

    ngOnInit(): void {
    }

    openCreate() {
        this.create.open();
    }

    openDelete(items: MultiClusterRepository[]) {
        this.delete.open(items);
    }

    refresh() {
        this.list.reset();
        this.list.refresh();
    }

    openErrorDetail(msg: string) {
        this.errDetail.open(msg);
    }

}
