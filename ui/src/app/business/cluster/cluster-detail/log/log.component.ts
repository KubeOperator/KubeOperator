import {Component, OnInit, ViewChild} from '@angular/core';
import {Cluster} from "../../cluster";
import {ActivatedRoute} from "@angular/router";
import {LogDetailComponent} from "./log-detail/log-detail.component";
import {LogListComponent} from "./log-list/log-list.component";
import {ClusterLog} from "./log";

@Component({
    selector: 'app-log',
    templateUrl: './log.component.html',
    styleUrls: ['./log.component.css']
})
export class LogComponent implements OnInit {

    constructor(private route: ActivatedRoute) {
    }

    currentCluster: Cluster;

    @ViewChild(LogDetailComponent, {static: true})
    detail: LogDetailComponent;

    @ViewChild(LogListComponent, {static: true})
    list: LogListComponent;


    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
        });
    }

    openDetail(item: ClusterLog) {
        this.detail.open(item);
    }

    refresh() {
        this.list.refresh();
    }


}
