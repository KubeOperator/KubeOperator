import {Component, OnInit} from '@angular/core';
import {Cluster, ClusterHealthCheck, ClusterRecoverItem} from "../cluster";
import {ClusterService} from "../cluster.service";

@Component({
    selector: 'app-cluster-health-check',
    templateUrl: './cluster-health-check.component.html',
    styleUrls: ['./cluster-health-check.component.css']
})
export class ClusterHealthCheckComponent implements OnInit {

    constructor(private clusterService: ClusterService) {
    }

    opened = false;
    cluster: Cluster = new Cluster();
    item: ClusterHealthCheck = new ClusterHealthCheck();
    loading = false;
    isRecover = false;
    recoverItems: ClusterRecoverItem[] = [];

    ngOnInit(): void {
    }

    open(cluster: Cluster) {
        this.cluster = cluster;
        this.opened = true;
        this.loading = true;
        this.isRecover = false;
        this.recoverItems = [];
        this.clusterService.healthCheck(cluster.name).subscribe(data => {
            this.loading = false;
            this.item = data;
        });
    }

    onRecover() {
        this.loading = true;
        this.isRecover = true;
        this.item = new ClusterHealthCheck();
        this.clusterService.recover(this.cluster.name).subscribe(data => {
            this.recoverItems = data;
            this.loading = false;
        });
    }
}
