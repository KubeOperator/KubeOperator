import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {ClusterService} from '../../../cluster.service';
import {Cluster} from '../../../cluster';
import {ActivatedRoute} from '@angular/router';

@Component({
    selector: 'app-monitor-enable',
    templateUrl: './monitor-enable.component.html',
    styleUrls: ['./monitor-enable.component.css']
})
export class MonitorEnableComponent implements OnInit {

    constructor(private clusterService: ClusterService, private route: ActivatedRoute) {
    }

    currentCluster: Cluster;
    @Output() created = new EventEmitter();

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
        });
    }

    onSubmit() {
        this.clusterService.createMonitor(this.currentCluster.name).subscribe(data => {
            this.created.emit();
        });
    }
}
