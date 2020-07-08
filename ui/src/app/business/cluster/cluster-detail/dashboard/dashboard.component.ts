import {Component, ElementRef, Input, OnInit, ViewChild} from '@angular/core';
import {Cluster, ClusterMonitor} from "../../cluster";
import {ActivatedRoute} from "@angular/router";
import {ClusterService} from "../../cluster.service";
import {DomSanitizer} from "@angular/platform-browser";
import {ClusterTool} from "../tools/tools";
import {ToolsService} from "../tools/tools.service";

@Component({
    selector: 'app-dashboard',
    templateUrl: './dashboard.component.html',
    styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {

    currentCluster: Cluster;
    ready = false;
    toolName = 'dashboard';
    item: ClusterTool;

    constructor(private toolService: ToolsService, private route: ActivatedRoute) {
    }

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
            this.toolService.list(this.currentCluster.name).subscribe(d => {
                for (const tool of d) {
                    if (tool.name === this.toolName) {
                        this.item = tool;
                        this.ready = tool.status === 'Running';
                    }
                }
            });
        });
    }

}
