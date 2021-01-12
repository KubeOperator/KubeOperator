import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {Cluster} from "../../../cluster";
import {ClusterLog} from "../log";
import {ClusterService} from "../../../cluster.service";

@Component({
    selector: 'app-log-list',
    templateUrl: './log-list.component.html',
    styleUrls: ['./log-list.component.css']
})
export class LogListComponent implements OnInit {

    constructor(private clusterService: ClusterService) {
    }

    @Input() currentCluster: Cluster;
    loading = false;
    items: ClusterLog[] = [];
    @Output() detailEvent = new EventEmitter<ClusterLog>();

    ngOnInit(): void {
        this.refresh();
    }

    refresh() {
        this.loading = true;
        this.clusterService.log(this.currentCluster.name).subscribe(data => {
            this.items = data;
            this.loading = false;
        });
    }

    onDetail(item: ClusterLog) {
        this.detailEvent.emit(item);
    }

}
