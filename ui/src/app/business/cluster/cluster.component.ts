import {Component, OnInit, ViewChild} from '@angular/core';
import {ClusterCreateComponent} from './cluster-create/cluster-create.component';
import {ClusterListComponent} from './cluster-list/cluster-list.component';
import {ClusterDeleteComponent} from './cluster-delete/cluster-delete.component';
import {Cluster} from './cluster';
import {ClusterConditionComponent} from './cluster-condition/cluster-condition.component';
import {ClusterUpgradeComponent} from './cluster-upgrade/cluster-upgrade.component';
import {ClusterHealthCheckComponent} from "./cluster-health-check/cluster-health-check.component";

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

    @ViewChild(ClusterConditionComponent, {static: true})
    condition: ClusterConditionComponent;

    @ViewChild(ClusterListComponent, {static: true})
    list: ClusterListComponent;

    @ViewChild(ClusterUpgradeComponent, {static: true})
    upgrade: ClusterUpgradeComponent;

    @ViewChild(ClusterHealthCheckComponent, {static: true})
    healthCheck: ClusterHealthCheckComponent;

    ngOnInit(): void {
    }

    openCreate() {
        this.create.open();
    }

    openHealthCheck(cluster: Cluster) {
        this.healthCheck.open(cluster);
    }

    openDelete(items: Cluster[]) {
        this.delete.open(items);
    }

    openStatusDetail(cluster: Cluster) {
        this.condition.open(cluster);
    }

    refresh() {
        this.list.reset();
        this.list.pageBy();
    }

    openUpgrade(item: Cluster) {
        this.upgrade.open(item);
    }
}
