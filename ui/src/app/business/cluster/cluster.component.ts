import {Component, OnInit, ViewChild} from '@angular/core';
import {ClusterCreateComponent} from './cluster-create/cluster-create.component';
import {ClusterListComponent} from './cluster-list/cluster-list.component';
import {ClusterDeleteComponent} from './cluster-delete/cluster-delete.component';
import {Cluster} from './cluster';

@Component({
    selector: 'app-cluster',
    templateUrl: './cluster.component.html',
    styleUrls: ['./cluster.component.css']
})
export class ClusterComponent implements OnInit {

    constructor() {
    }

    @ViewChild(ClusterCreateComponent, {static: true})
    create: ClusterCreateComponent;

    @ViewChild(ClusterDeleteComponent, {static: true})
    delete: ClusterDeleteComponent;

    @ViewChild(ClusterListComponent, {static: true})
    list: ClusterListComponent;

    ngOnInit(): void {
    }

    openCreate() {
        this.create.open();
    }

    openDelete(items: Cluster[]) {
        this.delete.open(items);
    }

    refresh() {
        this.list.reset();
        this.list.refresh();
    }

}
