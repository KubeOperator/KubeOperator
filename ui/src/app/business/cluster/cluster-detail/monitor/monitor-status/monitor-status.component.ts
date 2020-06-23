import {Component, EventEmitter, OnDestroy, OnInit, Output} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {ClusterService} from '../../../cluster.service';
import {Cluster, ClusterMonitor} from '../../../cluster';

@Component({
    selector: 'app-monitor-status',
    templateUrl: './monitor-status.component.html',
    styleUrls: ['./monitor-status.component.css']
})
export class MonitorStatusComponent implements OnInit, OnDestroy {

    currentCluster: Cluster;
    monitor: ClusterMonitor;
    cancel: any;

    @Output() completed = new EventEmitter();

    constructor(private route: ActivatedRoute, private clusterService: ClusterService) {
    }

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
            this.polling();
        });
    }

    ngOnDestroy(): void {
        this.cancel();
    }

    polling() {
        this.cancel = setInterval(() => {
            this.refresh();
            if (this.monitor.status !== 'Initializing') {
                this.cancel();
                this.completed.emit();
            }
        }, 1000);
    }


    refresh() {
        this.clusterService.monitor(this.currentCluster.name).subscribe(data => {
            this.monitor = data;
        });
    }

}
