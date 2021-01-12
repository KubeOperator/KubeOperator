import {Component, OnInit} from '@angular/core';
import {ClusterService} from '../../../cluster.service';
import {Cluster} from '../../../cluster';
import {ClusterLog} from '../../log/log';
import {ActivatedRoute} from '@angular/router';

@Component({
    selector: 'app-backup-log',
    templateUrl: './backup-log.component.html',
    styleUrls: ['./backup-log.component.css']
})
export class BackupLogComponent implements OnInit {

    constructor(private clusterService: ClusterService, private route: ActivatedRoute) {
    }

    currentCluster: Cluster;
    loading = false;
    items: ClusterLog[] = [];
    opened = false;
    detailItem: ClusterLog = new ClusterLog();

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
            this.refresh();
        });
    }

    refresh() {
        this.loading = true;
        this.clusterService.log(this.currentCluster.name).subscribe(data => {
            this.items = data;
            this.loading = false;
        });
    }

    onDetail(item: ClusterLog) {
        item.message = item.message.replace(/[\\]/g, '');
        this.detailItem = item;
        this.opened = true;
    }
}
