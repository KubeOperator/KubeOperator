import {Component, OnInit, ViewChild} from '@angular/core';
import {MultiClusterRepository} from "../../multi-cluster-repository";
import {ActivatedRoute} from "@angular/router";
import {MultiClusterLogListComponent} from "./multi-cluster-log-list/multi-cluster-log-list.component";
import {MultiClusterLogDetailComponent} from "./multi-cluster-log-detail/multi-cluster-log-detail.component";

@Component({
    selector: 'app-multi-cluster-log',
    templateUrl: './multi-cluster-log.component.html',
    styleUrls: ['./multi-cluster-log.component.css']
})
export class MultiClusterLogComponent implements OnInit {

    constructor(private route: ActivatedRoute) {
    }

    currentRepository: MultiClusterRepository;

    @ViewChild(MultiClusterLogListComponent, {static: true})
    list: MultiClusterLogListComponent;

    @ViewChild(MultiClusterLogDetailComponent, {static: true})
    detail: MultiClusterLogDetailComponent;

    ngOnInit(): void {
        this.route.parent.data.subscribe(d => {
            this.currentRepository = d.repo;
        });
    }

    openDetail(logId: string) {
        this.detail.open(logId);
    }
}
