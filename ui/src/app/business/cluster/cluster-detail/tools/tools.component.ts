import {Component, EventEmitter, OnInit, ViewChild} from '@angular/core';
import {ToolsService} from "./tools.service";
import {ActivatedRoute} from "@angular/router";
import {ClusterTool} from "./tools";
import {Cluster} from "../../cluster";
import {ToolsListComponent} from "./tools-list/tools-list.component";
import {ToolsEnableComponent} from "./tools-enable/tools-enable.component";
import {ToolsFailedComponent} from "./tools-failed/tools-failed.component";

@Component({
    selector: 'app-tools',
    templateUrl: './tools.component.html',
    styleUrls: ['./tools.component.css']
})
export class ToolsComponent implements OnInit {


    @ViewChild(ToolsListComponent, {static: true})
    list: ToolsListComponent;

    @ViewChild(ToolsEnableComponent, {static: true})
    enable: ToolsEnableComponent;

    @ViewChild(ToolsFailedComponent, {static: true})
    failed: ToolsFailedComponent;

    constructor(private service: ToolsService, private route: ActivatedRoute) {
    }

    currentCluster: Cluster;

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
        });
    }

    openEnable(item: ClusterTool) {
        this.enable.open(item);
    }

    openFailed(item: ClusterTool) {
        this.failed.open(item);
    }

    refresh() {
        this.list.refresh();
    }

}
