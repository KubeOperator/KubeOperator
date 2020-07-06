import {Component, EventEmitter, OnInit, ViewChild} from '@angular/core';
import {ToolsService} from "./tools.service";
import {ActivatedRoute} from "@angular/router";
import {ClusterTool} from "./tools";
import {Cluster} from "../../cluster";
import {ToolsListComponent} from "./tools-list/tools-list.component";
import {PrometheusEnableComponent} from "./tools-list/prometheus-enable/prometheus-enable.component";

@Component({
    selector: 'app-tools',
    templateUrl: './tools.component.html',
    styleUrls: ['./tools.component.css']
})
export class ToolsComponent implements OnInit {


    @ViewChild(ToolsListComponent, {static: true})
    list: ToolsListComponent;

    @ViewChild(PrometheusEnableComponent, {static: true})
    prometheus: PrometheusEnableComponent;

    constructor(private service: ToolsService, private route: ActivatedRoute) {
    }

    currentCluster: Cluster;

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
        });
    }

    openEnable(item: ClusterTool) {
        switch (item.name) {
            case 'Prometheus':
                this.prometheus.open(item);
                break;
        }
    }

    refresh() {
        this.list.refresh();
    }

}
