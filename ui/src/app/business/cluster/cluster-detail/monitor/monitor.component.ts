import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {ClusterService} from '../../cluster.service';
import {Cluster, ClusterMonitor} from '../../cluster';

@Component({
    selector: 'app-monitor',
    templateUrl: './monitor.component.html',
    styleUrls: ['./monitor.component.css']
})
export class MonitorComponent implements OnInit {

    currentCluster: Cluster;
    monitor: ClusterMonitor;

    constructor(private route: ActivatedRoute, private clusterService: ClusterService) {
    }

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
            this.getMonitor();
        });
    }

    getMonitor() {
        this.clusterService.monitor(this.currentCluster.name).subscribe(data => {
            console.log(data);
            this.monitor = data;
        });
    }

}
